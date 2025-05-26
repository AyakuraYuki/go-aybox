package crypt

import "crypto/cipher"

func ECBEncrypt(block cipher.Block, src []byte, padding string) ([]byte, error) {
	blockSize := block.BlockSize()
	src = Padding(padding, src, blockSize)
	dst := make([]byte, len(src))
	mode := NewECBEncrypter(block)
	mode.CryptBlocks(dst, src)
	return dst, nil
}

func ECBDecrypt(block cipher.Block, src []byte, padding string) ([]byte, error) {
	dst := make([]byte, len(src))
	mode := NewECBDecrypter(block)
	mode.CryptBlocks(dst, src)
	return UnPadding(padding, dst)
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{b: b, blockSize: b.BlockSize()}
}

type ecbEncrypter ecb

// NewECBEncrypter returns a [cipher.BlockMode] which encrypts in Electronic
// Codebook mode, using the given [cipher.Block].
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}

func (x *ecbEncrypter) BlockSize() int { return x.blockSize }

func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

// NewECBDecrypter returns a [cipher.BlockMode] which decrypts in Electronic
// Codebook mode, using the given [cipher.Block].
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}

func (x *ecbDecrypter) BlockSize() int { return x.blockSize }

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}
