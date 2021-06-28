package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/klaital/comics/pkg/comics"
	"github.com/klaital/comics/pkg/config"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)
type handlerConfig struct {
	logger *log.Entry
	db *sqlx.DB
	comicReadUpdates chan int

	cacheMutex sync.RWMutex
	comicsDataCache map[int]comics.Comic // used to preserve data between the fetch all and update read time calls
}
func (cfg *handlerConfig) cacheComicsData(comicSet []comics.Comic) {
	// Transform the slice into a map
	cacheData := make(map[int]comics.Comic, len(comicSet))
	for i, c := range comicSet {
		cacheData[c.ID] = comicSet[i]
	}

	// Update the in-memory cache
	cfg.cacheMutex.Lock()
	cfg.comicsDataCache = cacheData
	cfg.cacheMutex.Unlock()
}
func launchServer(cfg *config.Config) {
	serverCfg := handlerConfig{
		logger:           log.NewEntry(log.StandardLogger()),
		db:               nil,
		comicReadUpdates: make(chan int),
	}	
	db, err := cfg.ConnectDB()
	if err != nil {
		serverCfg.logger.WithError(err).Fatal("failed to connect to DB")
	}
	serverCfg.db = db
	// Asynchronously write the read updates to the DB
	go func() {
		for {
			// TODO: add batching
			readComicId := <-serverCfg.comicReadUpdates
			t := time.Now()
			// Update the DB
			err = comics.UpdateReadNow(readComicId, t, db)
			if err != nil {
				// TODO: detect DB disconnect and retry
				log.WithError(err).Error("Failed to update comic in DB")
			}
			// Update the in-memory cache
			serverCfg.cacheMutex.Lock()
			thisComic := serverCfg.comicsDataCache[readComicId]
			thisComic.LastRead = t
			serverCfg.comicsDataCache[readComicId] = thisComic
			serverCfg.cacheMutex.Unlock()
		}
	}()

	// Prime the cache
	comicSet, err := comics.FetchActiveComics(serverCfg.db)
	if err != nil {
		serverCfg.logger.WithError(err).Fatal("Failed to heat up the cache")
	}
	serverCfg.cacheComicsData(comicSet)

	http.HandleFunc("/api/comics", serverCfg.getActiveComicsHandler)
	http.HandleFunc("/api/read/", serverCfg.readComicHandler)
	http.Handle("/", http.FileServer(http.Dir("./web")))
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}



// readComicHandler will update the comic's read record in the DB, then redirect you to the comic's homepage.
// TODO: also auto-refresh the active comics list so that the least-recently-read comics bubble up.
func (cfg *handlerConfig) readComicHandler(w http.ResponseWriter, req *http.Request) {
	// Queue the async DB update
	idStr := strings.TrimPrefix(req.URL.Path, "/api/read/")
	if len(idStr) == 0 {
		cfg.logger.Error("No comic ID from path param")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		cfg.logger.WithError(err).Debug("Failed to parse comic ID from path param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cfg.comicReadUpdates <- int(id64)

	cfg.cacheMutex.RLock()
	thisComic, ok := cfg.comicsDataCache[int(id64)]
	cfg.cacheMutex.RUnlock()
	if !ok {
		cfg.logger.WithField("ComicID", id64).Error("No comic with ID in cache")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(w, req, thisComic.BaseURL, http.StatusSeeOther)
}

// getActiveComicsHandler will render an HTML document with the "to-be-read" comics
func (cfg *handlerConfig) getActiveComicsHandler(w http.ResponseWriter, req *http.Request) {
	logger := cfg.logger.WithFields(log.Fields{
		"operation": "getActiveComicsHandler",
	})
	if len(cfg.comicsDataCache) == 0 {
		comicSet, err := comics.FetchActiveComics(cfg.db)
		if err != nil {
			logger.WithError(err).Error("Failed to fetch active comics")
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			logger = logger.WithField("comics_count", len(comicSet))
			logger.Debug("Updating comicset")
		}
		cfg.cacheComicsData(comicSet)
	}

	tpl := `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		<h2>Comics For Today</h2>
		<table id="comicslist">
			{{range .Today}}<tr class="comic"><td>{{ .DaysAgoNow }}</td><td><a href="http://localhost:8080/api/read/{{ .ID }}">{{.Title}}</a></td></tr>{{end}}
		</table>
		<h2>The Rest</h2>
		<table id="comicslist">
			{{range .Items}}<tr class="comic"><td>{{ .DaysAgoNow }}</td><td><a href="http://localhost:8080/api/read/{{ .ID }}">{{.Title}}</a></td></tr>{{end}}
		</table>
	</body>
</html>`

	t, err := template.New("webpage").Parse(tpl)
	if err != nil {
		logger.WithError(err).Fatal("Failed to parse html template")
	}

	today, theRest := comics.SelectMapSubset(cfg.comicsDataCache, comics.GetTodaySelector())
	data := struct {
		Title string
		Today map[int]comics.Comic
		Items map[int]comics.Comic
	}{
		Title: "AF.net Dynamic Comics Home",
		Today: today,
		Items: theRest,
	}
	err = t.Execute(w, data)
	if err != nil {
		logger.WithError(err).WithField("tplData", data).Fatal("Failed to execute html template")
	} else {
		logger.Debug("Rendered comics list")
	}
}