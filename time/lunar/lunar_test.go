package lunar

import (
	aytime "github.com/AyakuraYuki/go-aybox/time"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFromStdTime(t *testing.T) {
	loc, _ := time.LoadLocation("PRC")

	t.Run("zero time", func(t *testing.T) {
		assert.Empty(t, FromStdTime(time.Time{}).String())
		assert.Empty(t, FromStdTime(time.Time{}.In(loc)).String())
	})

	t.Run("valid time", func(t *testing.T) {
		assert.Equal(t, "2020-04-01", FromStdTime(aytime.NewDate(2020, 5, 23, loc)).String())
		assert.Equal(t, "2020-05-01", FromStdTime(aytime.NewDate(2020, 6, 21, loc)).String())

		assert.Equal(t, "2020-06-16", FromStdTime(aytime.NewDate(2020, 8, 5, loc)).String())
		assert.Equal(t, "2023-02-11", FromStdTime(aytime.NewDate(2023, 3, 2, loc)).String())
		assert.Equal(t, "2023-02-11", FromStdTime(aytime.NewDate(2023, 4, 1, loc)).String())
	})
}
