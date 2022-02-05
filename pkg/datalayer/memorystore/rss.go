package memorystore

import (
	"context"
	"errors"
	"github.com/klaital/comics/pkg/datalayer"
	"time"
)

func (s *MemoryStore) AddRssItems(ctx context.Context, userId, comicId uint64, rssItems []datalayer.RssItem) error {
	existingComic, ok := s.comics[comicId]
	if !ok {
		return datalayer.ErrNotFound
	}
	if userId != existingComic.UserID {
		return errors.New("user mismatch")
	}

	for i, newItem := range rssItems {
		for _, existingItem := range existingComic.RssItems {
			if newItem.Guid != existingItem.Guid { // Filter out any dupes
				existingComic.RssItems = append(existingComic.RssItems, rssItems[i])
			}
		}
	}

	return nil
}

func (s *MemoryStore) MarkRssRead(ctx context.Context, userId uint64, rssItemGuid string) error {
	// Since we're doing a document store, this has to be done by brute-force
	for _, c := range s.comics {
		if c.UserID != userId {
			continue
		}
		for j, r := range c.RssItems {
			if r.Guid == rssItemGuid {
				tmp := time.Now()
				c.RssItems[j].ReadAt = &tmp
				return nil
			}
		}
	}
	return datalayer.ErrNotFound
}
