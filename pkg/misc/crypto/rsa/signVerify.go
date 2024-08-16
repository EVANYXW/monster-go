package rsa

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

// 签名
func Sign(plant []byte, private []byte) (string, error) {
	block, _ := pem.Decode([]byte(private))
	if block == nil {
		return "", nil
	}
	prk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	var hash crypto.Hash
	hash = crypto.SHA1
	h := hash.New()
	h.Write(plant)
	hashed := h.Sum(nil)
	signBytes, err := rsa.SignPKCS1v15(rand.Reader, prk, hash, hashed)
	if err != nil {
		return "", err
	}
	sign := base64.StdEncoding.EncodeToString(signBytes)
	return sign, err
}

// 签名
func Sign256(plant []byte, private []byte) (string, error) {
	block, _ := pem.Decode([]byte(private))
	if block == nil {
		return "", nil
	}
	prk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	var hash crypto.Hash
	hash = crypto.SHA256
	h := hash.New()
	h.Write(plant)
	hashed := h.Sum(nil)
	signBytes, err := rsa.SignPKCS1v15(rand.Reader, prk, hash, hashed)
	if err != nil {
		return "", err
	}
	sign := base64.StdEncoding.EncodeToString(signBytes)
	return sign, err
}

// 验证签名
func Verify(sign string, plant []byte, public []byte) (bool, error) {
	block, _ := pem.Decode([]byte(public))
	if block == nil {
		return false, nil
	}
	puk, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}
	var hash crypto.Hash
	hash = crypto.SHA1
	h := hash.New()
	h.Write(plant)
	hashed := h.Sum(nil)
	bsign, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false, err
	}
	err = rsa.VerifyPKCS1v15(puk.(*rsa.PublicKey), hash, hashed, bsign)
	if err != nil {
		return false, err
	}
	return true, nil
}

// 解密
func RsaDecode(data string, private []byte) ([]byte, error) {
	block, _ := pem.Decode(private)
	if block == nil {
		return nil, errors.New("证书异常")
	}
	prk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	str1, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	groups := grouping(str1, len(prk.PublicKey.N.Bytes()))
	buffer := bytes.Buffer{}
	for _, cipherTextBlock := range groups {
		plainText, err := rsa.DecryptPKCS1v15(rand.Reader, prk, cipherTextBlock)
		if err != nil {
			return nil, err
		}
		buffer.Write(plainText)
	}
	return buffer.Bytes(), nil
}

func RsaEncrypt(orgidata, publickey []byte) (string, error) {
	block, _ := pem.Decode(publickey)
	if block == nil {
		return "", errors.New("public key is bad")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	pub := pubInterface.(*rsa.PublicKey)

	d, err := rsa.EncryptPKCS1v15(rand.Reader, pub, orgidata) //加密
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(d), nil
}

func RsaEncryptBlock(src, publicKeyByte []byte) (string, error) {
	block, _ := pem.Decode(publicKeyByte)
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	pub := pubInterface.(*rsa.PublicKey)
	keySize, srcSize := pub.Size(), len(src)
	offSet, once := 0, keySize-11
	buffer := bytes.Buffer{}
	for offSet < srcSize {
		endIndex := offSet + once
		if endIndex > srcSize {
			endIndex = srcSize
		}
		// 加密一部分
		bytesOnce, err := rsa.EncryptPKCS1v15(rand.Reader, pub, src[offSet:endIndex])
		if err != nil {
			return "", err
		}
		buffer.Write(bytesOnce)
		offSet = endIndex
	}
	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

func RsaDecryptBlock(data string, privateKeyBytes []byte) (bytesDecrypt []byte, err error) {
	src, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privateKeyBytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}
	keySize := privateKey.Size()
	srcSize := len(src)
	var offSet = 0
	var buffer = bytes.Buffer{}
	for offSet < srcSize {
		endIndex := offSet + keySize
		if endIndex > srcSize {
			endIndex = srcSize
		}
		bytesOnce, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, src[offSet:endIndex])
		if err != nil {
			return nil, err
		}
		buffer.Write(bytesOnce)
		offSet = endIndex
	}
	bytesDecrypt = buffer.Bytes()
	return
}
