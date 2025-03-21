package aybox

import (
	"testing"
	"time"
)

func TestCalculateRealAge(t *testing.T) {
	tests := []time.Time{
		time.Date(1995, 1, 27, 0, 0, 0, 0, time.UTC),
		time.Date(1995, 9, 14, 0, 0, 0, 0, time.UTC),
		time.Date(1996, 1, 8, 0, 0, 0, 0, time.UTC),
		time.Date(1996, 12, 26, 0, 0, 0, 0, time.UTC),
	}
	for _, birthday := range tests {
		t.Log(birthday.Format(time.DateOnly), CalculateRealAge[int](birthday))
		t.Log(birthday.Format(time.DateOnly), CalculateRealAge[int](birthday, "Asia/Shanghai"), "Asia/Shanghai")
	}
}

func TestCalculateNominalAge(t *testing.T) {
	tests := []time.Time{
		time.Date(1995, 1, 10, 0, 0, 0, 0, time.UTC),
		time.Date(1995, 9, 14, 0, 0, 0, 0, time.UTC),
		time.Date(1996, 1, 8, 0, 0, 0, 0, time.UTC),
		time.Date(1996, 12, 26, 0, 0, 0, 0, time.UTC),
	}
	for _, birthday := range tests {
		t.Log(birthday.Format(time.DateOnly), CalculateNominalAge[int](birthday))
	}
}
