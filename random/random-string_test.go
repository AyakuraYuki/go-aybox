package random

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandString(t *testing.T) {
	assert.True(t, regexp.MustCompile(`[a-z]{16}`).MatchString(RandString(16, Lower)))
	assert.False(t, regexp.MustCompile(`[A-Z]{16}`).MatchString(RandString(16, Lower)))
	assert.True(t, regexp.MustCompile(`[A-Z]{16}`).MatchString(RandString(16, Upper)))
	assert.False(t, regexp.MustCompile(`[a-z]{16}`).MatchString(RandString(16, Upper)))
	assert.True(t, regexp.MustCompile(`[a-zA-Z0-9]{16}`).MatchString(RandString(16, LowerUpperDigit)))
}

func BenchmarkRandString_16(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		RandString(16, LowerUpperDigit)
	}
}

func BenchmarkRandString_30(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		RandString(30, LowerUpperDigit)
	}
}

func BenchmarkRandString_50(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		RandString(50, LowerUpperDigit)
	}
}
