package utils

import (
	"strings"
	"crypto/md5"
	"encoding/hex"
)

func GetReplaceSql(str string) string {
	str = strings.Replace(str, "'", "''", -1)
	return str
}

func GetReplaceSqlMysql(str string) string {
	str = strings.Replace(str, "'", "\\'", -1)
	return str
}

func MyMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
