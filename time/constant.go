package time

import (
	"bytes"
	"strings"
	"time"
)

const (
	Local = "Local" // 本地时间
	UTC   = "UTC"   // 世界协调时间

	CET  = "CET"  // 欧洲中部标准时间
	EET  = "EET"  // 欧洲东部标准时间
	EST  = "EST"  // 美国东部标准时间
	GMT  = "GMT"  // 格林尼治标准时间
	MET  = "MET"  // 欧洲中部标准时间
	MST  = "MST"  // 美国山地标准时间
	UCT  = "MST"  // 世界协调时间
	WET  = "WET"  // 欧洲西部标准时间
	Zulu = "Zulu" // 世界协调时间

	AsiaShanghai       = "Asia/Shanghai"       // 上海
	AsiaChongqing      = "Asia/Chongqing"      // 重庆
	AsiaHarbin         = "Asia/Harbin"         // 哈尔滨
	AsiaUrumqi         = "Asia/Urumqi"         // 乌鲁木齐
	AsiaHongKong       = "Asia/Hong_Kong"      // 香港
	AsiaMacao          = "Asia/Macao"          // 澳门
	AsiaTaipei         = "Asia/Taipei"         // 台北
	AsiaTokyo          = "Asia/Tokyo"          // 东京
	AsiaHoChiMinh      = "Asia/Ho_Chi_Minh"    // 胡志明
	AsiaHanoi          = "Asia/Hanoi"          // 河内
	AsiaSaigon         = "Asia/Saigon"         // 西贡
	AsiaSeoul          = "Asia/Seoul"          // 首尔
	AsiaPyongyang      = "Asia/Pyongyang"      // 平壤
	AsiaBangkok        = "Asia/Bangkok"        // 曼谷
	AsiaDubai          = "Asia/Dubai"          // 迪拜
	AsiaQatar          = "Asia/Qatar"          // 卡塔尔
	AsiaBangalore      = "Asia/Bangalore"      // 班加罗尔
	AsiaKolkata        = "Asia/Kolkata"        // 加尔各答
	AsiaMumbai         = "Asia/Mumbai"         // 孟买
	AmericaMexicoCity  = "America/Mexico_City" // 墨西哥
	AmericaNewYork     = "America/New_York"    // 纽约
	AmericaLosAngeles  = "America/Los_Angeles" // 洛杉矶
	AmericaChicago     = "America/Chicago"     // 芝加哥
	AmericaSaoPaulo    = "America/Sao_Paulo"   // 圣保罗
	EuropeMoscow       = "Europe/Moscow"       // 莫斯科
	EuropeLondon       = "Europe/London"       // 伦敦
	EuropeBerlin       = "Europe/Berlin"       // 柏林
	EuropeParis        = "Europe/Paris"        // 巴黎
	EuropeRome         = "Europe/Rome"         // 罗马
	AustraliaSydney    = "Australia/Sydney"    // 悉尼
	AustraliaMelbourne = "Australia/Melbourne" // 墨尔本
	AustraliaDarwin    = "Australia/Darwin"    // 达尔文
)

// month constants
// 月份常量
const (
	January   = "January"   // 一月
	February  = "February"  // 二月
	March     = "March"     // 三月
	April     = "April"     // 四月
	May       = "May"       // 五月
	June      = "June"      // 六月
	July      = "July"      // 七月
	August    = "August"    // 八月
	September = "September" // 九月
	October   = "October"   // 十月
	November  = "November"  // 十一月
	December  = "December"  // 十二月
)

// constellation constants
// 星座常量
const (
	Aries       = "Aries"       // 白羊座
	Taurus      = "Taurus"      // 金牛座
	Gemini      = "Gemini"      // 双子座
	Cancer      = "Cancer"      // 巨蟹座
	Leo         = "Leo"         // 狮子座
	Virgo       = "Virgo"       // 处女座
	Libra       = "Libra"       // 天秤座
	Scorpio     = "Scorpio"     // 天蝎座
	Sagittarius = "Sagittarius" // 射手座
	Capricorn   = "Capricorn"   // 摩羯座
	Aquarius    = "Aquarius"    // 水瓶座
	Pisces      = "Pisces"      // 双鱼座
)

