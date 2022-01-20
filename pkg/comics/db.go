package comics

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"time"
)

func FetchComicByID(db *sqlx.DB, comicID, userID int) (*ComicRecord, error) {
	logger := log.WithFields(log.Fields{
		"comic-id": comicID,
		"user-id":  userID,
	})
	sqlQuery := db.Rebind(`SELECT * FROM webcomic WHERE webcomic_id=? AND user_id=?`)
	var c ComicRecord
	err := db.Get(&c, sqlQuery, comicID, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		logger.WithError(err).Error("Error querying webcomic by ID")
		return nil, err
	}

	rssSqlQuery := db.Rebind(`SELECT item_id, webcomic_id, user_id, is_read, guid, title, link FROM rss_items WHERE webcomic_id = ? AND user_id = ?`)
	var rssItems []RssItem
	err = db.Select(&rssItems, rssSqlQuery, comicID, userID)
	if err == sql.ErrNoRows {
		c.RssItems = []RssItem{}
	} else if err != nil {
		logger.WithError(err).Error("Error querying RSS for webcomic by ID")
		return nil, err
	}
	c.RssItems = rssItems

	return &c, err
}

type fetchComicsRow struct {
	ComicRecord
	RssItem
}

// FetchComics queries the DB for comics matching the given filters.
// UserID is mandatory, but filtering on Active or NSFW flags is optional.
func FetchComics(ctx context.Context, db *sqlx.DB, userId int, filterActive *bool, filterNsfw *bool) ([]ComicRecord, error) {
	logger := log.WithFields(log.Fields{
		"userId": userId,
	})

	sqlStmt := `SELECT webcomic.webcomic_id, webcomic.user_id, webcomic.title, 
			webcomic.base_url, webcomic.first_comic_url, webcomic.latest_comic_url,
			webcomic.updates_monday, webcomic.updates_tuesday, webcomic.updates_wednesday,
			webcomic.updates_thursday, webcomic.updates_friday, webcomic.updates_saturday,
			webcomic.updates_sunday, webcomic.ordinal, webcomic.last_read, 
			webcomic.active, webcomic.nsfw,
			COALESCE(rss_items.guid, '') "guid", COALESCE(rss_items.title, '') "title", COALESCE(rss_items.link, '') "link"
		FROM webcomic LEFT OUTER JOIN rss_items
			ON webcomic.webcomic_id = rss_items.webcomic_id
		WHERE webcomic.user_id = ?`
	params := []interface{}{userId}
	if filterNsfw != nil {
		logger = logger.WithField("filterNsfw", *filterNsfw)
		sqlStmt += ` AND webcomic.nsfw=?`
		params = append(params, *filterNsfw)
	}
	if filterActive != nil {
		logger = logger.WithField("filterActive", *filterActive)
		sqlStmt += ` AND webcomic.active=?`
		params = append(params, *filterActive)
	}

	log.Debug("Composed filtered comics query")
	var rows []fetchComicsRow
	if err := db.SelectContext(ctx, &rows, db.Rebind(sqlStmt), params...); err != nil {
		logger.WithError(err).Error("Error fetching comics + RSS data")
		return nil, err
	}

	// Collate the RSS Items onto their base ComicRecord
	comicSet := make(map[int]ComicRecord, 0)
	for i, row := range rows {
		existingComic, ok := comicSet[row.WebcomicID]
		if !ok {
			existingComic = rows[i].ComicRecord
		}
		if row.Guid != "" {
			existingComic.RssItems = append(existingComic.RssItems, row.RssItem)
		}
		comicSet[row.WebcomicID] = existingComic
	}

	// Reformat the comic data into a flat array for return
	returnSet := make([]ComicRecord, 0, len(comicSet))
	for id := range comicSet {
		returnSet = append(returnSet, comicSet[id])
	}
	return returnSet, nil
}

//func FetchActiveComics(db *sqlx.DB, userID int) ([]ComicRecord, error) {
//	var comicData []ComicRecord
//	sqlQuery := `SELECT webcomic_id, user_id, title,
//			base_url, first_comic_url, latest_comic_url,
//			updates_monday, updates_tuesday, updates_wednesday,
//			updates_thursday, updates_friday, updates_saturday,
//			updates_sunday, ordinal, last_read, active, nsfw
//		FROM webcomic
//		WHERE active=1 AND user_id = ?
//		ORDER BY ordinal ASC`
//
//	sqlStmt, err := db.Preparex(sqlQuery)
//	if err != nil {
//		log.WithError(err).Errorf("Failed to prepare fetch query")
//		return nil, err
//	}
//
//	// Fetch the data!
//	err = sqlStmt.Select(&comicData, userID)
//	if err != nil {
//		log.WithError(err).Errorf("Failed to select comics data")
//		return nil, err
//	}
//
//	return comicData, nil
//}

