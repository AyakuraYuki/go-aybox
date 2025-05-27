package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

// AesECBEncrypt AES-ECB encrypt
func AesECBEncrypt(src, key []byte, padding string) ([]byte, error) {
	block, err := NewAesCipher(key)
	if err != nil {
		return nil, err
	}
	return ECBEncrypt(block, src, padding)
}

// AesECBDecrypt AES-ECB decrypt
func AesECBDecrypt(src, key []byte, padding string) ([]byte, error) {
	block, err := NewAesCipher(key)
	if err != nil {
		return nil, err
	}
	return ECBDecrypt(block, src, padding)
}

// AesCBCEncrypt AES-CBC encrypt
func AesCBCEncrypt(src, key, iv []byte, padding string) ([]byte, error) {
	block, err := NewAesCipher(key)
	if err != nil {
		return nil, err
	}
	return CBCEncrypt(block, src, iv, padding)
}

// AesCBCDecrypt AES-CBC decrypt
func AesCBCDecrypt(src, key, iv []byte, padding string) ([]byte, error) {
	block, err := NewAesCipher(key)
	if err != nil {
		return nil, err
	}
	return CBCDecrypt(block, src, iv, padding)
}

// NewAesCipher creates a new AES cipher.Block
func NewAesCipher(key []byte) (cipher.Block, error) {
	return aes.NewCipher(aesKeyPending(key))
}

func aesKeyPending(key []byte) []byte {
	k := len(key)
	count := 0
	switch true {
	case k <= 16:
		count = 16 - k
	case k <= 24:
		count = 24 - k
	case k <= 32:
		count = 32 - k
	default:
		return key[:32]
	}
	if count == 0 {
		return key
	}

	return append(key, bytes.Repeat([]byte{0}, count)...)
}
