package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsAM(t *testing.T) {
	assert.True(t, IsAM(time.Time{}))
	assert.True(t, IsAM(time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC)))
	assert.False(t, IsAM(time.Date(2020, 1, 1, 14, 0, 0, 0, time.UTC)))
}

func TestIsPM(t *testing.T) {
	assert.False(t, IsPM(time.Time{}))
	assert.False(t, IsPM(time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC)))
	assert.True(t, IsPM(time.Date(2020, 1, 1, 14, 0, 0, 0, time.UTC)))
}

func TestIsLeapYear(t *testing.T) {
	// invalid time
	assert.False(t, IsLeapYear(time.Time{}))

	assert.True(t, IsLeapYear(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))
	assert.False(t, IsLeapYear(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)))
}

func TestIsLongYear(t *testing.T) {
	t.Run("invalid time", func(t *testing.T) {
		assert.False(t, IsLongYear(time.Time{}))
	})

	t.Run("valid time", func(t *testing.T) {
		assert.True(t, IsLongYear(NewDate(2015, 1, 1, time.UTC)))
		assert.False(t, IsLongYear(NewDate(2016, 1, 1, time.UTC)))
	})
}

func TestCompare(t *testing.T) {
	t.Run("zero time", func(t *testing.T) {
		assert.True(t, Compare(time.Time{}, "==", time.Time{}))
		assert.False(t, Compare(time.Now(), "==", time.Time{}))
		assert.False(t, Compare(time.Time{}, "==", time.Now()))
	})

	t.Run("invalid operator", func(t *testing.T) {
		assert.False(t, Compare(Yesterday(), "", Tomorrow()))
		assert.False(t, Compare(Tomorrow(), "%", time.Now()))
	})

	t.Run("valid time", func(t *testing.T) {
		assert.True(t, Compare(StartOfDay(Yesterday()), "==", StartOfDay(Yesterday())))

		assert.True(t, Compare(StartOfDay(Yesterday()), "!=", StartOfDay(Tomorrow())))
		assert.True(t, Compare(StartOfDay(Yesterday()), "<>", StartOfDay(Tomorrow())))

		assert.True(t, Compare(time.Now(), ">", Yesterday()))
		assert.True(t, Compare(time.Now(), ">=", Yesterday()))

		assert.True(t, Compare(time.Now(), "<", Tomorrow()))
		assert.True(t, Compare(time.Now(), "<=", Tomorrow()))
	})
}

func TestBetween(t *testing.T) {
	t.Run("zero time", func(t *testing.T) {
		emptyTime := time.Time{}
		aprilFool := NewDate(2025, 4, 1, time.UTC)
		nextAprilFool := AddDay(aprilFool)

		assert.False(t, Between(emptyTime, emptyTime, emptyTime))
		assert.False(t, Between(emptyTime, aprilFool, emptyTime))
		assert.False(t, Between(aprilFool, emptyTime, emptyTime))
		assert.False(t, Between(aprilFool, emptyTime, aprilFool))
		assert.True(t, Between(aprilFool, emptyTime, nextAprilFool))
	})

	t.Run("valid time", func(t *testing.T) {
		sampleTime := time.Date(2025, 4, 1, 12, 0, 0, 0, time.UTC)
		lastHour := SubHour(sampleTime)
		nextHour := AddHour(sampleTime)

		assert.True(t, Between(sampleTime, lastHour, nextHour))
		assert.False(t, Between(sampleTime, sampleTime, nextHour))
		assert.False(t, Between(sampleTime, lastHour, sampleTime))
		assert.False(t, Between(sampleTime, sampleTime, sampleTime))
	})
}

