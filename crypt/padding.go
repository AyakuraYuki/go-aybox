package crypt

import (
	"bytes"
	"errors"
)

var ErrUnPadding = errors.New("UnPadding error")

const PaddingPKCS5 = "PKCS5"
const PaddingPKCS7 = "PKCS7"
const PaddingZEROS = "ZEROS"

func Padding(padding string, src []byte, blockSize int) []byte {
	switch padding {
	case PaddingPKCS5:
		src = PKCS5Padding(src, blockSize)
	case PaddingPKCS7:
		src = PKCS7Padding(src, blockSize)
	case PaddingZEROS:
		src = ZerosPadding(src, blockSize)
	}
	return src
}

func UnPadding(padding string, src []byte) ([]byte, error) {
	switch padding {
	case PaddingPKCS5:
		return PKCS5UnPadding(src)
	case PaddingPKCS7:
		return PKCS7UnPadding(src)
	case PaddingZEROS:
		return ZerosUnPadding(src)
	}
	return src, nil
}

func PKCS5Padding(src []byte, blockSize int) []byte {
	return PKCS7Padding(src, blockSize)
}

func PKCS5UnPadding(src []byte) ([]byte, error) {
	return PKCS7UnPadding(src)
}

func PKCS7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func PKCS7UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return src, ErrUnPadding
	}
	unpadding := int(src[length-1])
	if length < unpadding {
		return src, ErrUnPadding
	}
	return src[:(length - unpadding)], nil
}

func ZerosPadding(src []byte, blockSize int) []byte {
	paddingCount := blockSize - len(src)%blockSize
	if paddingCount == 0 {
		return src
	} else {
		return append(src, bytes.Repeat([]byte{byte(0)}, paddingCount)...)
	}
}

func ZerosUnPadding(src []byte) ([]byte, error) {
	for i := len(src) - 1; ; i-- {
		if src[i] != 0 {
			return src[:i+1], nil
		}
	}
}
