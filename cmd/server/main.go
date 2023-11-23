package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/klaital/comics/pkg/comicserver"
	"github.com/klaital/comics/pkg/config"
	"github.com/klaital/comics/pkg/datalayer"
	"github.com/klaital/comics/pkg/datalayer/postgresstore"
)

func main() {
	cfg := config.Load()

	// Connect to database
	db, err := cfg.ConnectPostgres()
	if err != nil {
		cfg.LogContext.WithError(err).Fatal("failed to connect to DB")
	}

	// TODO: Migrate DB schema

	// Init storage layer
	var storer datalayer.ComicDataSource
	storer = postgresstore.New(db)

	// Init business layer
	comicSrv := comicserver.New(storer, cfg.LogContext)
	
	// Start servers

	launchServer(&cfg)
}
