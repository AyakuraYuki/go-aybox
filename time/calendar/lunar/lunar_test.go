package lunar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testLoc, _ = time.LoadLocation("PRC")

func TestMaxValue(t *testing.T) {
	assert.Equal(t, "二一零零年腊月廿一", MaxValue().ToDateString())
}

func TestMinValue(t *testing.T) {
	assert.Equal(t, "一九零零年正月初一", MinValue().ToDateString())
}

func TestFromStdTime(t *testing.T) {
	t.Run("zero time", func(t *testing.T) {
		assert.Empty(t, FromStdTime(time.Time{}).String())
		assert.Empty(t, FromStdTime(time.Time{}.In(testLoc)).String())
	})

	t.Run("valid time", func(t *testing.T) {
		// 特殊边界
		assert.Equal(t, "2020-04-01", FromStdTime(time.Date(2020, 5, 23, 0, 0, 0, 0, testLoc)).String())
		assert.Equal(t, "2020-05-01", FromStdTime(time.Date(2020, 6, 21, 0, 0, 0, 0, testLoc)).String())

		assert.Equal(t, "2020-06-16", FromStdTime(time.Date(2020, 8, 5, 0, 0, 0, 0, testLoc)).String())
		assert.Equal(t, "2023-02-11", FromStdTime(time.Date(2023, 3, 2, 0, 0, 0, 0, testLoc)).String())
		assert.Equal(t, "2023-02-11", FromStdTime(time.Date(2023, 4, 1, 0, 0, 0, 0, testLoc)).String())
	})
}

func TestLunar_Gregorian(t *testing.T) {
	t.Run("invalid lunar", func(t *testing.T) {
		assert.Empty(t, new(Lunar).ToGregorian().String())
		assert.Empty(t, FromLunar(1800, 1, 1, false).ToGregorian().String())
	})

	t.Run("invalid timezone", func(t *testing.T) {
		assert.Empty(t, FromLunar(2023, 2, 11, false).ToGregorian("xxx").String())
		assert.Empty(t, FromLunar(3200, 1, 1, true).ToGregorian("xxx").String())
	})

	t.Run("without timezone", func(t *testing.T) {
		assert.Equal(t, "2023-03-01 16:00:00 +0000 UTC", FromLunar(2023, 2, 11, false).ToGregorian().String())
		assert.Equal(t, "2023-03-31 16:00:00 +0000 UTC", FromLunar(2023, 2, 11, true).ToGregorian().String())
	})

	t.Run("with timezone", func(t *testing.T) {
		assert.Equal(t, "2023-03-02 00:00:00 +0800 CST", FromLunar(2023, 2, 11, false).ToGregorian("PRC").String())
		assert.Equal(t, "2023-04-01 00:00:00 +0800 CST", FromLunar(2023, 2, 11, true).ToGregorian("PRC").String())
	})
}

func TestLunar_Animal(t *testing.T) {
	t.Run("invalid time", func(t *testing.T) {
		assert.Empty(t, new(Lunar).Animal())
		assert.Empty(t, FromLunar(1800, 1, 1, false).Animal())
	})

	t.Run("valid time", func(t *testing.T) {
		assert.Equal(t, "虎", FromStdTime(time.Date(2010, 8, 5, 0, 0, 0, 0, testLoc)).Animal())
		assert.Equal(t, "兔", FromStdTime(time.Date(2011, 8, 5, 0, 0, 0, 0, testLoc)).Animal())
		assert.Equal(t, "龙", FromStdTime(time.Date(2012, 8, 5, 0, 0, 0, 0, testLoc)).Animal())
		assert.Equal(t, "蛇", FromStdTime(time.Date(2013, 8, 5, 0, 0, 0, 0, testLoc)).Animal())
		assert.Equal(t, "马", FromStdTime(time.Date(2014, 8, 5, 0, 0, 0, 0, testLoc)).Animal())
		assert.Equal(t, "羊", FromStdTime(time.Date(2015, 8, 5, 0, 0, 0, 0, testLoc)).Animal())
		assert.Equal(t, "猴", FromStdTime(time.Date(2016, 8, 5, 0, 0, 0, 0, testLoc)).Animal())
		assert.Equal(t, "鸡", FromStdTime(time.Date(2017, 8, 5, 0, 0, 0, 0, testLoc)).Animal())
		assert.Equal(t, "狗", FromStdTime(time.Date(2018, 8, 5, 0, 0, 0, 0, testLoc)).Animal())
		assert.Equal(t, "猪", FromStdTime(time.Date(2019, 8, 5, 0, 0, 0, 0, testLoc)).Animal())
		assert.Equal(t, "鼠", FromStdTime(time.Date(2020, 8, 5, 0, 0, 0, 0, testLoc)).Animal())
		assert.Equal(t, "牛", FromStdTime(time.Date(2021, 8, 5, 0, 0, 0, 0, testLoc)).Animal())
	})
}

