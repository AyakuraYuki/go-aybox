package crypt

import (
	"github.com/emmansun/gmsm/sm4"
)

// region SM4 ECB

// SM4ECBEncrypt SM4 ECB encrypt
func SM4ECBEncrypt(src, key []byte) (dst []byte, err error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return ECBEncrypt(block, src, PaddingPKCS7)
}

// SM4ECBDecrypt SM4 ECB decrypt
func SM4ECBDecrypt(src, key []byte) (dst []byte, err error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return ECBDecrypt(block, src, PaddingPKCS7)
}

// endregion

// region SM4 CBC

// SM4CBCEncrypt SM4 CBC encrypt
func SM4CBCEncrypt(src, key, iv []byte) (dst []byte, err error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return CBCEncrypt(block, src, iv, PaddingPKCS7)
}

// SM4CBCDecrypt SM4 CBC decrypt
func SM4CBCDecrypt(src, key, iv []byte) (dst []byte, err error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return CBCDecrypt(block, src, iv, PaddingPKCS7)
}

// endregion

// region SM4 CTR

func SM4CTRXor(src, key, iv []byte) (dst []byte, err error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return CTRXor(block, src, iv)
}

func SM4CTREncrypt(src, key, iv []byte, padding string) (dst []byte, err error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return CTREncrypt(block, src, iv, padding)
}

func SM4CTRDecrypt(src, key, iv []byte, padding string) (dst []byte, err error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return CTRDecrypt(block, src, iv, padding)
}

// endregion
