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
