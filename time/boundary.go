package time

import "time"

func StartOfCentury(t time.Time) time.Time {
	return time.Date(t.Year()/YearsPerCentury*YearsPerCentury, 1, 1, 0, 0, 0, 0, t.Location())
}

func EndOfCentury(t time.Time) time.Time {
	return time.Date(t.Year()/YearsPerCentury*YearsPerCentury+99, 12, 31, 23, 59, 59, 999999999, t.Location())
}

func StartOfDecade(t time.Time) time.Time {
	return time.Date(t.Year()/YearsPerDecade*YearsPerDecade, 1, 1, 0, 0, 0, 0, t.Location())
}

func EndOfDecade(t time.Time) time.Time {
	return time.Date(t.Year()/YearsPerDecade*YearsPerDecade+9, 12, 31, 23, 59, 59, 999999999, t.Location())
}

func StartOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

func EndOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 12, 31, 23, 59, 59, 999999999, t.Location())
}

func StartOfQuarter(t time.Time) time.Time {
	year, quarter, day := t.Year(), Quarter(t), 1
	return time.Date(year, time.Month(3*quarter-2), day, 0, 0, 0, 0, t.Location())
}

func EndOfQuarter(t time.Time) time.Time {
	year, quarter, day := t.Year(), Quarter(t), 30
	if quarter == 1 || quarter == 4 {
		day = 31
	}
	return time.Date(year, time.Month(3*quarter), day, 23, 59, 59, 999999999, t.Location())
}

func StartOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

func EndOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month+1, 0, 23, 59, 59, 999999999, t.Location())
}

func StartOfWeek(t time.Time, weekStartsAt time.Weekday) time.Time {
	dayOfWeek := DayOfWeek(t)
	offset := -((DaysPerWeek + dayOfWeek - int(weekStartsAt)) % DaysPerWeek)
	return StartOfDay(t.AddDate(0, 0, offset))
}

func EndOfWeek(t time.Time, weekStartsAt time.Weekday) time.Time {
	dayOfWeek := DayOfWeek(t)
	weekEndsAt := weekStartsAt + DaysPerWeek - 1
	offset := (DaysPerWeek - dayOfWeek + int(weekEndsAt)) % DaysPerWeek
	return EndOfDay(t.AddDate(0, 0, offset))
}

func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}

func StartOfHour(t time.Time) time.Time {
	year, month, day, hour, _, _ := DateTime(t)
	return time.Date(year, month, day, hour, 0, 0, 0, t.Location())
}

func EndOfHour(t time.Time) time.Time {
	year, month, day, hour, _, _ := DateTime(t)
	return time.Date(year, month, day, hour, 59, 59, 999999999, t.Location())
}

func StartOfMinute(t time.Time) time.Time {
	year, month, day, hour, minute, _ := DateTime(t)
	return time.Date(year, month, day, hour, minute, 0, 0, t.Location())
}

func EndOfMinute(t time.Time) time.Time {
	year, month, day, hour, minute, _ := DateTime(t)
	return time.Date(year, month, day, hour, minute, 59, 999999999, t.Location())
}

func StartOfSecond(t time.Time) time.Time {
	year, month, day, hour, minute, second := DateTime(t)
	return time.Date(year, month, day, hour, minute, second, 0, t.Location())
}

func EndOfSecond(t time.Time) time.Time {
	year, month, day, hour, minute, second := DateTime(t)
	return time.Date(year, month, day, hour, minute, second, 999999999, t.Location())
}
