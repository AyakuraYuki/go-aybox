package time

import "fmt"

// returns a failed parse error.
// 解析失败错误
var failedParseError = func(value string) error {
	return fmt.Errorf("cannot parse %q as time.Time", value)
}

// returns a invalid timestamp error.
// 无效的时间戳错误
var invalidTimestampError = func(value string) error {
	return fmt.Errorf("invalid timestamp %s", value)
}

// returns a nil location error.
// 无效的位置错误
var nilLocationError = func() error {
	return fmt.Errorf("location cannot be nil")
}

// returns an nil language error.
// 无效的语言错误
var nilLanguageError = func() error {
	return fmt.Errorf("language cannot be nil")
}

// returns an empty timezone error.
// 空的时区错误
var emptyTimezoneError = func() error {
	return fmt.Errorf("timezone cannot be empty")
}

// returns an invalid timezone error.
// 无效的时区错误
var invalidTimezoneError = func(timezone string) error {
	return fmt.Errorf(`invalid timezone %q, please see the file "$GOROOT/lib/time/zoneinfo.zip" for all valid timezones`, timezone)
}

// returns an empty duration error.
// 空的时长错误
var emptyDurationError = func() error {
	return fmt.Errorf("duration cannot be empty")
}

// returns an invalid duration error.
// 无效的时长错误
var invalidDurationError = func(duration string) error {
	return fmt.Errorf("invalid duration %q", duration)
}

// returns an empty layout error.
// 空的布局模板错误
var emptyLayoutError = func() error {
	return fmt.Errorf("layout cannot be empty")
}

// returns an invalid layout error.
// 无效的布局模板错误
var invalidLayoutError = func(value, layout string) error {
	return fmt.Errorf("cannot parse string %q as time.Time by layout %q", value, layout)
}

// returns an empty format error.
// 空的格式模板错误
var emptyFormatError = func() error {
	return fmt.Errorf("format cannot be empty")
}

// returns an invalid format error.
// 无效的格式模板错误
var invalidFormatError = func(value, format string) error {
	return fmt.Errorf("cannot parse string %q as time.Time by format %q", value, format)
}

// returns an empty week starts day error.
// 空的周起始日期错误
var emptyWeekStartsDayError = func() error {
	return fmt.Errorf("week start day cannot be empty")
}

// returns an invalid week starts at day error.
// 无效的周起始日期错误
var invalidWeekStartsAtError = func(day string) error {
	return fmt.Errorf("invalid week starts at day %s", day)
}
