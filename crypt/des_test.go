package crypt

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testExpectedDesECBBase64 = "fF/QvLPpkgwE94ld2MWIfA=="
	testExpectedDesCBCBase64 = "u103iBiyRF3etJonMNMGNQ=="
)

func TestDesECBEncrypt(t *testing.T) {
	dst, err := DesECBEncrypt(testSrc(), testDesKey(), PaddingPKCS7)
	assert.NoError(t, err)
	assert.EqualValues(t, testExpectedDesECBBase64, base64.StdEncoding.EncodeToString(dst))
}

func TestDesECBDecrypt(t *testing.T) {
	src, _ := base64.StdEncoding.DecodeString(testExpectedDesECBBase64)
	dst, err := DesECBDecrypt(src, testDesKey(), PaddingPKCS7)
	assert.NoError(t, err)
	assert.EqualValues(t, testSrc(), dst)
}

func TestDesCBCEncrypt(t *testing.T) {
	dst, err := DesCBCEncrypt(testSrc(), testDesKey(), testDesIV(), PaddingPKCS7)
	assert.NoError(t, err)
	assert.EqualValues(t, testExpectedDesCBCBase64, base64.StdEncoding.EncodeToString(dst))
}

func TestDesCBCDecrypt(t *testing.T) {
	src, _ := base64.StdEncoding.DecodeString(testExpectedDesCBCBase64)
	dst, err := DesCBCDecrypt(src, testDesKey(), testDesIV(), PaddingPKCS7)
	assert.NoError(t, err)
	assert.EqualValues(t, testSrc(), dst)
}
