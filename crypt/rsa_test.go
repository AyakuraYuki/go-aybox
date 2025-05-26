package crypt

import (
	"crypto"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRSAEncryptDecrypt(t *testing.T) {
	priKey, err := RSAGenerateKey(2048)
	assert.NoError(t, err)
	t.Logf("private key: %s", priKey)

	pubKey, err := RSAGeneratePublicKey(priKey)
	assert.NoError(t, err)
	t.Logf("public key: %s", pubKey)

	src := []byte("apple")
	dst, err := RSAEncrypt(src, pubKey)
	assert.NoError(t, err)

	dst, err = RSADecrypt(dst, priKey)
	assert.NoError(t, err)

	assert.EqualValues(t, src, dst)
}

func TestRSASignVerify(t *testing.T) {
	priKey, err := RSAGenerateKey(2048)
	assert.NoError(t, err)
	t.Logf("private key: %s", priKey)

	pubKey, err := RSAGeneratePublicKey(priKey)
	assert.NoError(t, err)
	t.Logf("public key: %s", pubKey)

	src := []byte("apple")
	dst, err := RSASign(src, priKey, crypto.SHA512)
	assert.NoError(t, err)
	t.Logf("signature: %s", hex.EncodeToString(dst))

	err = RSAVerify(src, dst, pubKey, crypto.SHA512)
	assert.NoError(t, err)
}
