package comics

import (
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"time"
)

func FetchActiveComics(db *sqlx.DB) ([]Comic, error) {
	var comicData []Comic
	sqlQuery := `SELECT webcomic_id, title, 
			base_url, first_comic_url, latest_comic_url,
			updates_monday, updates_tuesday, updates_wednesday,
			updates_thursday, updates_friday, updates_saturday,
			updates_sunday, ordinal, last_read
		FROM webcomic
		WHERE active=1
		ORDER BY ordinal ASC`

	sqlStmt, err := db.Preparex(sqlQuery)
	if err != nil {
		log.WithError(err).Errorf("Failed to prepare fetch query")
		return nil, err
	}

	// Fetch the data!
	err = sqlStmt.Select(&comicData)
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

func InsertNewComic(c *Comic, db *sqlx.DB) error {
	logger := log.WithFields(log.Fields{
		"operation":   "InsertNewComic",
		"comic_title": c.Title,
		"ordinal":     c.Ordinal,
	})

	// First find out if the ordinal is already in use. If so, move it and anything below it all down one step.
	sql := db.Rebind(`SELECT webcomic_id, ordinal FROM webcomic WHERE ordinal >= ?`)
	var conflictingComics []Comic
	err := db.Select(&conflictingComics, sql, c.Ordinal)
	if err != nil {
		logger.WithError(err).Error("Error while checking for conflicting ordinals")
		return err
	}
	if len(conflictingComics) > 0 {
		// Increment the ordinals for any conflicts
		sql = db.Rebind(`UPDATE webcomic SET ordinal = ordinal + 1 WHERE ordinal >= ?`)
		res, err := db.Exec(sql, c.Ordinal)
		if err != nil {
			logger.WithError(err).Error("Error while updating conflicting ordinals")
			return err
		}
		if updateCount, err := res.RowsAffected(); err != nil {
			logger.WithError(err).Error("Error checking row count")
		} else {
			logger.WithField("rows_updated", updateCount).Debug("Updated conflicting ordinals")
		}
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
