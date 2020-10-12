package core

import (
	"time"
)

// DateTime ...
type DateTime struct {
}

// NewDateTime creates a new instance of the date time
func NewDateTime() DateTime {
	return DateTime{}
}

// Now returns the current date and time in UTC
func (dt DateTime) Now() time.Time {
	return time.Now().UTC()
}

// AddDuration will add duration to now date in utc
func (dt DateTime) AddDuration(duration time.Duration) time.Time {
	date := dt.Now().Add(duration)
	return date
}
