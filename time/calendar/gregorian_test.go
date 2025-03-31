package calendar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGregorian_String(t *testing.T) {
	t.Run("zero time", func(t *testing.T) {
		assert.Empty(t, Gregorian{}.String())
	})

	t.Run("valid time", func(t *testing.T) {
		g := Gregorian{Time: time.Date(2020, 8, 5, 0, 0, 0, 0, time.UTC)}
		assert.Equal(t, "2020-08-05 00:00:00 +0000 UTC", g.String())
	})
}
