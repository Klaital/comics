package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/klaital/comics/pkg/comics"
	"github.com/klaital/comics/pkg/config"
)




func main() {
	cfg := config.Load()
	logger := cfg.LogContext.WithField("operation", "main")
	db, err := cfg.ConnectDB()
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to db")
	}
	comicData, err := comics.FetchActiveComics(db)

	selector := comics.GetTodaySelector()
	today, theRest := comics.SelectSubset(comicData, selector)
	fmt.Printf("---- Today's Comics:\n")
	for _, c := range today {
		fmt.Printf("%s\n", c.ToString())
	}
	fmt.Printf("\n---- The Rest:\n")
	for _, c := range theRest {
		fmt.Printf("%s\n", c.ToString())
	}

	launchServer(&cfg)
}
