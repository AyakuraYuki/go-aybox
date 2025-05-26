package crypt

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testExpectedAes128ECBBase64 = "3Wc9SzUtEYmDZDGrgZrmiQ=="
	testExpectedAes192ECBBase64 = "rkxYCqB6vzL346abZMEZMA=="
	testExpectedAes256ECBBase64 = "A0iMsfS+C2UxzzyycSa3JQ=="

	testExpectedAes128CBCBase64 = "lZUlbp47X/Lbr7eCUZBm+w=="
	testExpectedAes192CBCBase64 = "kveXsZ33ZLtSu7TuFaa6xg=="
	testExpectedAes256CBCBase64 = "Ac76EhWrN0mWrjzqla4jyg=="
)

func TestAesECBEncrypt(t *testing.T) {
	// AES128/ECB/PKCS7
	dst, err := AesECBEncrypt(testSrc(), testAes128Key(), PaddingPKCS7)
	assert.NoError(t, err)
	assert.EqualValues(t, testExpectedAes128ECBBase64, base64.StdEncoding.EncodeToString(dst))

	// AES192/ECB/PKCS7
	dst, err = AesECBEncrypt(testSrc(), testAes192Key(), PaddingPKCS7)
	assert.NoError(t, err)
	assert.EqualValues(t, testExpectedAes192ECBBase64, base64.StdEncoding.EncodeToString(dst))

	// AES256/ECB/PKCS7
	dst, err = AesECBEncrypt(testSrc(), testAes256Key(), PaddingPKCS7)
	assert.NoError(t, err)
	assert.EqualValues(t, testExpectedAes256ECBBase64, base64.StdEncoding.EncodeToString(dst))
}

func TestAesECBDecrypt(t *testing.T) {
	// AES128/ECB/PKCS7
	src, _ := base64.StdEncoding.DecodeString(testExpectedAes128ECBBase64)
	dst, err := AesECBDecrypt(src, testAes128Key(), PaddingPKCS7)
	assert.NoError(t, err)
	assert.EqualValues(t, testSrc(), dst)

	// AES192/ECB/PKCS7
	src, _ = base64.StdEncoding.DecodeString(testExpectedAes192ECBBase64)
	dst, err = AesECBDecrypt(src, testAes192Key(), PaddingPKCS7)
	assert.NoError(t, err)
	assert.EqualValues(t, testSrc(), dst)

	// AES256/ECB/PKCS7
	src, _ = base64.StdEncoding.DecodeString(testExpectedAes256ECBBase64)
	dst, err = AesECBDecrypt(src, testAes256Key(), PaddingPKCS7)
	assert.NoError(t, err)
	assert.EqualValues(t, testSrc(), dst)
}

func TestAesCBCEncrypt(t *testing.T) {
	// AES128/CBC/PKCS7
	dst, err := AesCBCEncrypt(testSrc(), testAes128Key(), testAesIV(), PaddingPKCS7)
	assert.NoError(t, err)
	assert.EqualValues(t, testExpectedAes128CBCBase64, base64.StdEncoding.EncodeToString(dst))

	// AES192/CBC/PKCS7
	dst, err = AesCBCEncrypt(testSrc(), testAes192Key(), testAesIV(), PaddingPKCS7)
	assert.NoError(t, err)
	assert.EqualValues(t, testExpectedAes192CBCBase64, base64.StdEncoding.EncodeToString(dst))

	// AES256/CBC/PKCS7
	dst, err = AesCBCEncrypt(testSrc(), testAes256Key(), testAesIV(), PaddingPKCS7)
	assert.NoError(t, err)
	assert.EqualValues(t, testExpectedAes256CBCBase64, base64.StdEncoding.EncodeToString(dst))
}

func TestAesCBCDecrypt(t *testing.T) {
	// AES128/ECB/PKCS7
	src, _ := base64.StdEncoding.DecodeString(testExpectedAes128CBCBase64)
	dst, err := AesCBCDecrypt(src, testAes128Key(), testAesIV(), PaddingPKCS7)
	assert.NoError(t, err)
	assert.EqualValues(t, testSrc(), dst)

	// AES192/ECB/PKCS7
	src, _ = base64.StdEncoding.DecodeString(testExpectedAes192CBCBase64)
	dst, err = AesCBCDecrypt(src, testAes192Key(), testAesIV(), PaddingPKCS7)
	assert.NoError(t, err)
	assert.EqualValues(t, testSrc(), dst)

	// AES256/ECB/PKCS7
	src, _ = base64.StdEncoding.DecodeString(testExpectedAes256CBCBase64)
	dst, err = AesCBCDecrypt(src, testAes256Key(), testAesIV(), PaddingPKCS7)
	assert.NoError(t, err)
	assert.EqualValues(t, testSrc(), dst)
}
