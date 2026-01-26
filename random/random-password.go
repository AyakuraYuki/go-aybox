package random

import (
	"math/rand"
)

const (
	passwordLower    = "abcdefghijkmnopqrstuvwxyz"
	passwordLowerLen = len(passwordLower)

	passwordUpper    = "ABCDEFGHJKLMNPQRTUVWXYZ"
	passwordUpperLen = len(passwordUpper)

	passwordDigit    = "2346789"
	passwordDigitLen = len(passwordDigit)

	passwordSymbol    = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	passwordSymbolLen = len(passwordSymbol)

	passwordAll    = passwordLower + passwordUpper + passwordDigit + passwordSymbol
	passwordAllLen = len(passwordAll)
)

func RandPassword(length int) string {
	if length < 4 {
		length = 4
	}

	buf := make([]byte, length)

	buf[0] = passwordLower[rand.Intn(passwordLowerLen)]
	buf[1] = passwordUpper[rand.Intn(passwordUpperLen)]
	buf[2] = passwordSymbol[rand.Intn(passwordSymbolLen)]
	buf[3] = passwordDigit[rand.Intn(passwordDigitLen)]

	for i := 4; i < length; i++ {
		buf[i] = passwordAll[rand.Intn(passwordAllLen)]
	}

	return string(buf)
}