// week constants
// 星期常量
const (
	Monday    = "Monday"    // 周一
	Tuesday   = "Tuesday"   // 周二
	Wednesday = "Wednesday" // 周三
	Thursday  = "Thursday"  // 周四
	Friday    = "Friday"    // 周五
	Saturday  = "Saturday"  // 周六
	Sunday    = "Sunday"    // 周日
)

// season constants
// 季节常量
const (
	Spring = "Spring" // 春季
	Summer = "Summer" // 夏季
	Autumn = "Autumn" // 秋季
	Winter = "Winter" // 冬季
)

// number constants
// 数字常量
const (
	YearsPerMillennium = 1000   // 每千年1000年
	YearsPerCentury    = 100    // 每世纪100年
	YearsPerDecade     = 10     // 每十年10年
	QuartersPerYear    = 4      // 每年4个季度
	MonthsPerYear      = 12     // 每年12月
	MonthsPerQuarter   = 3      // 每季度3月
	WeeksPerNormalYear = 52     // 每常规年52周
	weeksPerLongYear   = 53     // 每长年53周
	WeeksPerMonth      = 4      // 每月4周
	DaysPerLeapYear    = 366    // 每闰年366天
	DaysPerNormalYear  = 365    // 每常规年365天
	DaysPerWeek        = 7      // 每周7天
	HoursPerWeek       = 168    // 每周168小时
	HoursPerDay        = 24     // 每天24小时
	MinutesPerDay      = 1440   // 每天1440分钟
	MinutesPerHour     = 60     // 每小时60分钟
	SecondsPerWeek     = 604800 // 每周604800秒
	SecondsPerDay      = 86400  // 每天86400秒
	SecondsPerHour     = 3600   // 每小时3600秒
	SecondsPerMinute   = 60     // 每分钟60秒
)

// layout constants
// 布局模板常量
const (
	AtomLayout     = RFC3339Layout
	ANSICLayout    = time.ANSIC
	CookieLayout   = "Monday, 02-Jan-2006 15:04:05 MST"
	KitchenLayout  = time.Kitchen
	RssLayout      = time.RFC1123Z
	RubyDateLayout = time.RubyDate
	UnixDateLayout = time.UnixDate
	W3cLayout      = RFC3339Layout

	RFC1036Layout      = "Mon, 02 Jan 06 15:04:05 -0700"
	RFC1123Layout      = time.RFC1123
	RFC1123ZLayout     = time.RFC1123Z
	RFC2822Layout      = time.RFC1123Z
	RFC3339Layout      = "2006-01-02T15:04:05Z07:00"
	RFC3339MilliLayout = "2006-01-02T15:04:05.999Z07:00"
	RFC3339MicroLayout = "2006-01-02T15:04:05.999999Z07:00"
	RFC3339NanoLayout  = "2006-01-02T15:04:05.999999999Z07:00"
	RFC7231Layout      = "Mon, 02 Jan 2006 15:04:05 MST"
	RFC822Layout       = time.RFC822
	RFC822ZLayout      = time.RFC822Z
	RFC850Layout       = time.RFC850

	ISO8601Layout      = "2006-01-02T15:04:05-07:00"
	ISO8601MilliLayout = "2006-01-02T15:04:05.999-07:00"
	ISO8601MicroLayout = "2006-01-02T15:04:05.999999-07:00"
	ISO8601NanoLayout  = "2006-01-02T15:04:05.999999999-07:00"

	ISO8601ZuluLayout      = "2006-01-02T15:04:05Z"
	ISO8601ZuluMilliLayout = "2006-01-02T15:04:05.999Z"
	ISO8601ZuluMicroLayout = "2006-01-02T15:04:05.999999Z"
	ISO8601ZuluNanoLayout  = "2006-01-02T15:04:05.999999999Z"

	FormattedDateLayout    = "Jan 2, 2006"
	FormattedDayDateLayout = "Mon, Jan 2, 2006"

	DayDateTimeLayout        = "Mon, Jan 2, 2006 3:04 PM"
	DateTimeLayout           = "2006-01-02 15:04:05"
	DateTimeMilliLayout      = "2006-01-02 15:04:05.999"
	DateTimeMicroLayout      = "2006-01-02 15:04:05.999999"
	DateTimeNanoLayout       = "2006-01-02 15:04:05.999999999"
	ShortDateTimeLayout      = "20060102150405"
	ShortDateTimeMilliLayout = "20060102150405.999"
	ShortDateTimeMicroLayout = "20060102150405.999999"
	ShortDateTimeNanoLayout  = "20060102150405.999999999"

	DateLayout           = "2006-01-02"
	DateMilliLayout      = "2006-01-02.999"
	DateMicroLayout      = "2006-01-02.999999"
	DateNanoLayout       = "2006-01-02.999999999"
	ShortDateLayout      = "20060102"
	ShortDateMilliLayout = "20060102.999"
	ShortDateMicroLayout = "20060102.999999"
	ShortDateNanoLayout  = "20060102.999999999"

	TimeLayout           = "15:04:05"
	TimeMilliLayout      = "15:04:05.999"
	TimeMicroLayout      = "15:04:05.999999"
	TimeNanoLayout       = "15:04:05.999999999"
	ShortTimeLayout      = "150405"
	ShortTimeMilliLayout = "150405.999"
	ShortTimeMicroLayout = "150405.999999"
	ShortTimeNanoLayout  = "150405.999999999"

	TimestampLayout      = "unix"
	TimestampMilliLayout = "unixMilli"
	TimestampMicroLayout = "unixMicro"
	TimestampNanoLayout  = "unixNano"
)

