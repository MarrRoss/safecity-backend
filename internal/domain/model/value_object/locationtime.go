package value_object

import (
	"errors"
	"time"
)

// TODO: нужна ли структура?
type LocationTime struct {
	Location Location
	Time     time.Time
}

func NewLocationTime(loc Location, t time.Time) (LocationTime, error) {
	if t.IsZero() {
		return LocationTime{}, errors.New("time is required")
	}
	now := time.Now()
	sixMonthsAgo := now.AddDate(0, -24, 0)
	if t.Before(sixMonthsAgo) {
		return LocationTime{}, errors.New("time must be within the last 2 years")
	}

	newLocTime := LocationTime{Location: loc, Time: t}
	return newLocTime, nil
}
