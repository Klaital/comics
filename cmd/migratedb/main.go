package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/klaital/comics/pkg/datalayer"
	"github.com/klaital/comics/pkg/datalayer/mysqlstore"
	"github.com/klaital/comics/pkg/datalayer/postgresstore"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
)

func main() {
	ctx := context.Background()
	// Connect to both databases
	olduser := os.Getenv("MYSQL_USER")
	oldpass := os.Getenv("MYSQL_PASS")
	oldStore, err := mysqlstore.New("mysql.abandonedfactory.net", olduser, oldpass, "webcomics", 3306)
	if err != nil {
		fmt.Printf("Failed to connect to old db: %v", err)
		os.Exit(1)
	}
	newUser := os.Getenv("PGUSER")
	newPass := os.Getenv("PGPASS")
	newConn, err := sql.Open("postgres", fmt.Sprintf("postgresql://%s:%s@localhost:5433/webcomics?sslmode=disable", newUser, newPass))
	if err != nil {
		fmt.Printf("Failed to connect to new db: %v", err)
		os.Exit(1)
	}
	newStore := postgresstore.New(newConn)

	// Load all comics from the old db
	rootUser, err := oldStore.GetUser(ctx, 1)
	if err != nil {
		fmt.Printf("Failed to fetch root user data: %v", err)
		os.Exit(1)
	}
	comics, err := oldStore.GetComics(ctx, rootUser.ID, nil, nil)
	if err != nil {
		fmt.Printf("Failed to fetch root user's comic data: %v", err)
		os.Exit(1)
	}

	// Insert the comics data into the new DB
	newStore.AddUser(ctx, rootUser.Name, rootUser.Email, rootUser.PasswordDigest)
	for _, c := range comics {
		err = newStore.AddComic(ctx, &datalayer.Comic{
			UserID:           c.UserID,
			Title:            c.Title,
			BaseURL:          c.BaseURL,
			FirstComicUrl:    c.FirstComicUrl,
			LatestComicUrl:   c.LatestComicUrl,
			RssUrl:           c.RssUrl,
			UpdatesMonday:    c.UpdatesMonday,
			UpdatesTuesday:   c.UpdatesTuesday,
			UpdatesWednesday: c.UpdatesWednesday,
			UpdatesThursday:  c.UpdatesThursday,
			UpdatesFriday:    c.UpdatesFriday,
			UpdatesSaturday:  c.UpdatesSaturday,
			UpdatesSunday:    c.UpdatesSunday,
			Ordinal:          c.Ordinal,
			LastRead:         c.LastRead,
			Active:           c.Active,
			Nsfw:             c.Nsfw,
			RssItems:         c.RssItems,
		})
		if err != nil {
			slog.Error("Failed to add comic to new DB", "err", err, "comic", c.Title)
		} else {
			fmt.Printf("#%d: %s\t%s\t%s\n", c.Ordinal, c.Title, c.BaseURL, c.RssUrl)
		}
	}
}
