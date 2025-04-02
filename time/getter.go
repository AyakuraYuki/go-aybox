package time

import "time"

func DaysInYear(t time.Time) int {
	if IsLeapYear(t) {
		return DaysPerLeapYear
	}
	return DaysPerNormalYear
}

func DaysInMonth(t time.Time) int {
	return EndOfMonth(t).Day()
}

func MonthOfYear(t time.Time) int {
	return int(t.Month())
}

func DayOfYear(t time.Time) int {
	return t.YearDay()
}

func DayOfMonth(t time.Time) int {
	return t.Day()
}

// DayOfWeek gets day of week like 6
func DayOfWeek(t time.Time) int {
	day := t.Weekday()
	if day == time.Sunday {
		return DaysPerWeek
	}
	return int(day)
}

func WeekOfYear(t time.Time) int {
	_, week := t.ISOWeek()
	return week
}

func WeekOfMonth(t time.Time) int {
	days := t.Day() + DayOfWeek(StartOfMonth(t)) - 1
	if days%DaysPerWeek == 0 {
		return days / DaysPerWeek
	}
	return days/DaysPerWeek + 1
}

func Date(t time.Time) (year int, month time.Month, day int) {
	return t.Date()
}

func DateTime(t time.Time) (year int, month time.Month, day, hour, minute, second int) {
	year, month, day = t.Date()
	hour, minute, second = t.Clock()
	return
}

func DateTimeMilli(t time.Time) (year int, month time.Month, day, hour, minute, second, millisecond int) {
	year, month, day = t.Date()
	hour, minute, second = t.Clock()
	millisecond = Millisecond(t)
	return
}

func DateTimeMicro(t time.Time) (year int, month time.Month, day, hour, minute, second, microsecond int) {
	year, month, day = t.Date()
	hour, minute, second = t.Clock()
	microsecond = Microsecond(t)
	return
}

func DateTimeNano(t time.Time) (year int, month time.Month, day, hour, minute, second, nanosecond int) {
	year, month, day = t.Date()
	hour, minute, second = t.Clock()
	nanosecond = t.Nanosecond()
	return
}

func DateMilli(t time.Time) (year int, month time.Month, day, millisecond int) {
	year, month, day, _, _, _, millisecond = DateTimeMilli(t)
	return
}

func DateMicro(t time.Time) (year int, month time.Month, day, microsecond int) {
	year, month, day, _, _, _, microsecond = DateTimeMicro(t)
	return
}

func DateNano(t time.Time) (year int, month time.Month, day, nanosecond int) {
	year, month, day, _, _, _, nanosecond = DateTimeNano(t)
	return
}

func Time(t time.Time) (hour, minute, second int) {
	return t.Clock()
}

func Millisecond(t time.Time) (millisecond int) {
	return t.Nanosecond() / 1e6
}

func Microsecond(t time.Time) (microsecond int) {
	return t.Nanosecond() / 1e3
}

func Century(t time.Time) int {
	return t.Year()/YearsPerCentury + 1
}

func Decade(t time.Time) int {
	return t.Year() % YearsPerCentury / YearsPerDecade * YearsPerDecade
}

// Quarter gets current quarter like 3
func Quarter(t time.Time) int {
	switch t.Month() {
	case time.January, time.February, time.March:
		return 1
	case time.April, time.May, time.June:
		return 2
	case time.July, time.August, time.September:
		return 3
	case time.October, time.November, time.December:
		return 4
	}
	return 0
}

// Week gets current week like 6, start from 0
func Week(t time.Time, weekStartsAt time.Weekday) int {
	return (DayOfWeek(t) + DaysPerWeek - int(weekStartsAt)) % DaysPerWeek
}

func ZoneOffset(t time.Time) int {
	if t.IsZero() {
		_, offset := time.Now().In(t.Location()).Zone()
		return offset
	}
	_, offset := t.Zone()
	return offset
}
