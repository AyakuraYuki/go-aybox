package attributes

import (
	"fmt"
	"strconv"
)

// About Attributes
//
// 1. The `attribute` is a bitmask.
//
// 2. The `position` is defined as a human-readable permission identifier,
//    using 1-based numbering (i.e., starts at 1). This means position values
//    should not be treated like zero-based list indices.
//
// 3. When performing bitwise operations, the Attribute system automatically
//    converts positions to zero-based indexes internally. Therefore:
//    - Values ≤ 0 are explicitly rejected for position parameters
//    - Client code should always use the documented 1-based positions

// ----------------------------------------------------------------------------

// SetBit modifies a specific bit position in a 64-bit attribute to enable (1)
// or disable (0) a feature flag.
//
// The 1-base position numbering provides human readability:
//
//   - Position 1 = Least Significant Bit (LSB, 2⁰)
//   - Position 64 = Most Significant Bit (MSB, 2⁶³)
//
// Parameters:
//
//   - attr: Original 64-bit value containing existing flags
//   - position: Bit position (1-64), values <= 0 are safely ignored
//   - flag: true = set bit to 1 (enable), false = set to 0 (disable)
//
// Returns: Modified attribute value, or returns original value unchanged if:
//
//   - position <= 0 (invalid 1-based index)
//   - position > 64 (exceeds uint64 capacity)
//
// Implementation note: Internally converts 1-based position to 0-based bitmask
// offsets.
//
// Example:
//
//	position = 3 -> bitmasks = 1<<2 (binary 0b100)
func SetBit(attr uint64, position int, flag bool) uint64 {
	if position <= 0 || position > 64 {
		return attr
	}
	if flag {
		// set to 1 (enable)
		return attr | uint64(1)<<(position-1)
	} else {
		// set to 0 (disable)
		return attr & ^(uint64(1) << (position - 1))
	}
}

// SetBits bulk-modifies multiple bit positions in a 64-bit attribute with
// uniform enable/disable state.
//
// This batch operation is equivalent to iteratively applying SetBit to each
// position:
//
//   - Processes positions in slice order
//   - Invalid positions are silently skipped
//   - Empty positions slice returns original value unchanged
//
// Parameters:
//
//   - attr: Original 64-bit value containing existing flags
//   - positions: Bit positions (1-64), duplicate entries are allowed but may
//     degrade performance with large sets
//   - flag: true = set bit to 1 (enable), false = set to 0 (disable)
//
// Returns: Modified attribute value after applying all valid position changes.
//
// Example:
//
//	SetBits(0b0010, []int{3, 5}, true) -> 0b10110 (decimal 22)
//
// Performance note: For large position sets (>100), consider pre-filtering
// positions:
//
//   - Remove duplicates
//   - Validate range before invocation
//   - Sort positions for cache-friendly access (optional)
func SetBits(attr uint64, positions []int, flag bool) uint64 {
	if len(positions) == 0 {
		return attr
	}
	for _, position := range positions {
		attr = SetBit(attr, position, flag)
	}
	return attr
}

// GetBit retrieves the state (0 or 1) of a specific bit position in a 64-bit
// attribute.
//
// The 1-based position numbering matches SetBit's human-readable convention:
//
//   - Position 1 = Least Significant Bit (LSB, 2⁰)
//   - Position 64 = Most Significant Bit (MSB, 2⁶³)
//
// Parameters:
//
//   - attr: 64-bit value containing feature flags
//   - position: Bit position to query (1-64)
//
// Returns:
//
//   - 1 if the bit is enabled (set)
//   - 0 if the bit is disabled (unset), or for invalid positions:
//     position <= 0 (underflow), position > 64 (overflow)
//
// Example:
//
//	GetBit(0b1010, 2) -> 1 // 2nd bit is set
//	GetBit(0b1010, 3) -> 0 // 3rd bit is unset
func GetBit(attr uint64, position int) int {
	if position <= 0 || position > 64 {
		return 0
	}
	res := (attr >> (position - 1)) & 1
	return int(res)
}

