package time

import (
	"strconv"
	"time"
)

// Parse parses a standard time string as a time.Time instance.
func Parse(value string, timezone ...string) (time.Time, error) {
	if value == "" {
		return time.Time{}, failedParseError("empty string")
	}

	var loc *time.Location
	var err error
	if len(timezone) > 0 {
		loc, err = loadLocationByTimezone(timezone[0])
	}
	if err != nil {
		return time.Time{}, err
	}

	now := time.Now()
	if loc != nil {
		now = now.In(loc)
	}

	switch value {
	case "now", "Now", "NOW", "today", "Today", "TODAY":
		return now, nil
	case "yesterday", "Yesterday", "YESTERDAY":
		return SubDay(now), nil
	case "tomorrow", "Tomorrow", "TOMORROW":
		return AddDay(now), nil
	}

	for _, layout := range defaultLayouts {
		var tt time.Time
		if loc == nil {
			tt, err = time.Parse(layout, value)
		} else {
			tt, err = time.ParseInLocation(layout, value, loc)
		}
		if err == nil {
			return tt, nil
		}
	}

	err = failedParseError(value)
	return time.Time{}, err
}

// ParseByFormat parses a time string as a time.Time instance by format.
func ParseByFormat(value, format string, timezone ...string) (time.Time, error) {
	if value == "" {
		return time.Time{}, failedParseError("empty string")
	}
	if format == "" {
		return time.Time{}, emptyFormatError()
	}
	t, err := ParseByLayout(value, format2layout(format), timezone...)
	if err != nil {
		err = invalidFormatError(value, format)
	}
	return t, err
}

// ParseWithFormats parses time string with formats as a time.Time instance.
func ParseWithFormats(value string, formats []string, timezone ...string) (time.Time, error) {
	if value == "" {
		return time.Time{}, failedParseError("empty string")
	}
	if len(formats) == 0 {
		return Parse(value, timezone...)
	}
	var l []string
	for _, format := range formats {
		l = append(l, format2layout(format))
	}
	return ParseWithLayouts(value, l, timezone...)
}

// ParseByLayout parses a time string as a time.Time instance by layout.
func ParseByLayout(value, layout string, timezone ...string) (time.Time, error) {
	if value == "" {
		return time.Time{}, failedParseError("empty string")
	}
	if layout == "" {
		return time.Time{}, emptyLayoutError()
	}

	var loc *time.Location
	var err error
	if len(timezone) > 0 {
		loc, err = loadLocationByTimezone(timezone[0])
	}
	if err != nil {
		return time.Time{}, err
	}

	var ts int64
	switch layout {
	case TimestampLayout:
		ts, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			err = invalidTimestampError(value)
			return time.Time{}, err
		}
		return CreateFromTimestamp(ts, timezone...)

	case TimestampMilliLayout:
		ts, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			err = invalidTimestampError(value)
			return time.Time{}, err
		}
		return CreateFromTimestampMilli(ts, timezone...)

	case TimestampMicroLayout:
		ts, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			err = invalidTimestampError(value)
			return time.Time{}, err
		}
		return CreateFromTimestampMicro(ts, timezone...)

	case TimestampNanoLayout:
		ts, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			err = invalidTimestampError(value)
			return time.Time{}, err
		}
		return CreateFromTimestampNano(ts, timezone...)

	}

	var tt time.Time
	if loc == nil {
		tt, err = time.Parse(layout, value)
	} else {
		tt, err = time.ParseInLocation(layout, value, loc)
	}
	if err == nil {
		return tt, nil
	}

	err = invalidLayoutError(value, layout)
	return time.Time{}, err
}

// ParseWithLayouts parses time string with layouts as a time.Time instance.
func ParseWithLayouts(value string, layouts []string, timezone ...string) (time.Time, error) {
	if value == "" {
		return time.Time{}, failedParseError("empty string")
	}
	if len(layouts) == 0 {
		return Parse(value, timezone...)
	}
	var tt time.Time
	var err error
	for _, layout := range layouts {
		tt, err = ParseByLayout(value, layout, timezone...)
		if err == nil {
			return tt, nil
		}
	}
	err = failedParseError(value)
	return time.Time{}, err
}
