package datalayer

import (
	"context"
	"errors"
)

type ComicDataSource interface {
	// AddUser will register a new user
	AddUser(ctx context.Context, username, email, password string) (*User, error)

	// GetUser will fetch a user's details. Use GetComics to fetch that user's subscriptions.
	GetUser(ctx context.Context, id uint64) (*User, error)

	// UpdateUser will modify a user's data
	UpdateUser(ctx context.Context, userId uint64, newData *User) error

	// DeleteUser will remove a user from the DB.
	// This will cascade to removing all of their comics/RSS data as well.
	DeleteUser(ctx context.Context, userId uint64) error

	// GetComics will retrieve all comics data for the user. If the supplied filters are not nil,
	// then the values are used as filters.
	GetComics(ctx context.Context, userId uint64, filterNsfw *bool, filterActive *bool) ([]Comic, error)

	// GetComic fetches the data on a single comic.
	GetComic(ctx context.Context, userId, comicId uint64) (*Comic, error)

	// AddComic will insert a new record. The supplied Comic will be updated with the generated unique ID.
	// An error will be returned in the event of a DB error, including trying to add a dupe record.
	AddComic(ctx context.Context, c *Comic) error

	// UpdateComic updates the data for the comic itself. Does not update anything about the RSS data.
	UpdateComic(ctx context.Context, c *Comic) error

	// DeleteComic removes a comic from a user's account. Will cascade to removing the associated RSS items.
	DeleteComic(ctx context.Context, userId, comicId uint64) error

	// MarkComicRead is a helper method to  update the "last read" record for a comic to "now".
	// Does not update any RSS items.
	MarkComicRead(ctx context.Context, userId, comicId uint64) error

	// AddRssItems will insert new RSS feed data records.
	// The supplied RSS Items will be updated with their generated unique IDs.
	AddRssItems(ctx context.Context, userId, comicId uint64, rssItems []RssItem) error

	// MarkRssRead updates a single RSS feed item as "read"
	MarkRssRead(ctx context.Context, userId uint64, rssItemGuid string) error
}

var ErrEmailAlreadyUsed error = errors.New("email already in use")
var ErrNotFound error = errors.New("not found")
var ErrNotValid error = errors.New("input data not valid")
var ErrNotImplemented = errors.New("operation not implemented")
