package attributes

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSetBit(t *testing.T) {
	var attr, want uint64

	// case: ignore invalid position
	attr, want = 0b1, 0b1
	if get := SetBit(attr, 0, true); get != want {
		t.Errorf("get = %b, want %b because the position is invalid", get, want)
	}
	if get := SetBit(attr, 65, true); get != want {
		t.Errorf("get = %b, want %b because the position is invalid", get, want)
	}

	// case: set enabled
	attr, want = 0b0, 0b1
	if get := SetBit(attr, 1, true); get != want {
		t.Errorf("get = %b, want %b", get, want)
	}

	// case: set disabled
	attr, want = 0b1001111, 0b1001011
	if get := SetBit(attr, 3, false); get != want {
		t.Errorf("get = %b, want %b", get, want)
	}

	// case: enable multiple positions
	attr, want = 0b1001111, 0b1100100001001111
	attr = SetBit(attr, 3, true)  // re-enable exist position 3
	attr = SetBit(attr, 12, true) // enable position 12
	attr = SetBit(attr, 15, true) // enable position 15
	attr = SetBit(attr, 16, true) // enable position 16
	if attr != want {
		t.Errorf("attr = %b, want %b", attr, want)
	}

	// case: disable multiple positions
	attr, want = 0b1001111, 0b1011
	attr = SetBit(attr, 3, false) // disable position 3
	attr = SetBit(attr, 7, false) // disable position 7
	attr = SetBit(attr, 8, false) // disable non-enabled position 8
	if attr != want {
		t.Errorf("attr = %b, want %b", attr, want)
	}
}

func TestSetBits(t *testing.T) {
	var attr, want uint64

	// case: enable multiple positions
	attr, want = 0b0, 0b101010
	if get := SetBits(attr, []int{2, 4, 6}, true); get != want {
		t.Errorf("get = %b, want %b", get, want)
	}

	// case: disable multiple positions
	attr, want = 0b1001111, 0b1000001
	if get := SetBits(attr, []int{2, 3, 4}, false); get != want {
		t.Errorf("get = %b, want %b", get, want)
	}
}

func TestGetBit(t *testing.T) {
	tests := []struct {
		attr     uint64
		position int
		want     int
	}{
		{0b0, 1, 0},
		{0b1, 1, 1},
		{0b1010, 2, 1},
		{0b1010, 3, 0},

		{139, 3, 0},
		{51343, 3, 1},
		{11, 2, 1},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d@%d", tt.position, tt.attr), func(t *testing.T) {
			if get := GetBit(tt.attr, tt.position); get != tt.want {
				t.Errorf("get = %d, want %d", get, tt.want)
			}
		})
	}
}

func TestEnabled(t *testing.T) {
	tests := []struct {
		attr     uint64
		position int
		want     bool
	}{
		{0b0, 1, false},
		{0b1, 1, true},
		{0b1010, 2, true},
		{0b1010, 3, false},

		{139, 3, false},
		{51343, 3, true},
		{11, 2, true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d@%d", tt.position, tt.attr), func(t *testing.T) {
			if get := Enabled(tt.attr, tt.position); get != tt.want {
				t.Errorf("get = %t, want %t", get, tt.want)
			}
		})
	}
}

func TestToAttr(t *testing.T) {
	tests := []struct {
		flags []int
		want  uint64
	}{
		{[]int{0, 1, 0, 1}, 0b1010},
		{make([]int, 100), 0},
		{[]int{3, 0, -1, 0}, 0b101},
	}

	for _, tt := range tests {
		if get := ToAttr(tt.flags); get != tt.want {
			t.Errorf("get = %d, want %d", get, tt.want)
		}
	}
}

func TestToFlag(t *testing.T) {
	tests := []struct {
		attr   uint64
		length int
		want   []int
	}{
		{0b0, 1, []int{0}},
		{0b11111, 3, []int{1, 1, 1}},
		{0b10101, 4, []int{1, 0, 1, 0}},
		{0b10101, 6, []int{1, 0, 1, 0, 1, 0}},

		{1109, 15, []int{1, 0, 1, 0, 1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0}},
		{0, 100, make([]int, 100)},

		{0b1010, 4, []int{0, 1, 0, 1}},
		{10, 8, []int{0, 1, 0, 1, 0, 0, 0, 0}},
		{0xFFFF, 20, []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0}},
		{123, -5, []int{}},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d[:%d]", tt.attr, tt.length), func(t *testing.T) {
			get := ToFlags(tt.attr, tt.length)
			if !reflect.DeepEqual(get, tt.want) {
				t.Errorf("get = %v, want %v", get, tt.want)
			}
		})
	}
}
