package time

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseClockDuration accepts a clock time string which formatted in
// `01:02:03.456` then parses into [time.Duration].
func ParseClockDuration(s string) (dur time.Duration, err error) {
	var (
		h, m int
		sec  float64
	)

	_, err = fmt.Sscanf(s, "%d:%d:%f", &h, &m, &sec)
	if err != nil {
		return 0, err
	}

	dur = time.Duration(h)*time.Hour +
		time.Duration(m)*time.Minute +
		time.Duration(sec*float64(time.Second))

	return dur, nil
}

// ParseFlexibleDuration accepts a clock time string which formatted in
// `[01:]02:03[.456]` then parses into [time.Duration]. The hour
// part and millisecond part are optional.
func ParseFlexibleDuration(s string) (dur time.Duration, err error) {
	parts := strings.Split(s, ":")
	if len(parts) < 2 || len(parts) > 3 {
		return 0, fmt.Errorf("malformed clock duration: %s", s)
	}

	var (
		h, m int
		sec  float64
	)

	if len(parts) == 3 {
		// hour
		h, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0, fmt.Errorf("malformed hour: %w", err)
		}
		parts = parts[1:]
	}

	// minute
	m, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("malformed minute: %w", err)
	}

	sec, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, fmt.Errorf("malformed second: %w", err)
	}

	return time.Duration(h)*time.Hour +
		time.Duration(m)*time.Minute +
		time.Duration(sec*float64(time.Second)), nil
}
