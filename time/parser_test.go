package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Run("empty time", func(t *testing.T) {
		tt, _ := Parse("")
		assert.Equal(t, time.Time{}.String(), tt.String())
	})

	t.Run("invalid timezone", func(t *testing.T) {
		tt, err := Parse("2025-04-01", "")
		assert.Equal(t, "timezone cannot be empty", err.Error())
		assert.Equal(t, time.Time{}.String(), tt.String())
		tt, err = Parse("2025-04-01", "xxx")
		assert.Equal(t, `invalid timezone "xxx", please see the file "$GOROOT/lib/time/zoneinfo.zip" for all valid timezones`, err.Error())
		assert.Equal(t, time.Time{}.String(), tt.String())
	})

	t.Run("invalid time", func(t *testing.T) {
		_, err := Parse("0")
		assert.Error(t, err)
		_, err = Parse("xxx")
		assert.Error(t, err)
	})

	t.Run("valid time without timezone", func(t *testing.T) {
		a := time.Now()
		b, _ := Parse("now")
		assert.Equal(t, a.Unix(), b.Unix())

		a = Yesterday()
		b, _ = Parse("yesterday")
		assert.Equal(t, a.Unix(), b.Unix())

		a = Tomorrow()
		b, _ = Parse("tomorrow")
		assert.Equal(t, a.Unix(), b.Unix())

		tt, _ := Parse("2020-8-5")
		assert.Equal(t, "2020-08-05 00:00:00 +0000 UTC", tt.String())
		tt, _ = Parse("2020-8-05")
		assert.Equal(t, "2020-08-05 00:00:00 +0000 UTC", tt.String())
		tt, _ = Parse("2020-08-05")
		assert.Equal(t, "2020-08-05 00:00:00 +0000 UTC", tt.String())
		tt, _ = Parse("2020-8-5 1:2:3")
		assert.Equal(t, "2020-08-05 01:02:03 +0000 UTC", tt.String())
		tt, _ = Parse("2020-08-05 1:2:03")
		assert.Equal(t, "2020-08-05 01:02:03 +0000 UTC", tt.String())
		tt, _ = Parse("2020-08-05 1:02:03")
		assert.Equal(t, "2020-08-05 01:02:03 +0000 UTC", tt.String())
		tt, _ = Parse("2020-08-05 01:02:03")
		assert.Equal(t, "2020-08-05 01:02:03 +0000 UTC", tt.String())
	})

	t.Run("valid time with timezone", func(t *testing.T) {
		a := time.Now().In(time.UTC)
		b, _ := Parse("now", "UTC")
		assert.Equal(t, a.UnixMilli(), b.UnixMilli())

		a = Yesterday().In(time.UTC)
		b, _ = Parse("yesterday", "UTC")
		assert.Equal(t, a.UnixMilli(), b.UnixMilli())

		a = Tomorrow().In(time.UTC)
		b, _ = Parse("tomorrow", "UTC")
		assert.Equal(t, a.UnixMilli(), b.UnixMilli())

		tt, _ := Parse("2020-8-5")
		assert.Equal(t, "2020-08-05 00:00:00 +0000 UTC", tt.String())
		tt, _ = Parse("2020-8-05")
		assert.Equal(t, "2020-08-05 00:00:00 +0000 UTC", tt.String())
		tt, _ = Parse("2020-08-05")
		assert.Equal(t, "2020-08-05 00:00:00 +0000 UTC", tt.String())
		tt, _ = Parse("2020-8-5 1:2:3")
		assert.Equal(t, "2020-08-05 01:02:03 +0000 UTC", tt.String())
		tt, _ = Parse("2020-08-05 1:2:03")
		assert.Equal(t, "2020-08-05 01:02:03 +0000 UTC", tt.String())
		tt, _ = Parse("2020-08-05 1:02:03")
		assert.Equal(t, "2020-08-05 01:02:03 +0000 UTC", tt.String())
		tt, _ = Parse("2020-08-05 01:02:03")
		assert.Equal(t, "2020-08-05 01:02:03 +0000 UTC", tt.String())
	})

	t.Run("issue202", func(t *testing.T) {
		var tt time.Time
		tt, _ = Parse("2023-01-08T09:02:48")
		assert.Equal(t, "2023-01-08 09:02:48 +0000 UTC", tt.String())
		tt, _ = Parse("2023-1-8T09:02:48")
		assert.Equal(t, "2023-01-08 09:02:48 +0000 UTC", tt.String())
		tt, _ = Parse("2023-01-08T9:2:48")
		assert.Equal(t, "2023-01-08 09:02:48 +0000 UTC", tt.String())
		tt, _ = Parse("2023-01-8T9:2:48")
		assert.Equal(t, "2023-01-08 09:02:48 +0000 UTC", tt.String())
	})

	t.Run("issue232", func(t *testing.T) {
		var tt time.Time
		tt, _ = Parse("0000-01-01 00:00:00")
		assert.Equal(t, "0000-01-01 00:00:00 +0000 UTC", tt.String())
		tt, _ = Parse("0001-01-01 00:00:00")
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", tt.String())
		_, err := Parse("0001-00-00 00:00:00")
		assert.Error(t, err)
	})
}