func TestLunar_Festival(t *testing.T) {
	t.Run("invalid time", func(t *testing.T) {
		assert.Empty(t, new(Lunar).Festival())
		assert.Empty(t, FromLunar(1800, 1, 1, false).Festival())
	})

	t.Run("valid time", func(t *testing.T) {
		assert.Equal(t, "春节", FromStdTime(time.Date(2021, 2, 12, 0, 0, 0, 0, testLoc)).Festival())
		assert.Equal(t, "元宵节", FromStdTime(time.Date(2021, 2, 26, 0, 0, 0, 0, testLoc)).Festival())
		assert.Equal(t, "端午节", FromStdTime(time.Date(2021, 6, 14, 0, 0, 0, 0, testLoc)).Festival())
		assert.Equal(t, "七夕节", FromStdTime(time.Date(2021, 8, 14, 0, 0, 0, 0, testLoc)).Festival())
		assert.Equal(t, "中元节", FromStdTime(time.Date(2021, 8, 22, 0, 0, 0, 0, testLoc)).Festival())
		assert.Equal(t, "中秋节", FromStdTime(time.Date(2021, 9, 21, 0, 0, 0, 0, testLoc)).Festival())
		assert.Equal(t, "重阳节", FromStdTime(time.Date(2021, 10, 14, 0, 0, 0, 0, testLoc)).Festival())
		assert.Equal(t, "寒衣节", FromStdTime(time.Date(2021, 11, 5, 0, 0, 0, 0, testLoc)).Festival())
		assert.Equal(t, "下元节", FromStdTime(time.Date(2021, 11, 19, 0, 0, 0, 0, testLoc)).Festival())
		assert.Equal(t, "腊八节", FromStdTime(time.Date(2022, 1, 10, 0, 0, 0, 0, testLoc)).Festival())
	})
}

func TestLunar_Year(t *testing.T) {
	t.Run("invalid time", func(t *testing.T) {
		assert.Empty(t, new(Lunar).Year())
		assert.Empty(t, FromLunar(1800, 1, 1, false).Year())
	})

	t.Run("valid time", func(t *testing.T) {
		l := FromStdTime(time.Date(2020, 8, 5, 0, 0, 0, 0, testLoc))
		assert.Equal(t, "2020-06-16", l.String())
		assert.Equal(t, 2020, l.Year())
	})
}

func TestLunar_Month(t *testing.T) {
	t.Run("invalid time", func(t *testing.T) {
		assert.Empty(t, new(Lunar).Month())
		assert.Empty(t, FromLunar(1800, 1, 1, false).Month())
	})

	t.Run("valid time", func(t *testing.T) {
		l := FromStdTime(time.Date(2020, 8, 5, 0, 0, 0, 0, testLoc))
		assert.Equal(t, "2020-06-16", l.String())
		assert.Equal(t, 6, l.Month())
	})
}

func TestLunar_LeapMonth(t *testing.T) {
	t.Run("invalid time", func(t *testing.T) {
		assert.Empty(t, new(Lunar).LeapMonth())
		assert.Empty(t, FromLunar(1800, 1, 1, false).LeapMonth())
	})
	t.Run("valid time", func(t *testing.T) {
		assert.Equal(t, 4, FromStdTime(time.Date(2020, 8, 5, 0, 0, 0, 0, testLoc)).LeapMonth())
		assert.Equal(t, 2, FromStdTime(time.Date(2023, 3, 2, 0, 0, 0, 0, testLoc)).LeapMonth())
		assert.Equal(t, 6, FromStdTime(time.Date(2025, 10, 7, 0, 0, 0, 0, testLoc)).LeapMonth())
	})
}

func TestLunar_IsLeapMonth(t *testing.T) {
	t.Run("invalid time", func(t *testing.T) {
		assert.False(t, new(Lunar).IsLeapMonth())
		assert.False(t, FromLunar(1800, 1, 1, false).IsLeapMonth())
	})
	t.Run("valid time", func(t *testing.T) {
		l1 := FromStdTime(time.Date(2020, 8, 5, 0, 0, 0, 0, testLoc))
		assert.Equal(t, 4, l1.LeapMonth())
		assert.False(t, l1.IsLeapMonth())

		l2 := FromStdTime(time.Date(2023, 4, 1, 0, 0, 0, 0, testLoc))
		assert.Equal(t, 2, l2.LeapMonth())
		assert.True(t, l2.IsLeapMonth())
	})
}
