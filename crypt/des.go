package crypt

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
)

// DesECBEncrypt DES-ECB encrypt
func DesECBEncrypt(src, key []byte, padding string) ([]byte, error) {
	block, err := NewDesCipher(key)
	if err != nil {
		return nil, err
	}
	return ECBEncrypt(block, src, padding)
}

// DesECBDecrypt DES-ECB decrypt
func DesECBDecrypt(src, key []byte, padding string) ([]byte, error) {
	block, err := NewDesCipher(key)
	if err != nil {
		return nil, err
	}
	return ECBDecrypt(block, src, padding)
}

// DesCBCEncrypt DES-CBC encrypt
func DesCBCEncrypt(src, key, iv []byte, padding string) ([]byte, error) {
	block, err := NewDesCipher(key)
	if err != nil {
		return nil, err
	}
	return CBCEncrypt(block, src, iv, padding)
}

// DesCBCDecrypt DES-CBC decrypt
func DesCBCDecrypt(src, key, iv []byte, padding string) ([]byte, error) {
	block, err := NewDesCipher(key)
	if err != nil {
		return nil, err
	}
	return CBCDecrypt(block, src, iv, padding)
}

// NewDesCipher creates a new DES cipher.Block
func NewDesCipher(key []byte) (cipher.Block, error) {
	if len(key) < des.BlockSize {
		key = append(key, bytes.Repeat([]byte{0}, des.BlockSize-len(key))...)
	} else {
		key = key[:des.BlockSize]
	}
	return des.NewCipher(key)
}
