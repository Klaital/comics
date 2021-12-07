package comics

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/klaital/comics/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ComicTestSuite struct {
	suite.Suite
	fixtures *testfixtures.Loader
	Cfg      config.Config
}

// Runs before all tests
func (suite *ComicTestSuite) SetupSuite() {
	suite.Cfg = config.Load()
	db, err := suite.Cfg.ConnectDB()
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to test DB")
	}
	// Run db migrations
	driver, err := mysql.WithInstance(db.DB, &mysql.Config{})
	if err != nil {
		log.WithError(err).Fatal("Failed to construct mysql migration driver")
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://../../db/migrations/",
		"mysql",
		driver)
	if err != nil {
		log.WithError(err).Fatal("Failed to construct DB migrator")
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.WithError(err).Fatal("Failed to migrate database")
	}

	// Load fixtures
	suite.fixtures, err = testfixtures.New(
		testfixtures.Database(db.DB),
		testfixtures.Dialect("mariadb"),
		testfixtures.Directory("testdata/fixtures"),
	)
	if err != nil {
		log.WithError(err).Fatal("Failed to construct fixture loader")
	}
	err = suite.fixtures.Load()
	if err != nil {
		log.WithError(err).Fatal("Failed to load initial fixtures")
	}

}

// Runs before each test
func (suite *ComicTestSuite) SetupTest() {
	// wipe the DB and reload fixtures
	if err := suite.fixtures.Load(); err != nil {
		log.WithError(err).Fatal("Failed to load test fixtures")
	}
}

// Run the suite
func TestComicTestSuite(t *testing.T) {
	suite.Run(t, new(ComicTestSuite))
}
