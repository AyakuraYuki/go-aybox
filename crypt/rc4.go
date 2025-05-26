package crypt

import (
	"crypto/rc4"
	"encoding/base64"
	"encoding/hex"
)

func RC4(src, key []byte) ([]byte, error) {
	cipher, err := rc4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	dst := make([]byte, len(src))
	cipher.XORKeyStream(dst, src)
	return dst, nil
}

func RC4ToHex(src, key []byte) (string, error) {
	dst, err := RC4(src, key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(dst), nil
}

func RC4FromHex(src string, key []byte) ([]byte, error) {
	bs, err := hex.DecodeString(src)
	if err != nil {
		return nil, err
	}
	return RC4(bs, key)
}

func RC4ToBase64(src, key []byte) (string, error) {
	dst, err := RC4(src, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(dst), nil
}

func RC4FromBase64(src string, key []byte) ([]byte, error) {
	bs, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, err
	}
	return RC4(bs, key)
}

func RC4ToURLBase64(src, key []byte) (string, error) {
	dst, err := RC4(src, key)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(dst), nil
}

func RC4FromURLBase64(src string, key []byte) ([]byte, error) {
	bs, err := base64.URLEncoding.DecodeString(src)
	if err != nil {
		return nil, err
	}
	return RC4(bs, key)
}
