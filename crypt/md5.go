package crypt

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 calculates the md5 hash of a string
func MD5(str string) []byte {
	h := md5.New()
	_, _ = h.Write([]byte(str))
	return h.Sum(nil)
}

// MD5ToString calculates the md5 hash of a string, returns hex string
func MD5ToString(str string) string {
	return hex.EncodeToString(MD5(str))
}

// BinaryMD5 calculates the md5 hash of a binary
func BinaryMD5(src []byte) []byte {
	h := md5.New()
	_, _ = h.Write(src)
	return h.Sum(nil)
}

// BinaryMD5ToString calculates the md5 hash of a binary, returns hex string
func BinaryMD5ToString(src []byte) string {
	return hex.EncodeToString(BinaryMD5(src))
}
