package crypt

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SM4_ECB_EncryptDecrypt(t *testing.T) {
	src := []byte("apple")
	key := []byte("1234567812345678")
	want := "uQkf7ZyRiQaUtzOgk2Vl7A==" // base64

	dst, err := SM4ECBEncrypt(src, key)
	assert.NoError(t, err)
	assert.EqualValues(t, want, base64.StdEncoding.EncodeToString(dst))

	dst, err = SM4ECBDecrypt(dst, key)
	assert.NoError(t, err)
	assert.EqualValues(t, src, dst)
}

func Test_SM4_CBC_EncryptDecrypt(t *testing.T) {
	src := []byte("apple")
	key := []byte("1234567812345678")
	iv := []byte("5678567856785678")
	want := "zLvNgbU8k82SVCr22wQ0lw==" // base64

	dst, err := SM4CBCEncrypt(src, key, iv)
	assert.NoError(t, err)
	assert.EqualValues(t, want, base64.StdEncoding.EncodeToString(dst))

	dst, err = SM4CBCDecrypt(dst, key, iv)
	assert.NoError(t, err)
	assert.EqualValues(t, src, dst)
}

func TestSM4CTRXor(t *testing.T) {
	src := []byte("apple")
	key := []byte("1234567812345678")
	iv := []byte("5678567856785678")
	want := "PAWoKg4=" // base64

	dst, err := SM4CTRXor(src, key, iv) // no padding to src
	assert.NoError(t, err)
	assert.EqualValues(t, want, base64.StdEncoding.EncodeToString(dst))

	dst, err = SM4CTRXor(dst, key, iv)
	assert.NoError(t, err)
	assert.EqualValues(t, src, dst)
}

func Test_SM4_CTR_EncryptDecrypt(t *testing.T) {
	src := []byte("apple")
	key := []byte("1234567812345678")
	iv := []byte("5678567856785678")
	want := "PAWoKg59C8+j6GO+mnnyPg==" // base64

	dst, err := SM4CTREncrypt(src, key, iv, PaddingPKCS7)
	assert.NoError(t, err)
	assert.EqualValues(t, want, base64.StdEncoding.EncodeToString(dst))

	dst, err = SM4CTRDecrypt(dst, key, iv, PaddingPKCS7)
	assert.NoError(t, err)
	assert.EqualValues(t, src, dst)
}
