package str

import "strings"

// TrimLower converts the string, the value of a string pointer, and the
// string slice to the trimmed, lower result.
//
// For example, a given string of `" ABC "` will be converted to `"abc"`;
// a given string pointer containing the value `" ABC "` will be converted
// to `"abc"` and results in a new string pointer; a given string slice
// with `["ABC", "Bar"]` will be converted to `["abc", "bar"]`.
func TrimLower[T ~string | ~*string | ~[]string](input T) T {
	return trimCase[T](input, strings.ToLower)
}

// TrimUpper converts the string, the value of a string pointer, and the
// string slice to the trimmed, upper result.
//
// For example, a given string of `" abc "` will be converted to `"ABC"`;
// a given string pointer containing the value `" abc "` will be converted
// to `"ABC"` and results in a new string pointer; a given string slice
// with `["abc", "Bar"]` will be converted to `["ABC", "BAR"]`.
func TrimUpper[T ~string | ~*string | ~[]string](input T) T {
	return trimCase[T](input, strings.ToUpper)
}

func trimCase[T ~string | ~*string | ~[]string](input T, toCase func(string) string) T {
	var result any

	switch s := any(input).(type) {
	case string:
		result = toCase(strings.TrimSpace(s))

	case *string:
		if s == nil {
			return input
		}
		result = new(toCase(strings.TrimSpace(*s)))

	case []string:
		trimmed := make([]string, len(s))
		for i, v := range s {
			trimmed[i] = toCase(strings.TrimSpace(v))
		}
		result = trimmed

	default:
		result = input
	}

	return result.(T)
}