// common formatting symbols
var formats = map[byte]string{
	'd': "02",      // Day:    Day of the month, 2 digits with leading zeros. Eg: 01 to 31.
	'D': "Mon",     // Day:    A textual representation of a day, three letters. Eg: Mon through Sun.
	'j': "2",       // Day:    Day of the month without leading zeros. Eg: 1 to 31.
	'l': "Monday",  // Day:    A full textual representation of the day of the week. Eg: Sunday through Saturday.
	'F': "January", // Month:  A full textual representation of a month, such as January or March. Eg: January through December.
	'm': "01",      // Month:  Numeric representation of a month, with leading zeros. Eg: 01 through 12.
	'M': "Jan",     // Month:  A short textual representation of a month, three letters. Eg: Jan through Dec.
	'n': "1",       // Month:  Numeric representation of a month, without leading zeros. Eg: 1 through 12.
	'Y': "2006",    // Year:   A full numeric representation of a year, 4 digits. Eg: 1999 or 2003.
	'y': "06",      // Year:   A two digit representation of a year. Eg: 99 or 03.
	'a': "pm",      // Time:   Lowercase morning or afternoon sign. Eg: am or pm.
	'A': "PM",      // Time:   Uppercase morning or afternoon sign. Eg: AM or PM.
	'g': "3",       // Time:   12-hour format of an hour without leading zeros. Eg: 1 through 12.
	'h': "03",      // Time:   12-hour format of an hour with leading zeros. Eg: 01 through 12.
	'H': "15",      // Time:   24-hour format of an hour with leading zeros. Eg: 00 through 23.
	'i': "04",      // Time:   Minutes with leading zeros. Eg: 00 to 59.
	's': "05",      // Time:   Seconds with leading zeros. Eg: 00 through 59.
	'O': "-0700",   // Zone:   Difference to Greenwich time (GMT) in hours. Eg: +0200.
	'P': "-07:00",  // Zone:   Difference to Greenwich time (GMT) with colon between hours and minutes. Eg: +02:00.
	'T': "MST",     // Zone:   Timezone abbreviation. Eg: UTC, EST, MDT ...

	'U': "timestamp",      // Timestamp with second. Eg: 1699677240.
	'V': "timestampMilli", // TimestampMilli with second. Eg: 1596604455666.
	'X': "timestampMicro", // TimestampMicro with second. Eg: 1596604455666666.
	'Z': "timestampNano",  // TimestampNano with second. Eg: 1596604455666666666.
}

