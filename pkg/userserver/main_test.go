package userserver

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/klaital/comics/pkg/config"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// TODO: launch test DB with docker
	// Connect to the test database
	cfg := config.Load()
	if cfg.Realm != "autotest" {
		cfg.LogContext.WithError(fmt.Errorf("invalid test realm: %s", cfg.Realm)).Fatal("Cannot run DB migrations against realms other than 'autotest'")
	}
	db, err := cfg.ConnectDB()
	if err != nil {
		cfg.LogContext.WithError(err).Fatal("Failed to connect to test DB")
	}

	// Run migrations
	driver, err := mysql.WithInstance(db.DB, &mysql.Config{})
	if err != nil {
		cfg.LogContext.WithError(err).Fatal("Failed to construct migration driver")
	}
	migrator, err := migrate.NewWithDatabaseInstance("file:///home/kit/devel/comics/db/migrations", cfg.DbName, driver)
	if err != nil {
		cfg.LogContext.WithError(err).Fatal("Failed to construct migrator instance")
	}
	err = migrator.Up()
	if err != nil {
		cfg.LogContext.WithError(err).Fatal("Failed to run migrations")
	}

	// TODO: Set up test data
	// Run Tests
	os.Exit(m.Run())
}
