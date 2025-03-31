package b64

import (
	"testing"
)

func TestTryDecodeString(t *testing.T) {
	tests := []struct {
		s  string
		ok bool
	}{
		{
			s:  "YWN0aW9uPXJjJm1pZD0yMDAwMDAwMCZlaWQ9aGpBeHhrR2ExTA==",
			ok: true,
		},
		{
			s:  "YWN0aW9uPXJjJm1pZD0yMDAwMDAwMCZlaWQ9aGpBeHhrR2ExTA",
			ok: true,
		},
		{
			s:  "YWN0aW9uPXJjJm1pZD0xJmVpZD1oakFKWldZRnk1",
			ok: true,
		},
		{
			s:  "AUYTOXZX)&^AUGSO",
			ok: false,
		},
		{
			s:  "FOOBARXYZ)(*&^*(&^",
			ok: false,
		},
	}

	for _, tt := range tests {
		bs, ok, _ := TryDecodeString(tt.s)
		if ok != tt.ok {
			t.Errorf("TryDecodeString(%q) = %v, want: %t, got: %t", tt.s, bs, tt.ok, ok)
		}
	}
}
