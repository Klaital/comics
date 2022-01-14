package comics

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"time"
)

func FetchComicByID(db *sqlx.DB, comicID, userID int) (*ComicRecord, error) {
	sqlQuery := db.Rebind(`SELECT * FROM webcomic WHERE webcomic_id=? AND user_id=?`)
	var c ComicRecord
	err := db.Get(&c, sqlQuery, comicID, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &c, err
}
func FetchActiveComics(db *sqlx.DB, userID int) ([]ComicRecord, error) {
	var comicData []ComicRecord
	sqlQuery := `SELECT webcomic_id, user_id, title, 
			base_url, first_comic_url, latest_comic_url,
			updates_monday, updates_tuesday, updates_wednesday,
			updates_thursday, updates_friday, updates_saturday,
			updates_sunday, ordinal, last_read
		FROM webcomic
		WHERE active=1 AND user_id = ?
		ORDER BY ordinal ASC`

	sqlStmt, err := db.Preparex(sqlQuery)
	if err != nil {
		log.WithError(err).Errorf("Failed to prepare fetch query")
		return nil, err
	}

	// Fetch the data!
	err = sqlStmt.Select(&comicData, userID)
	if err != nil {
		log.WithError(err).Errorf("Failed to select comics data")
		return nil, err
	}

	return comicData, nil
}

func UpdateReadNow(comicId int, readTime time.Time, db *sqlx.DB) error {
	logger := log.WithFields(log.Fields{
		"operation":   "UpdateReadNow",
		"read_time":   readTime.Unix(),
		"webcomic_id": comicId,
	})
	sqlQuery := db.Rebind(`UPDATE webcomic SET last_read = ? WHERE webcomic_id = ?`)

	// Fetch the data!
	_, err := db.Exec(sqlQuery, readTime, comicId)
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
	sqlStmt := db.Rebind(`SELECT COUNT(webcomic_id) FROM webcomic WHERE user_id = ? AND ordinal = ?`)
	var conflicts int64
	if err := db.Get(&conflicts, sqlStmt, userId, insertOrdinal); err != nil {
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
