package random

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandPassword(t *testing.T) {
	assert.True(t, regexp.MustCompile(`(\S+){16}`).MatchString(RandPassword(16)))
}

func BenchmarkRandPassword_16(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = RandPassword(16)
	}
}

func BenchmarkRandPasswordParallel_16(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = RandPassword(16)
		}
	})
}
