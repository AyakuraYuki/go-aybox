package crypt

import (
	"bytes"
	"crypto/cipher"
)

func GCMEncrypt(block cipher.Block, src, nonce, A []byte, padding string) (dst []byte, err error) {
	blockSize := block.BlockSize()
	src = Padding(padding, src, blockSize)
	mode, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := mode.NonceSize()
	if len(nonce) != nonceSize {
		nonce = gcmNoncePadding(nonce, nonceSize)
	}
	dst = mode.Seal(nil, nonce, src, A)
	return dst, nil
}

func GCMDecrypt(block cipher.Block, src, nonce, A []byte, padding string) (dst []byte, err error) {
	mode, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := mode.NonceSize()
	if len(nonce) != nonceSize {
		nonce = gcmNoncePadding(nonce, nonceSize)
	}
	dst, err = mode.Open(nil, nonce, src, A)
	if err != nil {
		return nil, err
	}
	return UnPadding(padding, dst)
}

func gcmNoncePadding(nonce []byte, nonceSize int) []byte {
	if k := len(nonce); k < nonceSize {
		return append(nonce, bytes.Repeat([]byte{0}, nonceSize-k)...)
	} else if k > nonceSize {
		return nonce[:nonceSize]
	}
	return nonce
}
