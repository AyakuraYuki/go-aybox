package time

import (
	"fmt"
	"time"
)

func ExampleFriendlyDuration() {
	tests := []int64{
		int64((3*time.Hour + 20*time.Minute).Seconds()),
		300,
		114514,
		int64((15*24*time.Hour + 45*time.Minute + 18*time.Second + 100*time.Millisecond).Seconds()),
		int64((18*24*time.Hour + 2*time.Hour + 17*time.Minute + 99*time.Second).Seconds()),
	}
	for _, tt := range tests {
		fmt.Println(FriendlyDuration(tt))
	}

	// Output:
	// 03:20:00
	// 00:05:00
	// 31:48:34
	// 360:45:18
	// 434:18:39
}

func ExampleExtractDuration() {
	tests := []time.Duration{
		12 * time.Minute,
		3*time.Minute + 12*time.Second,
		1*time.Hour + 2*time.Minute + 3*time.Second,
		2*24*time.Hour + 3*time.Minute + 4*time.Second + 999*time.Millisecond,
	}
	for _, dur := range tests {
		fmt.Println(ExtractDuration(int64(dur.Seconds())))
	}

	// Output:
	// 0 12 0
	// 0 3 12
	// 1 2 3
	// 48 3 4
}

func ExampleExtractDurationWithDays() {
	tests := []time.Duration{
		12 * time.Minute,
		3*time.Minute + 12*time.Second,
		1*time.Hour + 2*time.Minute + 3*time.Second,
		2*24*time.Hour + 3*time.Minute + 4*time.Second + 999*time.Millisecond,
		3*24*time.Hour + 4*time.Hour + 5*time.Minute + 6*time.Second + 999*time.Millisecond,
	}
	for _, dur := range tests {
		fmt.Println(ExtractDurationWithDays(int64(dur.Seconds())))
	}

	// Output:
	// 0 0 12 0
	// 0 0 3 12
	// 0 1 2 3
	// 2 0 3 4
	// 3 4 5 6
}
