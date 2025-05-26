package crypt

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha3"
	"crypto/sha512"
)

// region SHA1

func SHA1(src []byte) []byte {
	h := sha1.New()
	_, _ = h.Write(src)
	return h.Sum(nil)
}

// endregion

// region SHA256

func SHA256(src []byte) []byte {
	h := sha256.New()
	_, _ = h.Write(src)
	return h.Sum(nil)
}

func SHA224(src []byte) []byte {
	h := sha256.New224()
	_, _ = h.Write(src)
	return h.Sum(nil)
}

// endregion

// region SHA512

func SHA512(src []byte) []byte {
	h := sha512.New()
	_, _ = h.Write(src)
	return h.Sum(nil)
}

func SHA512224(src []byte) []byte {
	h := sha512.New512_224()
	_, _ = h.Write(src)
	return h.Sum(nil)
}

func SHA512256(src []byte) []byte {
	h := sha512.New512_256()
	_, _ = h.Write(src)
	return h.Sum(nil)
}

func SHA384(src []byte) []byte {
	h := sha512.New384()
	_, _ = h.Write(src)
	return h.Sum(nil)
}

// endregion

// region SHA3

func SHA3(src []byte) []byte {
	h := sha3.New512()
	_, _ = h.Write(src)
	return h.Sum(nil)
}

func SHA3224(src []byte) []byte {
	h := sha3.New224()
	_, _ = h.Write(src)
	return h.Sum(nil)
}

func SHA3256(src []byte) []byte {
	h := sha3.New256()
	_, _ = h.Write(src)
	return h.Sum(nil)
}

func SHA3384(src []byte) []byte {
	h := sha3.New384()
	_, _ = h.Write(src)
	return h.Sum(nil)
}

// endregion