func UpdateReadNow(comicId, userId int, readTime time.Time, db *sqlx.DB) error {
	logger := log.WithFields(log.Fields{
		"operation":   "UpdateReadNow",
		"read_time":   readTime.Unix(),
		"webcomic_id": comicId,
	})
	sqlQuery := db.Rebind(`UPDATE webcomic SET last_read = ? WHERE webcomic_id = ? AND user_id = ?`)
	_, err := db.Exec(sqlQuery, readTime, comicId, userId)
	if err != nil {
		logger.WithError(err).Errorf("Failed to update comics data")
		return err
	}

	sqlQuery = db.Rebind(`UPDATE rss_items SET is_read = true WHERE webcomic_id = ? AND user_id = ?`)
	_, err = db.Exec(sqlQuery, comicId, userId)
	if err != nil {
		logger.WithError(err).Errorf("Failed to update comics data")
		return err
	}

	// Success!
	return nil
}

func MoveComicsDownIfCollision(insertOrdinal, userId int, db *sqlx.DB) error {
	logger := log.WithFields(log.Fields{
		"operation":     "MoveComicsDownIfCollision",
		"insertOrdinal": insertOrdinal,
	})

	// Check if a reassignment is needed
	sqlStmt := db.Rebind(`SELECT COUNT(webcomic_id) FROM webcomic WHERE ordinal = ?`)
	var conflicts int64
	if err := db.Get(&conflicts, sqlStmt, insertOrdinal); err != nil {
		logger.WithError(err).Error("Error checking for ordinal conflicts")
		return err
	}

	if conflicts > 0 {
		sqlStmt = db.Rebind(`UPDATE webcomic SET ordinal = ordinal + 1 WHERE user_id = ? AND ordinal >= ?`)
		res, err := db.Exec(sqlStmt, userId, insertOrdinal)
		if err != nil {
			logger.WithError(err).Error("Error reassigning ordinals")
			return err
		}
		updateCount, err := res.RowsAffected()
		if err != nil {
			logger.WithError(err).Error("Error fetching update count")
			return err
		}
		if updateCount != conflicts {
			err = errors.New("update count does not match expected conflicts")
			logger.WithError(err).WithFields(log.Fields{
				"expectedConflicts": conflicts,
				"updateCount":       updateCount,
			}).Warn("Ordinal update count does not match counted conflicts")
		}
	}

	// Success! (or no-op if the count was zero)
	return nil

}

func InsertNewComic(c *ComicRecord, db *sqlx.DB) error {
	logger := log.WithFields(log.Fields{
		"operation":   "InsertNewComic",
		"comic_title": c.Title,
		"ordinal":     c.Ordinal,
	})

	// First find out if the ordinal is already in use. If so, move it and anything below it all down one step.
	if err := MoveComicsDownIfCollision(c.Ordinal, c.UserID, db); err != nil {
		logger.WithError(err).Error("Failed to update colliding ordinals")
		return err
	}
	sqlQuery := db.Rebind(`INSERT INTO webcomic (
							 title, 
							 user_id,
							 base_url, first_comic_url, latest_comic_url, rss_url, 
							 updates_monday, updates_tuesday, updates_wednesday,
							 updates_thursday, updates_friday, 
							 updates_saturday, updates_sunday,
							 ordinal, last_read
						 ) VALUES (
                             :title, 
							 :user_id,
                             :base_url, :first_comic_url, :latest_comic_url, :rss_url, 
                             :updates_monday, :updates_tuesday, :updates_wednesday,
                             :updates_thursday, :updates_friday, 
                             :updates_saturday, :updates_sunday,
                             :ordinal, :last_read           
						 )`)
	res, err := db.NamedExec(sqlQuery, c)
	if err != nil {
		logger.WithError(err).Error("Failed to insert new comic")
		return err
	}
	idTmp, err := res.LastInsertId()
	if err != nil {
		logger.WithError(err).Error("Failed to get ID created for new comic")
	}
	c.ID = int(idTmp)
	logger = logger.WithField("comic_ID", c.ID)
	logger.Debug("Successfully inserted new comic")
	return nil
}