// common layouts
var defaultLayouts = []string{
	DateTimeLayout, DateLayout, TimeLayout,
	ISO8601Layout, DayDateTimeLayout, ISO8601NanoLayout,
	DateTimeNanoLayout, ShortDateTimeLayout, ShortDateTimeNanoLayout,
	DateLayout, DateNanoLayout, ShortDateLayout, ShortDateNanoLayout,
	TimeMicroLayout, TimeMilliLayout, TimeNanoLayout,
	RFC822Layout, RFC822ZLayout, RFC850Layout, RFC1123Layout, RFC1123ZLayout, RFC3339Layout, RFC3339NanoLayout, RFC1036Layout, RFC7231Layout,
	KitchenLayout,
	CookieLayout,
	ANSICLayout,
	UnixDateLayout,
	RubyDateLayout,
	"2006",
	"2006-1", "2006-1-2", "2006-1-2 15", "2006-1-2 15:4", "2006-1-2 15:4:5", "2006-1-2 15:4:5.999999999",
	"2006.1", "2006.1.2", "2006.1.2 15", "2006.1.2 15:4", "2006.1.2 15:4:5", "2006.1.2 15:4:5.999999999",
	"2006/1", "2006/1/2", "2006/1/2 15", "2006/1/2 15:4", "2006/1/2 15:4:5", "2006/1/2 15:4:5.999999999",
	"2006-01-02 15:04:05 -0700 MST",
	"2006-01-02 15:04:05PM MST", "2006-01-02 15:04:05.999999999PM MST", "2006-1-2 15:4:5PM MST", "2006-1-2 15:4:5.999999999PM MST",
	"2006-01-02 15:04:05 PM MST", "2006-01-02 15:04:05.999999999 PM MST", "2006-1-2 15:4:5 PM MST", "2006-1-2 15:4:5.999999999 PM MST",
	"1/2/2006", "1/2/2006 15", "1/2/2006 15:4", "1/2/2006 15:4:5", "1/2/2006 15:4:5.999999999",
	"2006-1-2 15:4:5 -0700 MST", "2006-1-2 15:4:5.999999999 -0700 MST", "2006-1-2 15:04:05 -0700 MST", "2006-1-2 15:04:05.999999999 -0700 MST",
	"2006-01-02T15:04:05", "2006-01-02T15:04:05.999999999", "2006-1-2T3:4:5", "2006-1-2T3:4:5.999999999",
	"2006-01-02T15:04:05Z07", "2006-01-02T15:04:05.999999999Z07", "2006-1-2T15:4:5Z07", "2006-1-2T15:4:5.999999999Z07",
	"2006-01-02T15:04:05Z07:00", "2006-01-02T15:04:05.999999999Z07:00", "2006-1-2T15:4:5Z07:00", "2006-1-2T15:4:5.999999999Z07:00",
	"2006-01-02T15:04:05-07:00", "2006-01-02T15:04:05.999999999-07:00", "2006-1-2T15:4:5-07:00", "2006-1-2T15:4:5.999999999-07:00",
	"2006-01-02T15:04:05-0700", "2006-01-02T15:04:05.999999999-0700", "2006-1-2T3:4:5-0700", "2006-1-2T3:4:5.999999999-0700",
	"20060102150405-07:00", "20060102150405.999999999-07:00",
	"20060102150405Z07", "20060102150405.999999999Z07",
	"20060102150405Z07:00", "20060102150405.999999999Z07:00",
}

// converts format to layout
func format2layout(format string) string {
	buffer := bytes.NewBuffer(nil)
	for i := 0; i < len(format); i++ {
		if layout, ok := formats[format[i]]; ok {
			buffer.WriteString(layout)
		} else {
			switch format[i] {
			case '\\': // raw output, no parse
				buffer.WriteByte(format[i+1])
				i++
				continue
			default:
				buffer.WriteByte(format[i])
			}
		}
	}
	return buffer.String()
}

// gets a [time.Location] instance by a timezone string
func loadLocationByTimezone(timezone string) (*time.Location, error) {
	tz := strings.TrimSpace(timezone)
	if tz == "" {
		return nil, emptyTimezoneError()
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		err = invalidTimezoneError(tz)
	}
	return loc, err
}

// gets absolute value
func getAbsValue(value int64) int64 {
	return (value ^ value>>31) - value>>31
}
