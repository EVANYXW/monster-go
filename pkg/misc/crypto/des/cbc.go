/** 
 * cbc模式.
 *
 * User: zhangbob 
 * Date: 2018/7/5 
 * Time: 下午2:14 
 */
package des

import (
	"crypto/des"
	"crypto/cipher"
	"encoding/base64"
)

//CBC加密
func EncryptToCBC(src, key string) (string, error) {
	data := []byte(src)
	keyByte := []byte(key)
	block, err := des.NewCipher(keyByte)
	if err != nil {
		return "", err
	}
	data = PKCS5Padding(data, block.BlockSize())
	//data = ZeroPadding(data, block.BlockSize())

	blockMode := cipher.NewCBCEncrypter(block, keyByte)
	crypted := make([]byte, len(data))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	blockMode.CryptBlocks(crypted, data)

	return base64.StdEncoding.EncodeToString(crypted), nil
}

func DecryptToCBC(encodeString string, key []byte) (string, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	//base64解密
	crypted, err := base64.StdEncoding.DecodeString(encodeString)

	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return string(origData), nil
}
