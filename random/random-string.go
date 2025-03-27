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
	lower  = "abcdefghijklmnopqrstuvwxyz"
	upper  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digit  = "0123456789"
	symbol = "!#$%&()*+,-./:;<=>?@[]^_{|}~\\\"'"
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

func RandPassword(size int) string {
	letters := ""
	letters += lower
	letters += upper
	letters += digit
	letters += symbol
	letterSize := len(letters)
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		switch i {
		case 0:
			buf[i] = upper[rand.Intn(len(upper))] // first should be upper
		case 1:
			buf[i] = lower[rand.Intn(len(lower))] // second should be lower
		case 2:
			buf[i] = symbol[rand.Intn(len(symbol))] // third should be a symbol
		case 3:
			buf[i] = digit[rand.Intn(len(digit))] // forth should be a digit
		default:
			buf[i] = letters[rand.Intn(letterSize)]
		}
	}
	return string(buf)
}
