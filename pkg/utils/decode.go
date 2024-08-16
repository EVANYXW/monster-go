package utils

import (
	"encoding/base64"
	"github.com/evanyxw/monster-go/pkg/misc/crypto/des"
)

// 加密数据 DES(ECB,PKCS5) + base64
func EncodeData(data, key string) (string, error) {
	if len(data) == 0 {
		return "", nil
	}

	ret, err := des.Encrypt([]byte(data), []byte(key))
	if err != nil {
		return "", err
	}

	debyte := base64.StdEncoding.EncodeToString(ret)
	return string(debyte), nil
}

// 解密数据 DES(ECB,PKCS5) + base64
func DecodeData(data, key string) (string, error) {
	if len(data) == 0 {
		return "", nil
	}

	bdeb, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	det, err := des.Decrypt(bdeb, []byte(key))
	if err != nil {
		return "", err
	}

	return string(det), nil
}
