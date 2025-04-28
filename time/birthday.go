package time

import (
	"strings"
	"time"

	"github.com/AyakuraYuki/go-aybox/time/calendar/lunar"
)

type AgeNumber interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// CalculateRealAge returns the real age from birthday to now.
//
// 计算一个出生日期到现在的实岁。
func CalculateRealAge[T AgeNumber](birthday time.Time) (age T) {
	date := birthday
	now := time.Now().In(date.Location())

	if !IsLeapYear(now) && date.Month() == time.February && date.Day() == 29 {
		// in non-leap years, 02-29 is converted to 03-01
		date = date.AddDate(0, 0, 1)
	}

	if now.Before(date) {
		// man, are you sure your birthday is before than now?
		return 0
	}

	years := now.Year() - date.Year()
	if now.Month() < date.Month() || (now.Month() == date.Month() && now.Day() < date.Day()) {
		years--
	}
	return T(years)
}

// CalculateRealAgeFromString accepts a date string in time.DateOnly
// to calculate the real age. Use UTC timezone for calculation.
//
// 接受一个格式符合 time.DateOnly 的日期，计算现在到这个日期的实岁。使用UTC零时区计算。
func CalculateRealAgeFromString[T AgeNumber](date string) (age T) {
	date = strings.TrimSpace(date)
	if date == "" {
		return T(0)
	}
	birthday, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return T(0)
	}
	return CalculateRealAge[T](birthday)
}

// CalculateNominalAge returns the Xusui (known as nominal age) from some
// birthday to now.
//
// Xusui, the nominal age system, a traditional Chinese age reckoning
// method, is consistently calculated using the Asia/Shanghai time zone
// standard.
//
// 计算一个出生日期到现在的虚岁。虚岁是传统中国的记岁方法，固定使用 Asia/Shanghai 时区。
func CalculateNominalAge[T AgeNumber](birthday time.Time) (age T) {
	loc, _ := time.LoadLocation(AsiaShanghai)

	birthday = birthday.In(loc)
	lunarNewYearInBirthYear := lunar.FromLunar(birthday.Year(), 1, 1, false).ToGregorian(AsiaShanghai)

	now := time.Now().In(loc)
	lunarNewYear := lunar.FromLunar(now.Year(), 1, 1, false).ToGregorian(AsiaShanghai)

	years := now.Year() - birthday.Year()
	if birthday.Before(lunarNewYearInBirthYear.Time) {
		years++
	}
	if now.After(lunarNewYear.Time) {
		years++
	}
	return T(years)
}

// ----------------------------------------------------------------------------------------------------

// // CalculateRealAgeFromCarbon returns the real age from birthday to now.
// //
// // 计算一个出生日期到现在的实岁。
// func CalculateRealAgeFromCarbon[T AgeNumber](birthday *carbon.Carbon) (age T) {
// 	if birthday == nil {
// 		return T(0) // 0 when birthday is nil
// 	}
//
// 	date := birthday.Copy() // do not modify the original date
// 	if date.HasError() {
// 		// cannot calculate an age with broken time instance
// 		return T(0)
// 	}
//
// 	now := carbon.Now(date.Timezone())
// 	if !now.IsLeapYear() && date.Month() == 2 && date.Day() == 29 {
// 		// in non-leap years, 02-29 is converted to 03-01
// 		date = date.AddDay()
// 	}
//
// 	if now.Lt(date) {
// 		// man, are you sure your birthday is before than now?
// 		return 0
// 	}
//
// 	years := now.Year() - date.Year()
// 	if now.Month() < date.Month() || (now.Month() == date.Month() && now.Day() < date.Day()) {
// 		years--
// 	}
// 	return T(years)
// }
//
// // CalculateNominalAgeFromCarbon returns the Xusui (known as nominal age) from
// // some birthday to now.
// //
// // Xusui, the nominal age system, a traditional Chinese age reckoning
// // method, is consistently calculated using the Asia/Shanghai time zone
// // standard.
// //
// // 计算一个出生日期到现在的虚岁。虚岁是传统中国的记岁方法，固定使用 Asia/Shanghai 时区。
// func CalculateNominalAgeFromCarbon[T AgeNumber](birthday *carbon.Carbon) (age T) {
// 	timezone := carbon.Shanghai
//
// 	birthday = birthday.SetTimezone(timezone)
// 	lunarNewYearInBirthYear := lunarNewYearInThatYear(birthday)
//
// 	now := carbon.Now(timezone)
// 	lunarNewYear := lunarNewYearInThatYear(now)
//
// 	years := now.Year() - birthday.Year()
// 	if birthday.Lt(lunarNewYearInBirthYear) {
// 		years++
// 	}
// 	if now.Gte(lunarNewYear) {
// 		years++
// 	}
// 	return T(years)
// }
//
// func lunarNewYearInThatYear(date *carbon.Carbon) (lunarNewYear *carbon.Carbon) {
// 	timezone := carbon.Shanghai
//
// 	copied := date.Copy()
// 	copied = copied.SetTimezone(timezone)
//
// 	lunarNewYear = carbon.CreateFromLunar(copied.Year(), 1, 1, 0, 0, 0, false)
// 	lunarNewYear = lunarNewYear.SetTimezone(timezone)
//
// 	return
// }
