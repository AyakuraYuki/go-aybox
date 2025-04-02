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

// IsLongYear reports whether is a long year, see https://en.wikipedia.org/wiki/ISO_8601#Week_dates.
func IsLongYear(date time.Time) bool {
	_, w := time.Date(date.Year(), 12, 31, 0, 0, 0, 0, date.Location()).ISOWeek()
	return w == weeksPerLongYear
}

// Compare compares two by an operator. Available operators are:
//
//	==	equal
//	!=	not equal
//	<>	not equal
//	> 	left greater than right
//	>=	left greater than or equals to right
//	< 	left less than right
//	<=	left less than or equals to right
func Compare(left time.Time, operator string, right time.Time) bool {
	switch operator {
	case "==":
		return Eq(left, right)
	case "<>", "!=":
		return Ne(left, right)
	case ">":
		return Gt(left, right)
	case ">=":
		return Gte(left, right)
	case "<":
		return Lt(left, right)
	case "<=":
		return Lte(left, right)
	}
	return false
}

// Eq reports whether equal.
func Eq(left, right time.Time) bool { return left.Equal(right) }

// Ne reports whether not equal.
func Ne(left, right time.Time) bool { return !Eq(left, right) }

// Gt reports whether greater than.
func Gt(left, right time.Time) bool { return left.After(right) }

// Gte reports whether greater than or equal.
func Gte(left, right time.Time) bool { return Gt(left, right) || Eq(left, right) }

// Lt reports whether less than.
func Lt(left, right time.Time) bool { return left.Before(right) }

// Lte reports whether less than or equal.
func Lte(left, right time.Time) bool { return Lt(left, right) || Eq(left, right) }

// Between reports whether between two times, excluded the start and end time.
func Between(t, start, end time.Time) bool {
	if Gt(start, end) {
		return false
	}
	return Gt(t, start) && Lt(t, end)
}

// BetweenIncludedStart reports whether between two times, included the start time.
func BetweenIncludedStart(t, start, end time.Time) bool {
	if Gt(start, end) {
		return false
	}
	return Gte(t, start) && Lt(t, end)
}

// BetweenIncludedEnd reports whether between two times, included the end time.
func BetweenIncludedEnd(t, start, end time.Time) bool {
	if Gt(start, end) {
		return false
	}
	return Gt(t, start) && Lte(t, end)
}

// BetweenIncludedBoth reports whether between two times, included the start and end time.
func BetweenIncludedBoth(t, start, end time.Time) bool {
	if Gt(start, end) {
		return false
	}
	return Gte(t, start) && Lte(t, end)
}
