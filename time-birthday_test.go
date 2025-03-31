package aybox

import (
	"maps"
	"slices"
	"testing"
	"time"

	"github.com/dromara/carbon/v2"
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
		"1996-02-29": 31,
		"1996-12-26": 30,
	}
)

func testDateStringParseToTime(date string) time.Time {
	d, _ := time.Parse(time.DateOnly, date)
	return d
}

func testDateStringParseToCarbon(date string) *carbon.Carbon {
	c := carbon.ParseByLayout(date, time.DateOnly)
	if c.HasError() {
		return nil
	}
	return c
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
	now := carbon.Now(carbon.Shanghai)
	if lunarNewYear := lunarNewYearInThatYear(now); now.Gte(lunarNewYear) {
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

func TestCalculateNominalAgeFromCarbon(t *testing.T) {
	type testCase struct {
		birthday *carbon.Carbon
		want     int
	}
	var tests []testCase
	for _, date := range slices.Sorted(maps.Keys(testBirthdayNominalAgeMapping)) {
		tests = append(tests, testCase{carbon.ParseByLayout(date, time.DateOnly), testBirthdayNominalAgeMapping[date]})
	}

	// adjust want
	now := carbon.Now(carbon.Shanghai)
	if lunarNewYear := lunarNewYearInThatYear(now); now.Gte(lunarNewYear) {
		for i := range tests {
			tests[i].want += testGetDeltaYears()
		}
	}

	for _, tt := range tests {
		t.Run(tt.birthday.Format("Y-m-d"), func(t *testing.T) {
			if get := CalculateNominalAgeFromCarbon[int](tt.birthday); get != tt.want {
				t.Errorf("unexpected nominal age, want %v, but got %v", tt.want, get)
			}
		})
	}
}

// ----------------------------------------------------------------------------------------------------
