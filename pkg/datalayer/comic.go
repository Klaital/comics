package datalayer

import "time"

// Comic contains the data to be saved/retrieved from the database.
// Transform this as needed for API interface types.
type Comic struct {
	ID               uint64    `db:"webcomic_id"`
	UserID           uint64    `db:"user_id"`
	Title            string    `db:"title"`
	BaseURL          string    `db:"base_url"`
	FirstComicUrl    *string   `db:"first_comic_url"`
	LatestComicUrl   *string   `db:"latest_comic_url"`
	RssUrl           *string   `db:"rss_url"`
	UpdatesMonday    bool      `db:"updates_monday"`
	UpdatesTuesday   bool      `db:"updates_tuesday"`
	UpdatesWednesday bool      `db:"updates_wednesday"`
	UpdatesThursday  bool      `db:"updates_thursday"`
	UpdatesFriday    bool      `db:"updates_friday"`
	UpdatesSaturday  bool      `db:"updates_saturday"`
	UpdatesSunday    bool      `db:"updates_sunday"`
	Ordinal          int       `db:"ordinal"`
	LastRead         time.Time `db:"last_read"`
	Active           bool      `db:"active"`
	Nsfw             bool      `db:"nsfw"`

	RssItems []RssItem
}
