package aybox

import (
	"strings"
	"time"

	"github.com/dromara/carbon/v2"
)

type AgeNumber interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64
}

// region Real Age

// CalculateRealAge returns the real age from some birthday to now, accept
// given timezone for the calculation.
//
// There are some examples of timezone, like `America/Los_Angeles` and
// `Asia/Shanghai`.
//
// 计算一个出生日期到现在的实岁。可以指定时区。
// 这里给出一些时区的例子，如 `America/Los_Angeles`，或者 `Asia/Shanghai`。
func CalculateRealAge[T AgeNumber](birthday time.Time, timezone ...string) (age T) {
	now := time.Now()
	if len(timezone) > 0 {
		loc, err := time.LoadLocation(timezone[0])
		if err != nil {
			return 0 // unknown timezone resulting 0 age
		}
		birthday = birthday.In(loc)
		now = now.In(loc)
	}
	if now.Before(birthday) {
		return 0
	}
	years := now.Year() - birthday.Year()
	if now.Month() < birthday.Month() || (now.Month() == birthday.Month() && now.Day() < birthday.Day()) {
		years--
	}
	return T(years)
}

// CalculateRealAgeFromString accepts a date string in layout `2006-01-02`
// to calculate the real age, accept given timezone for the calculation.
//
// There are some examples of timezone, like `America/Los_Angeles` and
// `Asia/Shanghai`.
//
// 接受一个格式符合 `2006-01-02` 的日期，计算现在到这个日期的实岁。可以指定时区。
// 这里给出一些时区的例子，如 `America/Los_Angeles`，或者 `Asia/Shanghai`。
func CalculateRealAgeFromString[T AgeNumber](date string, timezone ...string) (age T) {
	if strings.TrimSpace(date) == "" {
		return T(0)
	}
	birthday := carbon.ParseByLayout(date, time.DateOnly, timezone...)
	return CalculateRealAgeFromCarbon[T](birthday, timezone...)
}

// CalculateRealAgeFromCarbon returns the real age from some birthday to now, accept
// given timezone for the calculation.
//
// There are some examples of timezone, like `America/Los_Angeles` and
// `Asia/Shanghai`.
//
// 计算一个出生日期到现在的实岁。可以指定时区。
// 这里给出一些时区的例子，如 `America/Los_Angeles`，或者 `Asia/Shanghai`。
func CalculateRealAgeFromCarbon[T AgeNumber](birthday carbon.Carbon, timezone ...string) (age T) {
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

// endregion

// region Nominal Age

// CalculateNominalAge returns the Xusui (known as nominal age) from some
// birthday to now.
//
// Xusui, the nominal age system, a traditional Chinese age reckoning
// method, is consistently calculated using the Asia/Shanghai time zone
// standard.
//
// 计算一个出生日期到现在的虚岁。虚岁是传统中国的记岁方法，固定使用 Asia/Shanghai 时区。
func CalculateNominalAge[T AgeNumber](birthday time.Time) (age T) {
	carbonBirthday := carbon.CreateFromStdTime(birthday, carbon.Shanghai)
	return CalculateNominalAgeFromCarbon[T](carbonBirthday)
}

// CalculateNominalAgeFromCarbon returns the Xusui (known as nominal age) from
// some birthday to now.
//
// Xusui, the nominal age system, a traditional Chinese age reckoning
// method, is consistently calculated using the Asia/Shanghai time zone
// standard.
//
// 计算一个出生日期到现在的虚岁。虚岁是传统中国的记岁方法，固定使用 Asia/Shanghai 时区。
func CalculateNominalAgeFromCarbon[T AgeNumber](birthday carbon.Carbon) (age T) {
	timezone := carbon.Shanghai

	birthday = birthday.SetTimezone(timezone)
	lunarNewYearInBirthYear := carbon.CreateFromLunar(birthday.Year(), 1, 1, 0, 0, 0, birthday.IsLeapYear())
	lunarNewYearInBirthYear = lunarNewYearInBirthYear.SetTimezone(timezone)

	now := carbon.Now(timezone)
	lunarNewYear := carbon.CreateFromLunar(now.Year(), 1, 1, 0, 0, 0, now.IsLeapYear())
	lunarNewYear = lunarNewYear.SetTimezone(timezone)

	years := now.Year() - birthday.Year()
	if birthday.Lt(lunarNewYearInBirthYear) {
		years++
	}
	if now.Gte(lunarNewYear) {
		years++
	}
	return T(years)
}

// endregion
