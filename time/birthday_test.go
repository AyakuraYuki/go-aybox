package time

import (
	"maps"
	"slices"
	"testing"
	"time"

	"github.com/AyakuraYuki/go-aybox/time/calendar/lunar"
)

const testCreatedAtYear = 2025

func testGetDeltaYears() int { return time.Now().Year() - testCreatedAtYear }

func testAdjustWantRealAge(given time.Time, age int) int {
	delta := testGetDeltaYears()
	if delta <= 0 {
		return age
	}
	now := time.Now()
	if now.Month() < given.Month() || (now.Month() == given.Month() && now.Day() < given.Day()) {
		return age + delta - 1
	}
	return age + delta
}

// created at 2025
var (
	testBirthdayAgeMapping = map[string]int{
		"1995-03-14": 30,
		"1995-11-05": 29,
		"1996-01-19": 29,
		"1996-02-29": 29,
		"1996-12-26": 28,
	}

	testBirthdayNominalAgeMapping = map[string]int{
		"1995-03-14": 31,
		"1995-11-05": 31,
		"1996-01-19": 31,
		"1996-02-29": 30,
		"1996-12-26": 30,
	}
)

func testDateStringParseToTime(date string) time.Time {
	d, _ := time.Parse(time.DateOnly, date)
	return d
}

// ----------------------------------------------------------------------------------------------------

func TestCalculateRealAge(t *testing.T) {
	type testCase struct {
		birthday time.Time
		want     int
	}
	var tests []testCase
	for _, date := range slices.Sorted(maps.Keys(testBirthdayAgeMapping)) {
		birthday := testDateStringParseToTime(date)
		tests = append(tests, testCase{birthday, testAdjustWantRealAge(birthday, testBirthdayAgeMapping[date])})
	}

	for _, tt := range tests {
		t.Run(tt.birthday.Format(time.DateOnly), func(t *testing.T) {
			if get := CalculateRealAge[int](tt.birthday); get != tt.want {
				t.Errorf("unexpected age, want %v, but got %v", tt.want, get)
			}
		})
	}
}

func TestCalculateRealAgeFromString(t *testing.T) {
	type testCase struct {
		birthday string
		want     int
	}
	var tests []testCase
	for _, date := range slices.Sorted(maps.Keys(testBirthdayAgeMapping)) {
		tests = append(tests, testCase{date, testAdjustWantRealAge(testDateStringParseToTime(date), testBirthdayAgeMapping[date])})
	}

	for _, tt := range tests {
		t.Run(tt.birthday, func(t *testing.T) {
			if get := CalculateRealAgeFromString[int](tt.birthday); get != tt.want {
				t.Errorf("unexpected age, want %v, but got %v", tt.want, get)
			}
		})
	}
}

// ----------------------------------------------------------------------------------------------------

func TestCalculateNominalAge(t *testing.T) {
	type testCase struct {
		birthday time.Time
		want     int
	}
	var tests []testCase
	for _, date := range slices.Sorted(maps.Keys(testBirthdayNominalAgeMapping)) {
		tests = append(tests, testCase{testDateStringParseToTime(date), testBirthdayNominalAgeMapping[date]})
	}

	// adjust want
	loc, _ := time.LoadLocation(AsiaShanghai)
	now := time.Now().In(loc)
	if lunarNewYear := lunar.FromStdTime(now); now.After(lunarNewYear.ToGregorian(AsiaShanghai).Time) {
		for i := range tests {
			tests[i].want += testGetDeltaYears()
		}
	}

	for _, tt := range tests {
		t.Run(tt.birthday.Format(time.DateOnly), func(t *testing.T) {
			if get := CalculateNominalAge[int](tt.birthday); get != tt.want {
				t.Errorf("unexpected nominal age, want %v, but got %v", tt.want, get)
			}
		})
	}
}
