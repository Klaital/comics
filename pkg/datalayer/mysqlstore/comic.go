package mysqlstore

import (
	"context"
	"errors"
	"github.com/klaital/comics/pkg/datalayer"
	"time"
)

func (s *MysqlStore) GetComics(ctx context.Context, userId uint64, filterNsfw *bool, filterActive *bool) ([]datalayer.Comic, error) {
	// Fetch all comics for the user
	comics := make([]datalayer.Comic, 0)
	for id, c := range s.comics {
		if c.UserID != userId {
			continue
		}

		// Apply filters
		if filterActive != nil {
			if c.Active != *filterActive {
				continue
			}
		}
		if filterNsfw != nil {
			if c.Nsfw != *filterNsfw {
				continue
			}
		}
		comics = append(comics, s.comics[id])
	}

	return comics, nil
}

func (s *MysqlStore) GetComic(ctx context.Context, userId, comicId uint64) (*datalayer.Comic, error) {
	c, ok := s.comics[comicId]
	if !ok {
		return nil, datalayer.ErrNotFound
	}
	if c.UserID != userId {
		return nil, datalayer.ErrNotFound
	}
	return &c, nil
}

func (s *MysqlStore) AddComic(ctx context.Context, c *datalayer.Comic) error {
	s.comics[c.ID] = *c
	return nil
}

func (s *MysqlStore) UpdateComic(ctx context.Context, c *datalayer.Comic) error {
	existingComic, ok := s.comics[c.ID]
	if !ok {
		return datalayer.ErrNotFound
	}
	if existingComic.UserID != c.UserID {
		return errors.New("user mismatch")
	}
	s.comics[c.ID] = *c
	return nil
}

func (s *MysqlStore) DeleteComic(ctx context.Context, userId, comicId uint64) error {
	existingComic, ok := s.comics[comicId]
	if !ok {
		return datalayer.ErrNotFound
	}
	if userId != existingComic.UserID {
		return errors.New("user mismatch")
	}
	delete(s.comics, comicId)
	return nil
}
func (s *MysqlStore) MarkComicRead(ctx context.Context, userId, comicId uint64) error {
	existingComic, ok := s.comics[comicId]
	if !ok {
		return datalayer.ErrNotFound
	}
	if userId != existingComic.UserID {
		return errors.New("user mismatch")
	}

	existingComic.LastRead = time.Now()
	s.comics[comicId] = existingComic
	return nil
}
