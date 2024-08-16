package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

// aes加密算法，CBC , pkcs7
func Encrypt(plantText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key) //选择加密算法
	if err != nil {
		return nil, err
	}
	plantText = PKCS7Padding(plantText, block.BlockSize())

	blockModel := cipher.NewCBCEncrypter(block, key)

	ciphertext := make([]byte, len(plantText))

	blockModel.CryptBlocks(ciphertext, plantText)
	return ciphertext, nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func Decrypt(ciphertext, key []byte) ([]byte, error) {
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes) //选择加密算法
	if err != nil {
		return nil, err
	}
	blockModel := cipher.NewCBCDecrypter(block, keyBytes)
	plantText := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(plantText, ciphertext)
	plantText = PKCS7UnPadding(plantText, block.BlockSize())
	return plantText, nil
}

func PKCS7UnPadding(plantText []byte, blockSize int) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

//解密字符串
func CbcDecrypt(src, key, iv []byte) (strDesc string, err error) {
	aesBlockDecrypter, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockModel := cipher.NewCBCDecrypter(aesBlockDecrypter, iv)
	plantText := make([]byte, len(src))
	blockModel.CryptBlocks(plantText, src)
	plantText = PKCS7UnPadding(plantText, aesBlockDecrypter.BlockSize())
	return string(plantText), nil
}
