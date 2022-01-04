package comics

import (
	"time"
)

func (suite *ComicTestSuite) TestInsertNewComic() {
	db, err := suite.Cfg.ConnectDB()
	suite.NoError(err, "Failed to connect to test db")

	// Insert a new comic at the end of the list, check that no other ordinals have changed
	newComic := ComicRecord{
		UserID:         1,
		Ordinal:        100,
		Title:          "TestInsertNewComic",
		BaseURL:        "http://TestInsertNewComic",
		FirstComicUrl:  strPtr("http://TestInsertNewComic/1"),
		LatestComicUrl: strPtr("http://TestInsertNewComic/3"),
		RssUrl:         strPtr("http://TestInsertNewComic/rss"),
		UpdatesFriday:  true,
		LastRead:       time.Now(),
	}

	err = InsertNewComic(&newComic, db)
	suite.NoError(err, "Error inserting new comic at the end of the list")

	expectedOrdinals := map[int]int{ // maps comic id -> ordinal
		1: 1,
		2: 2,
		3: 3,
	}
	for comicId, expectedOrdinal := range expectedOrdinals {
		c, err := FetchComicByID(db, comicId, 1)
		suite.NoError(err, "Error fetching comic for validation")
		suite.Equal(expectedOrdinal, c.Ordinal, "Ordinal should not have changed")
	}

	// Insert a new comic at the middle of the list, check that all larger ordinals have changed
	newComic = ComicRecord{
		UserID:         1,
		Ordinal:        2,
		Title:          "TestInsertNewComic2",
		BaseURL:        "http://TestInsertNewComic2",
		FirstComicUrl:  strPtr("http://TestInsertNewComic2/1"),
		LatestComicUrl: strPtr("http://TestInsertNewComic2/3"),
		RssUrl:         strPtr("http://TestInsertNewComic2/rss"),
		UpdatesFriday:  true,
		LastRead:       time.Now(),
	}

	err = InsertNewComic(&newComic, db)
	suite.NoError(err, "Error inserting new comic at the middle of the list")

	expectedOrdinals = map[int]int{ // maps comic id -> ordinal
		1: 1,
		2: 3,
		3: 4,
	}
	for comicId, expectedOrdinal := range expectedOrdinals {
		c, err := FetchComicByID(db, comicId, 1)
		suite.NoError(err, "Error fetching comic for validation")
		suite.Equal(expectedOrdinal, c.Ordinal, "Ordinal should not have changed")
	}
}

func (suite *ComicTestSuite) TestFetchComicByID() {
	var err error
	var c *ComicRecord

	db, err := suite.Cfg.ConnectDB()
	suite.NoError(err, "Failed to connect to test db")

	// Happy path
	c, err = FetchComicByID(db, 1, 1)
	suite.NoError(err, "Error when fetching first comic")
	suite.Equal(c.Title, "Test Comic 1")
	suite.Equal(c.UserID, 1)

	// User does not own the particular record
	c, err = FetchComicByID(db, 1, 2)
	suite.NoError(err, "Error when fetching first comic")
	suite.Nil(c, "Expected no comic returned with the wrong user is asking")
}
