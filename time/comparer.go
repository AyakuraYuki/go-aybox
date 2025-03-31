package time

import "time"

func IsAM(t time.Time) bool {
	return t.Format("pm") == "am"
}

func IsPM(t time.Time) bool {
	return t.Format("pm") == "pm"
}

// IsLeapYear reports whether is a leap year
func IsLeapYear(date time.Time) bool {
	year := date.Year()
	if year%400 == 0 || (year%4 == 0 && year%100 != 0) {
		return true
	}
	return false
}
