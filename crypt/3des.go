package crypt

import "crypto/des"

// TripleDesECBEncrypt 3DES-ECB encrypt
func TripleDesECBEncrypt(src, key []byte, padding string) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	return ECBEncrypt(block, src, padding)
}

// TripleDesECBDecrypt 3DES-ECB decrypt
func TripleDesECBDecrypt(src, key []byte, padding string) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	return ECBDecrypt(block, src, padding)
}

// TripleDesCBCEncrypt 3DES-CBC encrypt
func TripleDesCBCEncrypt(src, key, iv []byte, padding string) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	return CBCEncrypt(block, src, iv, padding)
}

// TripleDesCBCDecrypt 3DES-CBC decrypt
func TripleDesCBCDecrypt(src, key, iv []byte, padding string) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	return CBCDecrypt(block, src, iv, padding)
}
