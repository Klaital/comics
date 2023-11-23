package postgresstore

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"github.com/klaital/comics/pkg/datalayer"
	"github.com/klaital/comics/pkg/datalayer/postgresstore/queries"
	"log/slog"
)

//go:embed migrations/*
var Migrations embed.FS

type Store struct {
	q *queries.Queries
}

func New(db *sql.DB) *Store {
	return &Store{q: queries.New(db)}
}

func StringPtrToSql(s *string) sql.NullString {
	ns := sql.NullString{
		String: "",
		Valid:  false,
	}
	if s != nil {
		ns.Valid = true
		ns.String = *s
	}
	return ns
}

func (s *Store) AddUser(ctx context.Context, username, email, password string) (*datalayer.User, error) {
	//// Check for existing user in cache
	//for _, u := range s.users {
	//	if u.Email == email {
	//		return nil, datalayer.ErrEmailAlreadyUsed
	//	}
	//}
	//u := datalayer.User{
	//	ID:             uint64(rand.Int63()),
	//	Email:          email,
	//	PasswordDigest: password, // TODO: digest the password here? Require it to be digested before being passed in?
	//}

	// TODO: cache the user's login session

	err := s.q.AddUser(ctx, queries.AddUserParams{
		Username: username,
		Email:    email,
		Passwd:   password,
	})
	if err != nil {
		slog.Error("Failed to add user", "err", err, "email", email)
		return nil, fmt.Errorf("creating user: %w", err)
	}

	// TODO: convert to RETURNING id and populate return user
	return nil, nil
}

func (s *Store) AddComic(ctx context.Context, c *datalayer.Comic) error {
	params := queries.AddComicParams{
		Ordinal: sql.NullInt32{
			Valid: true,
			Int32: int32(c.Ordinal),
		},
		UserID:           int32(c.UserID),
		Title:            c.Title,
		BaseUrl:          sql.NullString{Valid: true, String: c.BaseURL},
		FirstComicUrl:    StringPtrToSql(c.FirstComicUrl),
		LatestComicUrl:   StringPtrToSql(c.LatestComicUrl),
		RssUrl:           StringPtrToSql(c.RssUrl),
		UpdatesMonday:    c.UpdatesMonday,
		UpdatesTuesday:   c.UpdatesTuesday,
		UpdatesWednesday: c.UpdatesWednesday,
		UpdatesThursday:  c.UpdatesThursday,
		UpdatesFriday:    c.UpdatesFriday,
		UpdatesSaturday:  c.UpdatesSaturday,
		UpdatesSunday:    c.UpdatesSunday,
		Active:           c.Active,
		Nsfw:             c.Nsfw,
		LastRead:         sql.NullTime{Valid: false},
	}

	err := s.q.AddComic(ctx, params)
	if err != nil {
		return fmt.Errorf("adding comic: %w", err)
	}

	//Success!
	return nil
}

func (s *Store) GetUser(ctx context.Context, id uint64) (*datalayer.User, error) {
	return nil, datalayer.ErrNotImplemented
}

func (s *Store) UpdateUser(ctx context.Context, userId uint64, newData *datalayer.User) error {
	return datalayer.ErrNotImplemented
}

func (s *Store) DeleteUser(ctx context.Context, userId uint64) error {
	return datalayer.ErrNotImplemented
}

func (s *Store) GetComics(ctx context.Context, userId uint64, filterNsfw *bool, filterActive *bool) ([]datalayer.Comic, error) {
	return nil, datalayer.ErrNotImplemented
}

func (s *Store) GetComic(ctx context.Context, userId, comicId uint64) (*datalayer.Comic, error) {
	return nil, datalayer.ErrNotImplemented
}

func (s *Store) UpdateComic(ctx context.Context, c *datalayer.Comic) error {
	return datalayer.ErrNotImplemented
}

func (s *Store) DeleteComic(ctx context.Context, userId, comicId uint64) error {
	return datalayer.ErrNotImplemented
}

func (s *Store) MarkComicRead(ctx context.Context, userId, comicId uint64) error {
	return datalayer.ErrNotImplemented
}

func (s *Store) AddRssItems(ctx context.Context, userId, comicId uint64, rssItems []datalayer.RssItem) error {
	return datalayer.ErrNotImplemented
}

func (s *Store) MarkRssRead(ctx context.Context, userId uint64, rssItemGuid string) error {
	return datalayer.ErrNotImplemented
}
