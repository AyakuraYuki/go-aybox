package aybox

import (
	"github.com/dromara/carbon/v2"
	"time"
)

type AgeType interface {
	int | int32 | int64 | uint | uint32 | uint64
}

// CalculateRealAgeFromString accepts a date string in layout `2006-01-02`
// to calculate the real age, accept given timezone for the calculation.
func CalculateRealAgeFromString[T AgeType](date string, timezone ...string) (age T) {
	birthday := carbon.ParseByLayout(date, time.DateOnly, timezone...)
	return CalculateRealAge[T](birthday, timezone...)
}

// CalculateRealAge returns the real age from some birthday to now, accept
// given timezone for the calculation.
func CalculateRealAge[T AgeType](birthday carbon.Carbon, timezone ...string) (age T) {
	if len(timezone) > 0 {
		birthday.SetTimezone(timezone[0])
	}
	now := carbon.Now(timezone...)
	if now.Lt(birthday) {
		return 0
	}
	years := now.Year() - birthday.Year()
	if now.Month() < birthday.Month() || (now.Month() == birthday.Month() && now.Day() < birthday.Day()) {
		years--
	}
	return T(years)
}
