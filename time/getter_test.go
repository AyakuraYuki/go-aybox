package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDateTimeMilli(t *testing.T) {
	tt, _ := time.Parse("2006-01-02 15:04:05.999", "2020-02-02 14:03:04.123")
	year, month, day, hour, minute, second, milli := DateTimeMilli(tt)
	assert.Equal(t, 2020, year)
	assert.Equal(t, time.February, month)
	assert.Equal(t, 2, day)
	assert.Equal(t, 14, hour)
	assert.Equal(t, 3, minute)
	assert.Equal(t, 4, second)
	assert.Equal(t, 123, milli)
}

func TestDateTimeMicro(t *testing.T) {
	tt, _ := time.Parse("2006-01-02 15:04:05.999", "2020-02-02 14:03:04.123456")
	year, month, day, hour, minute, second, micro := DateTimeMicro(tt)
	assert.Equal(t, 2020, year)
	assert.Equal(t, time.February, month)
	assert.Equal(t, 2, day)
	assert.Equal(t, 14, hour)
	assert.Equal(t, 3, minute)
	assert.Equal(t, 4, second)
	assert.Equal(t, 123456, micro)
}

func TestDateTimeNano(t *testing.T) {
	tt, _ := time.Parse("2006-01-02 15:04:05.999", "2020-02-02 14:03:04.123456789")
	year, month, day, hour, minute, second, nano := DateTimeNano(tt)
	assert.Equal(t, 2020, year)
	assert.Equal(t, time.February, month)
	assert.Equal(t, 2, day)
	assert.Equal(t, 14, hour)
	assert.Equal(t, 3, minute)
	assert.Equal(t, 4, second)
	assert.Equal(t, 123456789, nano)
}

func TestZoneOffset(t *testing.T) {
	assert.Zero(t, ZoneOffset(NewDate(2020, time.February, 1, time.UTC)))
	assert.Zero(t, ZoneOffset(time.Time{}))
	assert.Equal(t, 28800, ZoneOffset(NewDate(2020, time.February, 1, testCST)))
}