func TestBetweenIncludedStart(t *testing.T) {
	t.Run("zero time", func(t *testing.T) {
		emptyTime := time.Time{}
		aprilFool := NewDate(2025, 4, 1, time.UTC)
		nextAprilFool := AddDay(aprilFool)

		assert.False(t, BetweenIncludedStart(emptyTime, emptyTime, emptyTime))
		assert.False(t, BetweenIncludedStart(emptyTime, aprilFool, emptyTime))
		assert.False(t, BetweenIncludedStart(aprilFool, emptyTime, emptyTime))
		assert.False(t, BetweenIncludedStart(aprilFool, emptyTime, aprilFool))
		assert.True(t, BetweenIncludedStart(aprilFool, emptyTime, nextAprilFool))
		assert.True(t, BetweenIncludedStart(aprilFool, aprilFool, nextAprilFool))
	})

	t.Run("valid time", func(t *testing.T) {
		sampleTime := time.Date(2025, 4, 1, 12, 0, 0, 0, time.UTC)
		lastHour := SubHour(sampleTime)
		nextHour := AddHour(sampleTime)

		assert.True(t, BetweenIncludedStart(sampleTime, lastHour, nextHour))
		assert.True(t, BetweenIncludedStart(sampleTime, sampleTime, nextHour))
		assert.False(t, BetweenIncludedStart(sampleTime, lastHour, sampleTime))
		assert.False(t, BetweenIncludedStart(sampleTime, sampleTime, sampleTime))
	})
}

func TestBetweenIncludedEnd(t *testing.T) {
	t.Run("zero time", func(t *testing.T) {
		emptyTime := time.Time{}
		aprilFool := NewDate(2025, 4, 1, time.UTC)
		nextAprilFool := AddDay(aprilFool)

		assert.False(t, BetweenIncludedEnd(emptyTime, emptyTime, emptyTime))
		assert.False(t, BetweenIncludedEnd(emptyTime, aprilFool, emptyTime))
		assert.False(t, BetweenIncludedEnd(aprilFool, emptyTime, emptyTime))
		assert.False(t, BetweenIncludedEnd(aprilFool, aprilFool, nextAprilFool))
		assert.True(t, BetweenIncludedEnd(aprilFool, emptyTime, aprilFool))
		assert.True(t, BetweenIncludedEnd(aprilFool, emptyTime, nextAprilFool))
	})

	t.Run("valid time", func(t *testing.T) {
		sampleTime := time.Date(2025, 4, 1, 12, 0, 0, 0, time.UTC)
		lastHour := SubHour(sampleTime)
		nextHour := AddHour(sampleTime)

		assert.True(t, BetweenIncludedEnd(sampleTime, lastHour, nextHour))
		assert.False(t, BetweenIncludedEnd(sampleTime, sampleTime, nextHour))
		assert.True(t, BetweenIncludedEnd(sampleTime, lastHour, sampleTime))
		assert.False(t, BetweenIncludedEnd(sampleTime, sampleTime, sampleTime))
	})
}

func TestBetweenIncludedBoth(t *testing.T) {
	t.Run("zero time", func(t *testing.T) {
		emptyTime := time.Time{}
		aprilFool := NewDate(2025, 4, 1, time.UTC)
		nextAprilFool := AddDay(aprilFool)

		assert.True(t, BetweenIncludedBoth(emptyTime, emptyTime, emptyTime))
		assert.False(t, BetweenIncludedBoth(emptyTime, aprilFool, emptyTime))
		assert.False(t, BetweenIncludedBoth(aprilFool, emptyTime, emptyTime))
		assert.True(t, BetweenIncludedBoth(aprilFool, aprilFool, nextAprilFool))
		assert.True(t, BetweenIncludedBoth(aprilFool, emptyTime, aprilFool))
		assert.True(t, BetweenIncludedBoth(aprilFool, emptyTime, nextAprilFool))
	})

	t.Run("valid time", func(t *testing.T) {
		sampleTime := time.Date(2025, 4, 1, 12, 0, 0, 0, time.UTC)
		lastHour := SubHour(sampleTime)
		nextHour := AddHour(sampleTime)

		assert.True(t, BetweenIncludedBoth(sampleTime, lastHour, nextHour))
		assert.True(t, BetweenIncludedBoth(sampleTime, sampleTime, nextHour))
		assert.True(t, BetweenIncludedBoth(sampleTime, lastHour, sampleTime))
		assert.True(t, BetweenIncludedBoth(sampleTime, sampleTime, sampleTime))
	})
}
