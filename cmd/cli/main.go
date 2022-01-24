package main

import (
	"context"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/klaital/comics/pkg/comics"
	"github.com/klaital/comics/pkg/config"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

func addNewComic(cfg *config.Config) {
	logger := cfg.LogContext.WithField("operation", "addNewComic")
	c := comics.ComicRecord{
		ID:               0,
		Title:            "",
		BaseURL:          "",
		FirstComicUrl:    nil,
		LatestComicUrl:   nil,
		RssUrl:           nil,
		UpdatesMonday:    false,
		UpdatesTuesday:   false,
		UpdatesWednesday: false,
		UpdatesThursday:  false,
		UpdatesFriday:    false,
		UpdatesSaturday:  false,
		UpdatesSunday:    false,
		Ordinal:          0,
		LastRead:         time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	var updateScheduleString string
	var firstComicUrl string
	var latestComicUrl string
	var rssUrl string
	fs := flag.NewFlagSet("AddNewComic", flag.ExitOnError)
	fs.StringVar(&updateScheduleString, "updates", "", "Update schedule. Specify like 'SuMTuWThFSa'")
	fs.IntVar(&c.Ordinal, "ord", 9999, "Sort ordinal") // TODO: ensure uniqueness, allow for "insert at next available ordinal" or "insert here, bump collisions" modes.
	fs.StringVar(&firstComicUrl, "first", "", "First ComicRecord URL.")
	fs.StringVar(&latestComicUrl, "latest", "", "Most recently read comic URL.")
	fs.StringVar(&latestComicUrl, "rss", "", "ComicRecord's RSS feed URL.")
	fs.StringVar(&c.Title, "title", "", "ComicRecord title")
	fs.StringVar(&c.BaseURL, "base", "", "Base URL for the comic's website. Preferably the 'newest comic' page")
	if err := fs.Parse(os.Args[2:]); err != nil {
		logger.WithError(err).Fatal("failed to parse flagset")
	}

	if len(updateScheduleString) > 0 {
		parseDateString(updateScheduleString, &c)
		logger.WithFields(log.Fields{
			"comic":    c,
			"schedule": updateScheduleString,
		}).Debug("updated schedule")
	}
	if len(firstComicUrl) > 0 {
		c.FirstComicUrl = &firstComicUrl
	}
	if len(latestComicUrl) > 0 {
		c.LatestComicUrl = &latestComicUrl
	}
	if len(rssUrl) > 0 {
		c.RssUrl = &rssUrl
	}

	if err := c.IsValid(); err != nil {
		logger.WithFields(log.Fields{
			"args":  os.Args,
			"comic": c,
		}).WithError(err).Error("Failed to construct a valid comic to insert")
		flag.PrintDefaults()
		return
	}

	db, err := cfg.ConnectDB()
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to db")
	}

	if err = comics.InsertNewComic(&c, db); err != nil {
		logger.WithError(err).Fatal("Failed to insert new comic")
	}

}
func parseDateString(s string, c *comics.ComicRecord) {
	c.UpdatesSunday = strings.Index(s, "Su") > 0
	c.UpdatesMonday = strings.Index(s, "M") > 0
	c.UpdatesTuesday = strings.Index(s, "Tu") > 0
	c.UpdatesWednesday = strings.Index(s, "W") > 0
	c.UpdatesThursday = strings.Index(s, "Th") > 0
	c.UpdatesFriday = strings.Index(s, "F") > 0
	c.UpdatesSaturday = strings.Index(s, "Sa") > 0
}

func listActiveComics(cfg *config.Config) {
	db, err := cfg.ConnectDB()
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to DB")
	}

	activeComics, err := comics.FetchComics(context.Background(), db, userId, &filterActive, nil)
	if err != nil {
		log.WithError(err).Fatal("Failed to fetch comics list")
	}

	for _, c := range activeComics {
		fmt.Printf("%s\n", c.ToString())
	}

}

func markComicRead(cfg *config.Config) {
	logger := log.WithFields(log.Fields{
		"userId":  userId,
		"comicId": comicId,
	})
	db, err := cfg.ConnectDB()
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to DB")
	}

	if userId == 0 || comicId == 0 {
		flag.Usage()
		logger.Fatalf("Missing parameter user-id or comic-id")
	}

	err = comics.UpdateReadNow(comicId, userId, time.Now(), db)
	if err != nil {
		logger.WithError(err).Fatal("Failed to update read now for comic")
	}
}

var userId int
var comicId int
var filterActive bool

func main() {
	cfg := config.Load()
	//logger := cfg.LogContext.WithField("operation", "main")

	flag.IntVar(&userId, "user-id", 0, "User ID to update")
	flag.IntVar(&comicId, "comic-id", 0, "Comic ID to update")
	flag.BoolVar(&filterActive, "active", true, "Only show active comics")
	flag.Parse()

	switch os.Args[1] {
	case "add":
		addNewComic(&cfg)
	case "list":
		listActiveComics(&cfg)
	case "mark-read":
		markComicRead(&cfg)
	}
}
