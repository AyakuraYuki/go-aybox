package random

import "math/rand"

const (
	Lower           = 1 << 0
	Upper           = 1 << 1
	Digit           = 1 << 2
	LowerUpper      = Lower | Upper
	LowerDigit      = Lower | Digit
	UpperDigit      = Upper | Digit
	LowerUpperDigit = LowerUpper | Digit
)

const (
	lower = "abcdefghijklmnopqrstuvwxyz"
	upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digit = "0123456789"
)

func RandString(size, set int) string {
	letters := ""
	if set&Lower > 0 {
		letters += lower
	}
	if set&Upper > 0 {
		letters += upper
	}
	if set&Digit > 0 {
		letters += digit
	}
	letterSize := len(letters)
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = letters[rand.Intn(letterSize)]
	}
	return string(buf)
}
