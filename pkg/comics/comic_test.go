package comics

import (
	"database/sql"
	"time"
)

func (suite *ComicTestSuite) TestInsertNewComic() {
	db, err := suite.Cfg.ConnectDB()
	suite.NoError(err, "Failed to connect to test db")

	// Insert a new comic at the end of the list, check that no other ordinals have changed
	newComic := Comic{
		UserID:         1,
		Ordinal:        100,
		Title:          "TestInsertNewComic",
		BaseURL:        "http://TestInsertNewComic",
		FirstComicUrl:  sql.NullString{String: "http://TestInsertNewComic/1"},
		LatestComicUrl: sql.NullString{String: "http://TestInsertNewComic/3"},
		RssUrl:         sql.NullString{String: "http://TestInsertNewComic/rss"},
		UpdatesFriday:  true,
		LastRead:       time.Now(),
	}

	err = InsertNewComic(&newComic, db)
	suite.NoError(err, "Error inserting new comic at the end of the list")
}
