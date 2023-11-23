package mysqlstore

import (
	"context"
	"errors"
	"fmt"
	"github.com/klaital/comics/pkg/datalayer"
	"time"
)

func (s *MysqlStore) GetComics(ctx context.Context, userId uint64, filterNsfw *bool, filterActive *bool) ([]datalayer.Comic, error) {
	// Fetch all comics for the user
	var comics []datalayer.Comic
	var err error
	err = s.db.Select(&comics, "SELECT * FROM webcomic")
	if err != nil {
		return nil, fmt.Errorf("fetching old comics: %w", err)
	}
	// Success!
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
