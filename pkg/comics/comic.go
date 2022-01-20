package comics

import (
	"fmt"
	"github.com/klaital/intmath"
	log "github.com/sirupsen/logrus"
	"time"
)

type ComicRecord struct {
	ID               int       `db:"webcomic_id"`
	UserID           int       `db:"user_id"`
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

func strPtr(s string) *string {
	if len(s) == 0 {
		return nil
	}

	return &s
}

// IsValid checks whether the comic struct is prima facia valid for insertion into the database.
//Does not check against unique indices, as might be on the ID or ordinal fields.
func (c ComicRecord) IsValid() error {
	if len(c.Title) == 0 {
		return fmt.Errorf("missing required field: title")
	}
	if len(c.BaseURL) == 0 {
		return fmt.Errorf("missing required field: base_url")
	}
	if c.Ordinal == 0 {
		return fmt.Errorf("missing required field: ordinal")
	}
	// Success!
	return nil
}

func (c ComicRecord) SupportedDayCodes() string {
	s := ""
	if c.UpdatesSunday {
		s += "Su"
	}
	if c.UpdatesMonday {
		s += "M"
	}
	if c.UpdatesTuesday {
		s += "Tu"
	}
	if c.UpdatesWednesday {
		s += "W"
	}
	if c.UpdatesThursday {
		s += "Th"
	}
	if c.UpdatesFriday {
		s += "F"
	}
	if c.UpdatesSaturday {
		s += "Sa"
	}
	return s
}

func startOfToday() time.Time {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Fatal(err)
	}
	y, m, d := time.Now().In(loc).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, loc)
}
func GetTodaySelector() func(ComicRecord) bool {
	switch startOfToday().Weekday() {
	case time.Monday:
		return func(c ComicRecord) bool { return c.UpdatesMonday && c.LastRead.Unix() < startOfToday().Unix() }
	case time.Tuesday:
		return func(c ComicRecord) bool { return c.UpdatesTuesday && c.LastRead.Unix() < startOfToday().Unix() }
	case time.Wednesday:
		return func(c ComicRecord) bool { return c.UpdatesWednesday && c.LastRead.Unix() < startOfToday().Unix() }
	case time.Thursday:
		return func(c ComicRecord) bool { return c.UpdatesThursday && c.LastRead.Unix() < startOfToday().Unix() }
	case time.Friday:
		return func(c ComicRecord) bool { return c.UpdatesFriday && c.LastRead.Unix() < startOfToday().Unix() }
	case time.Saturday:
		return func(c ComicRecord) bool { return c.UpdatesSaturday && c.LastRead.Unix() < startOfToday().Unix() }
	case time.Sunday:
		return func(c ComicRecord) bool { return c.UpdatesSunday && c.LastRead.Unix() < startOfToday().Unix() }
	}
	return func(c ComicRecord) bool { return c.LastRead.Unix() < startOfToday().Unix() }
}

func daysAgo(t time.Time) string {
	d := time.Since(t)
	if d.Hours() > 24 {
		if d.Hours() > 48 {
			return fmt.Sprintf("%d days ago", intmath.RoundToInt(d.Hours()/24))
		} else {
			return "yesterday"
		}
	}
	if d.Hours() > 1 {
		return fmt.Sprintf("%d hours ago", intmath.RoundToInt(d.Hours()))
	}
	if d.Seconds() < 60 {
		return "just now"
	}
	return fmt.Sprintf("%d minutes ago", intmath.RoundToInt(d.Minutes()))
}
func (c ComicRecord) ToString() string {
	return fmt.Sprintf("[%11s] (%s) %s", c.SupportedDayCodes(), daysAgo(c.LastRead), c.BaseURL)
}
func (c ComicRecord) DaysAgoNow() string {
	return daysAgo(c.LastRead)
}

func SelectSubset(set []ComicRecord, selector func(ComicRecord) bool) (selected []ComicRecord, theRest []ComicRecord) {
	selected = make([]ComicRecord, 0)
	theRest = make([]ComicRecord, 0, len(set))
	for _, c := range set {
		if selector(c) {
			selected = append(selected, c)
		} else {
			theRest = append(theRest, c)
		}
	}
	return selected, theRest
}
func SelectMapSubset(set map[int]ComicRecord, selector func(ComicRecord) bool) (selected map[int]ComicRecord, theRest map[int]ComicRecord) {
	selected = make(map[int]ComicRecord, 0)
	theRest = make(map[int]ComicRecord, len(set))
	for i, c := range set {
		if selector(c) {
			selected[c.Ordinal] = set[i]
		} else {
			theRest[c.Ordinal] = set[i]
		}
	}
	return selected, theRest
}
func SortByOrdinal(set []ComicRecord) map[int]ComicRecord {
	retval := make(map[int]ComicRecord, len(set))
	for i, c := range set {
		retval[c.Ordinal] = set[i]
	}
	return retval
}
