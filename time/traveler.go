package time

import (
	"time"
)

// Tomorrow returns a time.Time instance for tomorrow, or empty time when the
// timezone is wrong.
func Tomorrow(timezone ...string) time.Time {
	now := time.Now()
	if len(timezone) > 0 {
		loc, err := loadLocationByTimezone(timezone[0])
		if err != nil || loc == nil {
			return time.Time{}
		}
		now = now.In(loc)
	}
	return AddDay(now)
}

// Yesterday returns a time.Time instance for yesterday, or empty time when the
// timezone is wrong.
func Yesterday(timezone ...string) time.Time {
	now := time.Now()
	if len(timezone) > 0 {
		loc, err := loadLocationByTimezone(timezone[0])
		if err != nil || loc == nil {
			return time.Time{}
		}
		now = now.In(loc)
	}
	return SubDay(now)
}

func AddMonths(t time.Time, months int) time.Time {
	return t.AddDate(0, months, 0)
}

func AddMonth(t time.Time) time.Time {
	return AddMonths(t, 1)
}

func AddMonthsNoOverflow(t time.Time, months int) time.Time {
	nanoseconds := t.Nanosecond()
	year, month, day, hour, minute, second := DateTime(t)
	// get the last day after N months
	lastYear, lastMonth, lastDay := Date(time.Date(year, month+time.Month(months)+1, 0, hour, minute, second, nanoseconds, t.Location()))
	if day > lastDay {
		day = lastDay
	}
	return time.Date(lastYear, lastMonth, day, hour, minute, second, nanoseconds, t.Location())
}

func AddMonthNoOverflow(t time.Time) time.Time {
	return AddMonthsNoOverflow(t, 1)
}

func SubMonths(t time.Time, months int) time.Time {
	return AddMonths(t, -months)
}

func SubMonth(t time.Time) time.Time {
	return SubMonths(t, 1)
}

func SubMonthsNoOverflow(t time.Time, months int) time.Time {
	return AddMonthsNoOverflow(t, -months)
}

func SubMonthNoOverflow(t time.Time) time.Time {
	return SubMonthsNoOverflow(t, 1)
}

func AddWeeks(t time.Time, weeks int) time.Time {
	return AddDays(t, weeks*DaysPerWeek)
}

func AddWeek(t time.Time) time.Time {
	return AddWeeks(t, 1)
}

func SubWeeks(t time.Time, weeks int) time.Time {
	return SubDays(t, weeks*DaysPerWeek)
}

func SubWeek(t time.Time) time.Time {
	return SubWeeks(t, 1)
}

func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

func AddDay(t time.Time) time.Time {
	return AddDays(t, 1)
}

func SubDays(t time.Time, days int) time.Time {
	return AddDays(t, -days)
}

func SubDay(t time.Time) time.Time {
	return SubDays(t, 1)
}

func AddHours(t time.Time, hours int) time.Time {
	return t.Add(time.Duration(hours) * time.Hour)
}

func AddHour(t time.Time) time.Time {
	return AddHours(t, 1)
}

func SubHours(t time.Time, hours int) time.Time {
	return AddHours(t, -hours)
}

func SubHour(t time.Time) time.Time {
	return SubHours(t, 1)
}

func AddMinutes(t time.Time, minutes int) time.Time {
	return t.Add(time.Duration(minutes) * time.Minute)
}

func AddMinute(t time.Time) time.Time {
	return AddMinutes(t, 1)
}

func SubMinutes(t time.Time, minutes int) time.Time {
	return AddMinutes(t, -minutes)
}

func SubMinute(t time.Time) time.Time {
	return SubMinutes(t, 1)
}

func AddSeconds(t time.Time, seconds int) time.Time {
	return t.Add(time.Duration(seconds) * time.Second)
}

func AddSecond(t time.Time) time.Time {
	return AddSeconds(t, 1)
}

func SubSeconds(t time.Time, seconds int) time.Time {
	return AddSeconds(t, -seconds)
}

func SubSecond(t time.Time) time.Time {
	return SubSeconds(t, 1)
}

func AddMilliseconds(t time.Time, milliseconds int) time.Time {
	return t.Add(time.Duration(milliseconds) * time.Millisecond)
}

func AddMillisecond(t time.Time) time.Time {
	return AddMilliseconds(t, 1)
}

func SubMilliseconds(t time.Time, milliseconds int) time.Time {
	return AddMilliseconds(t, -milliseconds)
}

func SubMillisecond(t time.Time) time.Time {
	return SubMilliseconds(t, 1)
}

func AddMicroseconds(t time.Time, microseconds int) time.Time {
	return t.Add(time.Duration(microseconds) * time.Microsecond)
}

func AddMicrosecond(t time.Time) time.Time {
	return AddMicroseconds(t, 1)
}

func SubMicroseconds(t time.Time, microseconds int) time.Time {
	return AddMicroseconds(t, -microseconds)
}

func SubMicrosecond(t time.Time) time.Time {
	return SubMicroseconds(t, 1)
}

func AddNanoseconds(t time.Time, nanoseconds int) time.Time {
	return t.Add(time.Duration(nanoseconds) * time.Nanosecond)
}

func AddNanosecond(t time.Time) time.Time {
	return AddNanoseconds(t, 1)
}

func SubNanoseconds(t time.Time, nanoseconds int) time.Time {
	return AddNanoseconds(t, -nanoseconds)
}

func SubNanosecond(t time.Time) time.Time {
	return SubNanoseconds(t, 1)
}
