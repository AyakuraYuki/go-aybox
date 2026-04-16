package calendar

import "time"

type Gregorian struct {
	Time  time.Time
	Error error
}

func (g Gregorian) String() string {
	if g.Time.IsZero() {
		return ""
	}
	return g.Time.String()
}
