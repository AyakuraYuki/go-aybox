package str

import "strings"

// TrimLower converts the string, the value of a string pointer, and the
// string slice to the trimmed, lower result.
//
// For example, a given string of `" ABC "` will be converted to `"abc"`;
// a given string pointer containing the value `" ABC "` will be converted
// to `"abc"` and results in a new string pointer; a given string slice
// with `["ABC", "Bar"]` will be converted to `["abc", "bar"]`.
func TrimLower[T string | *string | []string](input T) T {
	var result any

	switch s := any(input).(type) {
	case string:
		result = strings.ToLower(strings.TrimSpace(s))

	case *string:
		if s == nil {
			return input
		}
		trimmedLower := strings.ToLower(strings.TrimSpace(*s))
		result = &trimmedLower

	case []string:
		trimmedLower := make([]string, len(s))
		for i, v := range s {
			trimmedLower[i] = strings.ToLower(strings.TrimSpace(v))
		}
		result = trimmedLower

	default:
		result = input
	}
	return result.(T)
}

// TrimUpper converts the string, the value of a string pointer, and the
// string slice to the trimmed, upper result.
//
// For example, a given string of `" abc "` will be converted to `"ABC"`;
// a given string pointer containing the value `" abc "` will be converted
// to `"ABC"` and results in a new string pointer; a given string slice
// with `["abc", "Bar"]` will be converted to `["ABC", "BAR"]`.
func TrimUpper[T string | *string | []string](input T) T {
	var result any

	switch s := any(input).(type) {
	case string:
		result = strings.ToUpper(strings.TrimSpace(s))

	case *string:
		if s == nil {
			return input
		}
		trimmedUpper := strings.ToUpper(strings.TrimSpace(*s))
		result = &trimmedUpper

	case []string:
		trimmedUpper := make([]string, len(s))
		for i, v := range s {
			trimmedUpper[i] = strings.ToUpper(strings.TrimSpace(v))
		}
		result = trimmedUpper

	default:
		result = input
	}

	return result.(T)
}
