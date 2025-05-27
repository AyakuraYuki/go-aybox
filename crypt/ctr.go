package crypt

import (
	"bytes"
	"crypto/cipher"
)

func CTREncrypt(block cipher.Block, src, iv []byte, padding string) ([]byte, error) {
	blockSize := block.BlockSize()
	src = Padding(padding, src, blockSize)
	dst := make([]byte, len(src))
	if len(iv) != blockSize {
		iv = ctrIVPadding(iv, blockSize)
	}
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(dst, src)
	return dst, nil
}

func CTRDecrypt(block cipher.Block, src, iv []byte, padding string) ([]byte, error) {
	blockSize := block.BlockSize()
	dst := make([]byte, len(src))
	if len(iv) != blockSize {
		iv = ctrIVPadding(iv, blockSize)
	}
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(dst, src)
	return UnPadding(padding, dst)
}

func ctrIVPadding(iv []byte, blockSize int) []byte {
	if k := len(iv); k < blockSize {
		return append(iv, bytes.Repeat([]byte{0}, blockSize-k)...)
	} else if k > blockSize {
		return iv[:blockSize]
	}
	return iv
}
