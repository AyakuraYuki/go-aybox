package crypt

import (
	"bytes"
	"crypto/cipher"
)

func CBCEncrypt(block cipher.Block, src, iv []byte, padding string) ([]byte, error) {
	blockSize := block.BlockSize()
	src = Padding(padding, src, blockSize)
	dst := make([]byte, len(src))
	if len(iv) != blockSize {
		iv = cbcIVPadding(iv, blockSize)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(dst, src)
	return dst, nil
}

func CBCDecrypt(block cipher.Block, src, iv []byte, padding string) ([]byte, error) {
	blockSize := block.BlockSize()
	dst := make([]byte, len(src))
	if len(iv) != blockSize {
		iv = cbcIVPadding(iv, blockSize)
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(dst, src)
	return UnPadding(padding, dst)
}

// cbcIVPadding auto pad length to block size
func cbcIVPadding(iv []byte, blockSize int) []byte {
	if k := len(iv); k < blockSize {
		return append(iv, bytes.Repeat([]byte{0}, blockSize-k)...)
	} else if k > blockSize {
		return iv[:blockSize]
	}
	return iv
}
