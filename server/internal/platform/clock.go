package platform

import (
	"time"
	_ "time/tzdata" // embed zone data so Asia/Kathmandu loads in minimal containers
)

const AppTimezone = "Asia/Kathmandu"

func LoadAppLocation() (*time.Location, error) {
	return time.LoadLocation(AppTimezone)
}

// DayBounds returns the start (inclusive) and end (exclusive) of the local day containing t.
func DayBounds(t time.Time, loc *time.Location) (time.Time, time.Time) {
	local := t.In(loc)
	start := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc)
	return start, start.AddDate(0, 0, 1)
}
