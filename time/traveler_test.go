package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTomorrow(t *testing.T) {
	t.Run("invalid timezone", func(t *testing.T) {
		assert.Equal(t, time.Time{}.String(), Tomorrow("abc").String())
		assert.Equal(t, time.Time{}.String(), Tomorrow("xxx").String())
	})

	t.Run("without timezone", func(t *testing.T) {
		assert.Equal(t, time.Now().Add(24*time.Hour).Format(time.DateOnly), Tomorrow().Format(time.DateOnly))
	})

	t.Run("with timezone", func(t *testing.T) {
		assert.Equal(t, time.Now().In(time.UTC).Add(24*time.Hour).Format(time.DateOnly), Tomorrow("UTC").Format(time.DateOnly))
	})
}

func TestYesterday(t *testing.T) {
	t.Run("invalid timezone", func(t *testing.T) {
		assert.Equal(t, time.Time{}.String(), Yesterday("abc").String())
		assert.Equal(t, time.Time{}.String(), Yesterday("xxx").String())
	})

	t.Run("without timezone", func(t *testing.T) {
		assert.Equal(t, time.Now().Add(-24*time.Hour).Format(time.DateOnly), Yesterday().Format(time.DateOnly))
	})

	t.Run("with timezone", func(t *testing.T) {
		assert.Equal(t, time.Now().In(time.UTC).Add(-24*time.Hour).Format(time.DateOnly), Yesterday("UTC").Format(time.DateOnly))
	})
}

func TestAddMonths(t *testing.T) {
	t.Run("zero time", func(t *testing.T) {
		assert.Equal(t, "0001-03-01", AddMonths(time.Time{}, 2).Format(time.DateOnly))
	})

	t.Run("valid time", func(t *testing.T) {
		assert.Equal(t, "2025-01-01", AddMonths(StartOfYear(time.Now()), 0).Format(time.DateOnly))
		assert.Equal(t, "2025-02-01", AddMonths(StartOfYear(time.Now()), 1).Format(time.DateOnly))
		assert.Equal(t, "2025-03-01", AddMonths(StartOfYear(time.Now()), 2).Format(time.DateOnly))
		assert.Equal(t, "2000-05-29", AddMonths(NewDate(2000, 2, 29, time.UTC), 3).Format(time.DateOnly))
		assert.Equal(t, "2025-10-31", AddMonths(NewDate(2025, 8, 31, time.UTC), 2).Format(time.DateOnly))
	})
}

func TestAddMonthsNoOverflow(t *testing.T) {
	t.Run("zero time", func(t *testing.T) {
		assert.Equal(t, "0001-03-01", AddMonthsNoOverflow(time.Time{}, 2).Format(time.DateOnly))
	})

	t.Run("valid time", func(t *testing.T) {
		assert.Equal(t, "2025-01-01", AddMonthsNoOverflow(StartOfYear(time.Now()), 0).Format(time.DateOnly))
		assert.Equal(t, "2025-02-01", AddMonthsNoOverflow(StartOfYear(time.Now()), 1).Format(time.DateOnly))
		assert.Equal(t, "2025-03-01", AddMonthsNoOverflow(StartOfYear(time.Now()), 2).Format(time.DateOnly))
		assert.Equal(t, "2000-05-29", AddMonthsNoOverflow(NewDate(2000, 2, 29, time.UTC), 3).Format(time.DateOnly))
		assert.Equal(t, "2025-10-31", AddMonthsNoOverflow(NewDate(2025, 8, 31, time.UTC), 2).Format(time.DateOnly))
	})
}
