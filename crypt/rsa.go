package crypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// RSAGenerateKey generates RSA private key, returns bytes
func RSAGenerateKey(bits int) (priKey []byte, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}

	x509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)

	block := pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509PrivateKey}

	return pem.EncodeToMemory(&block), nil
}

// RSAGeneratePublicKey generates RSA public keys, returns bytes
func RSAGeneratePublicKey(priKey []byte) (pubKey []byte, err error) {
	block, _ := pem.Decode(priKey)
	if block == nil {
		return nil, errors.New("key is invalid format")
	}

	// x509 parse
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.PublicKey
	x509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return nil, err
	}

	publicBlock := pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509PublicKey}

	return pem.EncodeToMemory(&publicBlock), nil
}

// RSAEncrypt RSA encrypt by public key
func RSAEncrypt(src, pubKey []byte) ([]byte, error) {
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return nil, errors.New("key is invalid format")
	}

	// x509 parse
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("the kind of key is not a rsa.PublicKey")
	}

	// encrypt
	dst, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, src)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

// RSADecrypt RSA decrypt by private key
func RSADecrypt(src, priKey []byte) ([]byte, error) {
	block, _ := pem.Decode(priKey)
	if block == nil {
		return nil, errors.New("key is invalid format")
	}

	// x509 parse
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	dst, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, src)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

// RSASign RSA sign by private key
func RSASign(src, priKey []byte, hash crypto.Hash) ([]byte, error) {
	block, _ := pem.Decode(priKey)
	if block == nil {
		return nil, errors.New("key is invalid format")
	}

	// x509 parse
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	h := hash.New()
	if _, err = h.Write(src); err != nil {
		return nil, err
	}

	bytes := h.Sum(nil)
	sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, hash, bytes)
	if err != nil {
		return nil, err
	}

	return sign, nil
}

// RSAVerify RSA verify sign by public key
func RSAVerify(src, sign, pubKey []byte, hash crypto.Hash) error {
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return errors.New("key is invalid format")
	}

	// x509 parse
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return errors.New("the kind of key is not a rsa.PublicKey")
	}

	h := hash.New()
	if _, err = h.Write(src); err != nil {
		return err
	}

	bytes := h.Sum(nil)

	return rsa.VerifyPKCS1v15(publicKey, hash, bytes, sign)
}
