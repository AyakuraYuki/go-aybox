package crypt

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testExpected3DesECBBase64 = "goTWM84H8NxbrJ4sDaHiRw=="
	testExpected3DesCBCBase64 = "g6Im7UM3ObVsJ3HMJW/6Fw=="
)

func TestTripleDesECBEncrypt(t *testing.T) {
	dst, err := TripleDesECBEncrypt(testSrc(), test3DesKey(), PaddingPKCS7)
	assert.NoError(t, err)
	assert.EqualValues(t, testExpected3DesECBBase64, base64.StdEncoding.EncodeToString(dst))
}

func TestTripleDesECBDecrypt(t *testing.T) {
	src, _ := base64.StdEncoding.DecodeString(testExpected3DesECBBase64)
	dst, err := TripleDesECBDecrypt(src, test3DesKey(), PaddingPKCS7)
	assert.NoError(t, err)
	assert.EqualValues(t, testSrc(), dst)
}

func TestTripleDesCBCEncrypt(t *testing.T) {
	dst, err := TripleDesCBCEncrypt(testSrc(), test3DesKey(), testDesIV(), PaddingPKCS7)
	assert.NoError(t, err)
	assert.EqualValues(t, testExpected3DesCBCBase64, base64.StdEncoding.EncodeToString(dst))
}

func TestTripleDesCBCDecrypt(t *testing.T) {
	src, _ := base64.StdEncoding.DecodeString(testExpected3DesCBCBase64)
	dst, err := TripleDesCBCDecrypt(src, test3DesKey(), testDesIV(), PaddingPKCS7)
	assert.NoError(t, err)
	assert.EqualValues(t, testSrc(), dst)
}
