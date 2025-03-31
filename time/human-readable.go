package time

import (
	"fmt"
	"time"
)

// FriendlyDuration returns a duration string
func FriendlyDuration(durationInSeconds int64) (desc string) {
	duration := time.Duration(durationInSeconds) * time.Second

	hours := int64(duration.Hours())
	duration -= time.Duration(hours) * time.Hour

	minutes := int64(duration.Minutes())
	duration -= time.Duration(minutes) * time.Minute

	seconds := int64(duration.Seconds())
	duration -= time.Duration(seconds) * time.Second

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

// ExtractDuration extracts a given duration seconds into hours, minutes and
// seconds
func ExtractDuration(durationInSeconds int64) (hours, minutes, seconds int64) {
	duration := time.Duration(durationInSeconds) * time.Second

	hours = int64(duration.Hours())
	duration -= time.Duration(hours) * time.Hour

	minutes = int64(duration.Minutes())
	duration -= time.Duration(minutes) * time.Minute

	seconds = int64(duration.Seconds())
	duration -= time.Duration(seconds) * time.Second

	return
}

// ExtractDurationWithDays extracts a given duration seconds into days,
// hours (in day), minutes and seconds
func ExtractDurationWithDays(durationInSeconds int64) (days, hours, minutes, seconds int64) {
	duration := time.Duration(durationInSeconds) * time.Second

	hours = int64(duration.Hours())
	duration -= time.Duration(hours) * time.Hour
	days = hours / 24
	hours = hours - (days * 24)

	minutes = int64(duration.Minutes())
	duration -= time.Duration(minutes) * time.Minute

	seconds = int64(duration.Seconds())
	duration -= time.Duration(seconds) * time.Second

	return
}
