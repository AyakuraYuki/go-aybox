package crypt

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRC4(t *testing.T) {
	src := []byte("apple")
	key := []byte("secret")
	dst, err := RC4(src, key)
	assert.NoError(t, err)
	want, _ := hex.DecodeString("8c46a270e7")
	assert.EqualValues(t, want, dst)
}

func TestRC4ToHex(t *testing.T) {
	src := []byte("apple")
	key := []byte("secret")
	dst, err := RC4ToHex(src, key)
	assert.NoError(t, err)
	assert.EqualValues(t, "8c46a270e7", dst)
}

func TestRC4FromHex(t *testing.T) {
	key := []byte("secret")
	dst, err := RC4FromHex("8c46a270e7", key)
	assert.NoError(t, err)
	assert.EqualValues(t, "apple", string(dst))
}

func TestRC4ToBase64(t *testing.T) {
	src := []byte("apple")
	key := []byte("secret")
	dst, err := RC4ToBase64(src, key)
	assert.NoError(t, err)
	assert.EqualValues(t, "jEaicOc=", dst)
}

func TestRC4FromBase64(t *testing.T) {
	key := []byte("secret")
	dst, err := RC4FromBase64("jEaicOc=", key)
	assert.NoError(t, err)
	assert.EqualValues(t, "apple", string(dst))
}

func TestRC4ToURLBase64(t *testing.T) {
	src := []byte("apple")
	key := []byte("secret")
	dst, err := RC4ToURLBase64(src, key)
	assert.NoError(t, err)
	assert.EqualValues(t, "jEaicOc=", dst)
}

func TestRC4FromURLBase64(t *testing.T) {
	key := []byte("secret")
	dst, err := RC4FromURLBase64("jEaicOc=", key)
	assert.NoError(t, err)
	assert.EqualValues(t, "apple", string(dst))
}
