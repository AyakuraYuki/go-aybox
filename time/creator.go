package time

import "time"

func CreateFromTimestamp(timestamp int64, timezone ...string) (time.Time, error) {
	var loc *time.Location
	var err error
	if len(timezone) > 0 {
		loc, err = loadLocationByTimezone(timezone[0])
		if err != nil {
			return time.Time{}, err
		}
	}
	t := time.Unix(timestamp, 0)
	if loc != nil {
		t = t.In(loc)
	}
	return t, nil
}

func CreateFromTimestampMilli(timestampMilli int64, timezone ...string) (time.Time, error) {
	var loc *time.Location
	var err error
	if len(timezone) > 0 {
		loc, err = loadLocationByTimezone(timezone[0])
		if err != nil {
			return time.Time{}, err
		}
	}
	t := time.UnixMilli(timestampMilli)
	if loc != nil {
		t = t.In(loc)
	}
	return t, nil
}

func CreateFromTimestampMicro(timestampMicro int64, timezone ...string) (time.Time, error) {
	var loc *time.Location
	var err error
	if len(timezone) > 0 {
		loc, err = loadLocationByTimezone(timezone[0])
		if err != nil {
			return time.Time{}, err
		}
	}
	t := time.UnixMicro(timestampMicro)
	if loc != nil {
		t = t.In(loc)
	}
	return t, nil
}

func CreateFromTimestampNano(timestampNano int64, timezone ...string) (time.Time, error) {
	var loc *time.Location
	var err error
	if len(timezone) > 0 {
		loc, err = loadLocationByTimezone(timezone[0])
		if err != nil {
			return time.Time{}, err
		}
	}
	t := time.Unix(timestampNano/1e9, timestampNano%1e9)
	if loc != nil {
		t = t.In(loc)
	}
	return t, nil
}
