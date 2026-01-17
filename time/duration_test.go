package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseClockDuration(t *testing.T) {
	tests := []struct {
		s    string
		want time.Duration
	}{
		{"00:00:00.000", 0},
		{"00:01:10.123", time.Duration(70.123 * float64(time.Second))},
		{"00:00:01.456", time.Duration(1.456 * float64(time.Second))},
		{"", 0},
		{"00:01.123", 0},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got, _ := ParseClockDuration(tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseFlexibleDuration(t *testing.T) {
	tests := []struct {
		s    string
		want time.Duration
	}{
		{"00:00:00.000", 0},
		{"00:01:10.123", time.Duration(70.123 * float64(time.Second))},
		{"00:00:01.456", time.Duration(1.456 * float64(time.Second))},
		{"", 0},
		{"00:01", time.Second},
		{"00:01.123", time.Duration(1.123 * float64(time.Second))},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got, _ := ParseFlexibleDuration(tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}