// Enabled retrieves the state of a specific bit position in a 64-bit attribute
// in bool form.
//
// The 1-based position numbering matches SetBit's human-readable convention:
//
//   - Position 1 = Least Significant Bit (LSB, 2⁰)
//   - Position 64 = Most Significant Bit (MSB, 2⁶³)
//
// Parameters:
//
//   - attr: 64-bit value containing feature flags
//   - position: Bit position to query (1-64)
//
// Returns:
//
//   - true if the bit is enabled (set)
//   - false if the bit is disabled (unset), or for invalid positions:
//     position <= 0 (underflow), position > 64 (overflow)
//
// Example:
//
//	Enabled(0b1010, 2) -> true // 2nd bit is set
//	Enabled(0b1010, 3) -> false // 3rd bit is unset
func Enabled(attr uint64, position int) bool {
	return GetBit(attr, position) == 1
}

// ToAttr converts a slice of binary flags into a packed 64-bit attribute
// (bitmask).
//
// The slice index directly maps to 1-based bit positions:
//
//   - flags[0] -> position 1 (LSB, 2⁰)
//   - flags[n] -> position n+1 (2ⁿ)
//
// Input rules:
//
//   - Empty slice returns 0
//   - Values != 0 are treated as 1 (enabled)
//   - Values == 0 are treated as 0 (disabled)
//   - Beyond 64 elements are automatically truncated (positions 65+ ignored)
//
// Returns: A uint64 bitmask where each bit state corresponds to the slice's
// flag value at the equivalent position.
//
// Examples:
//
//	ToAttr([]int{0, 1, 0, 1})  -> 0b1010
//	ToAttr(make([]int, 100))   -> 0
//	ToAttr([]int{3, 0, -1, 0}) -> 0b101 (non-zero as enabled)
//
// Performance note: For slices larger than 64 elements, pre-truncate to avoid
// wasted iterations:
//
//	if len(flags) > 64 {
//		flags = flags[:64]
//	}
func ToAttr(flags []int) (attr uint64) {
	attr = 0
	if len(flags) == 0 {
		return attr
	}
	for i := range flags {
		attr = SetBit(attr, i+1, flags[i] != 0)
	}
	return attr
}

// ToFlags reconstructs a binary flag slice from a 64-bit attribute (bitmask),
// capturing bit states from position 1 (LSB) up to the specified length.
//
// Key characteristics:
//
//   - Result slice index [0] corresponds to position 1 (LSB, 2⁰)
//   - Slice index [n] corresponds to position n+1 (2ⁿ)
//   - Positions beyond 64 (MSB, 2⁶³) are always 0 (unset)
//
// Parameters:
//
//   - attr: 64-bit value containing feature flags
//   - length: Number of positions to extract. Handles edge cases: length <= 0
//     returns empty slice, length > 64 excess positions filled with 0 (unset)
//
// Returns: A []int where each element represents the bit state (0 or 1) of the
// corresponding position, preserving the original bitmask's LSB-to-MSB order.
//
// Examples:
//
//	ToFlags(0b1010, 4)  -> [0, 1, 0, 1]               // Positions 1-4
//	ToFlags(10, 8)      -> [0, 1, 0, 1, 0, 0, 0, 0]   // 10 = 0b1010
//	ToFlags(0xFFFF, 20) -> [1, 1, ..., 1, 0, 0, 0, 0] // first 16 bits 1, rest 0
//	ToFlags(123, -5)    -> []                         // invalid length
//
// Security note: Extremely large length values (e.g. > 1e6) may cause memory
// exhaustion. Consider capping length to reasonable ranges:
//
//	if length > 1024 {
//		length = 1024 // Application-specific limit
//	}
func ToFlags(attr uint64, length int) (flags []int) {
	flags = make([]int, 0)
	for position := 1; position <= length; position++ {
		flags = append(flags, GetBit(attr, position))
	}
	return flags
}

// ----------------------------------------------------------------------------

// ToHex converts an attribute into hexadecimal string
func ToHex(attr uint64) string {
	return strconv.FormatUint(attr, 16)
}

// ParseHex converts a hexadecimal string into attribute
func ParseHex(hex string) (uint64, error) {
	attr, err := strconv.ParseUint(hex, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert hex into attr, err: %v", err)
	}
	return attr, nil
}
