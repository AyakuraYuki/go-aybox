package str

import (
	"reflect"
	"testing"
)

func TestTrimLower(t *testing.T) {
	input := " ABC "
	want := "abc"
	output := TrimLower(input)
	if output != want {
		t.Errorf("TrimLower(%q) = %q, want %q", input, output, want)
	}

	outputPtr := TrimLower(&input)
	if outputPtr == nil || *outputPtr != want {
		t.Errorf("TrimLower(%q) as ptr = %v, want %q", input, outputPtr, want)
	}

	inputSlice := []string{"ABC", "Bar"}
	wantSlice := []string{"abc", "bar"}
	outputSlice := TrimLower(inputSlice)
	if !reflect.DeepEqual(outputSlice, wantSlice) {
		t.Errorf("TrimLower(%v) = %v, want %v", inputSlice, outputSlice, wantSlice)
	}
}

func TestTrimUpper(t *testing.T) {
	input := " abc "
	want := "ABC"
	output := TrimUpper(input)
	if output != want {
		t.Errorf("TrimUpper(%q) = %q, want %q", input, output, want)
	}

	outputPtr := TrimUpper(&input)
	if outputPtr == nil || *outputPtr != want {
		t.Errorf("TrimUpper(%q) as ptr = %v, want %q", input, outputPtr, want)
	}

	inputSlice := []string{"abc", "Bar"}
	wantSlice := []string{"ABC", "BAR"}
	outputSlice := TrimUpper(inputSlice)
	if !reflect.DeepEqual(outputSlice, wantSlice) {
		t.Errorf("TrimUpper(%v) = %v, want %v", inputSlice, outputSlice, wantSlice)
	}
}
