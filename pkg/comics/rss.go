package comics

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"time"
)

type RssItem struct {
	ID         int        `db:"item_id" json:"-"`
	UserID     int        `db:"user_id" json:"-"`
	WebcomicID int        `db:"webcomic_id" json:"webcomic_id"`
	Posted     *time.Time `json:"posted"`
	Guid       string     `db:"guid" json:"guid"`
	IsRead     bool       `db:"is_read" json:"is_read"`
	Title      string     `db:"title" json:"title"`
	Link       string     `db:"link" json:"link"`
}

// FetchRssFeed fetches the current RSS data for a comic
func (c *ComicRecord) FetchRssFeed() ([]RssItem, error) {
	if c.RssUrl == nil {
		return []RssItem{}, nil
	}
	logger := log.WithField("rss_url", c.RssUrl)

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(*c.RssUrl)
	if err != nil {
		logger.WithError(err).Error("Failed to parse the RSS feed")
		return nil, err
	}

	feedItems := make([]RssItem, len(feed.Items))
	for i := range feed.Items {
		feedItems[i] = RssItem{
			UserID:     c.UserID,
			WebcomicID: c.ID,
			Posted:     feed.Items[i].PublishedParsed,
			Guid:       feed.Items[i].GUID,
			IsRead:     false,
			Link:       feed.Items[i].Link,
			Title:      feed.Items[i].Title,
		}
	}

	return feedItems, nil
}

// Insert writes a newly-published RSS item to the DB.
func (item *RssItem) Insert(ctx context.Context, db *sqlx.DB) error {
	sqlStmt := `INSERT INTO rss_items (user_id, webcomic_id, guid, is_read, link, title) VALUES (:user_id, :webcomic_id, :guid, :is_read, :link, :title)`
	_, err := db.NamedExecContext(ctx, sqlStmt, item)
	return err
}

// FetchKnownRssItems retrieves all RSS items from the DB for this user/comic combination
func (c *ComicRecord) FetchKnownRssItems(db *sqlx.DB) ([]RssItem, error) {
	logger := log.WithFields(log.Fields{
		"comic":  c.Title,
		"userId": c.UserID,
	})

	sqlStmt := db.Rebind(`SELECT * FROM rss_items WHERE user_id = ? AND webcomic_id = ?`)
	var items []RssItem
	if err := db.Select(&items, sqlStmt, c.UserID, c.ID); err != nil {
		if err == sql.ErrNoRows {
			return []RssItem{}, nil
		}
		logger.WithError(err).Error("Error fetching rss data from DB")
		return nil, err
	}

	return items, nil
}

// MarkRead updates the DB to note that the Item has been read.
func (item *RssItem) MarkRead(ctx context.Context, db *sqlx.DB) error {
	sqlStmt := db.Rebind(`UPDATE rss_items SET is_read = true WHERE guid = ?`)
	_, err := db.ExecContext(ctx, sqlStmt, item.Guid)
	return err
}
