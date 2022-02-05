package datalayer

import "time"

type RssItem struct {
	UserID      uint64
	ComicID     uint64
	Guid        string
	ReadAt      *time.Time
	Link        string
	Description string
}
