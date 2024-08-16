package rsa

import (
	"bytes"
	"crypto"
	"errors"
)



type Cipher interface {
	Encrypt(plainText []byte) ([]byte, error)
	Decrypt(cipherText []byte) ([]byte, error)
	Sign(src []byte, hash crypto.Hash) ([]byte, error)
	Verify(src []byte, sign []byte, hash crypto.Hash) error
}

func NewCipher(key Key, padding Padding, cipherMode CipherMode, signMode SignMode) Cipher {
	return &cipher{key: key, padding: padding, cipherMode: cipherMode, sign: signMode}
}

type cipher struct {
	key        Key
	cipherMode CipherMode
	sign       SignMode
	padding    Padding
}

func (cipher *cipher) Encrypt(plainText []byte) ([]byte, error) {
	groups := cipher.padding.Padding(plainText)
	buffer := bytes.Buffer{}
	for _, plainTextBlock := range groups {
		cipherText, err := cipher.cipherMode.Encrypt(plainTextBlock, cipher.key.PublicKey())
		if err != nil {
			return nil, err
		}
		buffer.Write(cipherText)
	}
	return buffer.Bytes(), nil
}

func (cipher *cipher) Decrypt(cipherText []byte) ([]byte, error) {
	if len(cipherText) == 0 {
		return nil, errors.New("密文不能为空")
	}
	groups := grouping(cipherText, cipher.key.Modulus())
	buffer := bytes.Buffer{}
	for _, cipherTextBlock := range groups {
		plainText, err := cipher.cipherMode.Decrypt(cipherTextBlock, cipher.key.PrivateKey())
		if err != nil {
			return nil, err
		}
		buffer.Write(plainText)
	}
	return buffer.Bytes(), nil
}

func (cipher *cipher) Sign(src []byte, hash crypto.Hash) ([]byte, error) {
	return cipher.sign.Sign(src, hash, cipher.key.PrivateKey())
}

func (cipher *cipher) Verify(src []byte, sign []byte, hash crypto.Hash) error {
	return cipher.sign.Verify(src, sign, hash, cipher.key.PublicKey())
}


func NewRSA(key Key) (Cipher, error) {
	padding := NewPKCS1Padding(key.Modulus())
	cipherMode := NewPKCS1v15Cipher()
	signMode := NewPKCS1v15Sign()
	return NewCipher(key, padding, cipherMode, signMode), nil
}

func NewRSAWith(key Key, padding Padding, cipherMode CipherMode, signMode SignMode) (Cipher, error) {
	return NewCipher(key, padding, cipherMode, signMode), nil
}