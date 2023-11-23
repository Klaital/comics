package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/klaital/comics/pkg/comics"
	"github.com/klaital/comics/pkg/config"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type handlerConfig struct {
	logger           *log.Entry
	db               *sql.DB
	ctx              context.Context
	comicReadUpdates chan int

	// used to generate callback links
	hostname string
	port     int

	cacheMutex      sync.RWMutex
	comicsDataCache map[int]comics.ComicRecord // used to preserve data between the fetch all and update read time calls
}

func (cfg *handlerConfig) cacheComicsData(comicSet []comics.ComicRecord) {
	// Transform the slice into a map
	cacheData := make(map[int]comics.ComicRecord, len(comicSet))
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
		ctx:              context.Background(),
		comicReadUpdates: make(chan int),
		hostname:         cfg.Hostname,
		port:             cfg.Port,
	}
	db, err := cfg.ConnectPostgres()
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
			err = comics.UpdateReadNow(readComicId, 1, t, db) // TODO: use the real userID!
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
	comicSet, err := comics.FetchComics(serverCfg.ctx, serverCfg.db, 1, nil, nil)
	if err != nil {
		serverCfg.logger.WithError(err).Fatal("Failed to heat up the cache")
	}
	serverCfg.cacheComicsData(comicSet)

	http.HandleFunc("/healthz", serverCfg.healthCheckHandler)
	http.HandleFunc("/api/comics", serverCfg.getActiveComicsHandler)
	http.HandleFunc("/api/read/", serverCfg.readComicHandler)
	http.Handle("/", http.FileServer(http.Dir("./web")))
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func (cfg *handlerConfig) healthCheckHandler(w http.ResponseWriter, req *http.Request) {
	if cfg.db == nil {
		cfg.logger.WithError(errors.New("database not connected")).Error("Database not connected")
		http.Error(w, "database not connected", http.StatusInternalServerError)
		return
	}

	cfg.logger.Debug("Healthy!")
	w.WriteHeader(200)
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

type GetActiveComicsResponse struct {
	Today   []comics.ComicRecord `json:"today"`
	TheRest []comics.ComicRecord `json:"therest"`
}

// getActiveComicsHandler will fetch the active comics list for today
func (cfg *handlerConfig) getActiveComicsHandler(w http.ResponseWriter, req *http.Request) {
	logger := cfg.logger.WithFields(log.Fields{
		"operation": "getActiveComicsHtmlHandler",
	})
	activeComicsFilter := true
	nsfwComicsFilter := false
	if len(cfg.comicsDataCache) == 0 {
		comicSet, err := comics.FetchComics(cfg.ctx, cfg.db, 1, &activeComicsFilter, &nsfwComicsFilter)
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

	today, theRest := comics.SelectMapSubset(cfg.comicsDataCache, comics.GetTodaySelector())
	resp := GetActiveComicsResponse{
		Today:   make([]comics.ComicRecord, 0, len(today)),
		TheRest: make([]comics.ComicRecord, 0, len(theRest)),
	}
	for ordinal := range today {
		resp.Today = append(resp.Today, today[ordinal])
	}
	for ordinal := range theRest {
		resp.TheRest = append(resp.TheRest, theRest[ordinal])
	}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.WithError(err).Error("Failed to marshal comics list")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}
