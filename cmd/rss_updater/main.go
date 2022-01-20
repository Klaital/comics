package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/klaital/comics/pkg/comics"
	"github.com/klaital/comics/pkg/config"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	var userId int
	var comicId int
	flag.IntVar(&userId, "user-id", 0, "ID of user to update")
	flag.IntVar(&comicId, "comic-id", 0, "ID of comic to update")
	flag.Parse()
	if userId <= 0 || comicId <= 0 {
		log.Fatal("Missing required parameters user-id and comic-id")
	}

	cfg := config.Load()
	logger := cfg.LogContext.WithFields(log.Fields{
		"user-id":  userId,
		"comic-id": comicId,
	})
	db, err := cfg.ConnectDB()
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to DB")
	}

	comic, err := comics.FetchComicByID(db, comicId, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Infof("Bad input. Comic ID %d does not exist under user %d", comicId, userId)
		}
		logger.WithError(err).Fatal("Failed to query comic data.")
	}

	logger = logger.WithField("comic", comic.Title)

	rssFeed, err := comic.FetchRssFeed()
	if err != nil {
		logger.WithError(err).Fatal("Failed to fetch rss feed")
	}
	var latest *comics.RssItem
	if len(rssFeed) > 0 {
		latest = &rssFeed[0]
		for i, c := range rssFeed {
			if c.Posted == nil {
				continue
			}
			if c.Posted.Unix() > latest.Posted.Unix() {
				latest = &rssFeed[i]
			}
		}
	}
	logger.WithField("item_count", len(rssFeed)).Debug("Fetched current RSS feed")
	if latest != nil {
		logger.WithFields(log.Fields{
			"Posted": latest.Posted.Format(time.RFC3339),
			"IsRead": latest.IsRead,
		}).Info("Latest Post")
	}

	knownItems, err := comic.FetchKnownRssItems(db)
	if err != nil {
		logger.WithError(err).Fatal("Failed to fetch existing RSS data from the DB")
	}
	logger.WithField("existing_data_count", len(knownItems)).Info("Fetched existing RSS data")

	// Convert the feed into a map using the GUID as a key for reconciling the two lists
	allRssData := make(map[string]comics.RssItem, len(rssFeed))
	for i := range knownItems {
		allRssData[knownItems[i].Guid] = knownItems[i]
	}

	// Merge the existing data with the new data fetched from the site's RSS feed
	for i := range rssFeed {
		if _, ok := allRssData[rssFeed[i].Guid]; !ok {
			allRssData[rssFeed[i].Guid] = rssFeed[i]
		}
	}

	// Save a list of the unread items
	//unreadItems := make([]comics.RssItem, 0)
	var unreadItemCount uint
	var newItemCout uint

	// Re-save the new rows
	for _, item := range allRssData {
		if !item.IsRead {
			//unreadItems = append(unreadItems, allRssData[guid])
			unreadItemCount++
		}
		if item.ID == 0 {
			err = item.Insert(context.Background(), db)
			if err != nil {
				logger.WithError(err).Fatal("Error inserting RSS item")
			}
			fmt.Print(".")
			newItemCout++
		}
	}
	fmt.Print("\n")

	log.WithFields(log.Fields{
		"undread_count":  unreadItemCount,
		"new_item_count": newItemCout,
	}).Info("Saved new items")

}
