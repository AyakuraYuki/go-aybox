package crypt

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
)

// region HmacSHA1

func HmacSHA1(src, key []byte) []byte {
	h := hmac.New(sha1.New, key)
	_, _ = h.Write(src)
	return h.Sum(nil)
}

// endregion

// region HmacSHA256

func HmacSHA256(src, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	_, _ = h.Write(src)
	return h.Sum(nil)
}

func HmacSHA256224(src, key []byte) []byte {
	h := hmac.New(sha256.New224, key)
	_, _ = h.Write(src)
	return h.Sum(nil)
}

// endregion

// region HmacSHA512

func HmacSHA512(src, key []byte) []byte {
	h := hmac.New(sha512.New, key)
	_, _ = h.Write(src)
	return h.Sum(nil)
}

func HmacSHA512224(src, key []byte) []byte {
	h := hmac.New(sha512.New512_224, key)
	_, _ = h.Write(src)
	return h.Sum(nil)
}

func HmacSHA512256(src, key []byte) []byte {
	h := hmac.New(sha512.New512_256, key)
	_, _ = h.Write(src)
	return h.Sum(nil)
}

func HmacSHA384(src, key []byte) []byte {
	h := hmac.New(sha512.New384, key)
	_, _ = h.Write(src)
	return h.Sum(nil)
}

// endregion
