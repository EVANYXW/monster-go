package utils

import (
	"math/rand"
	"time"
	"fmt"
	"strings"
)

func getRandom(num int) string {
	var str string
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < num; i++ {
		str += fmt.Sprintf("%d", r.Intn(9))
	}
	return str
}

func GetSerialNumber(fix string, userId int64, num int) string {
	timeLayout := "20060102150405"
	sr := time.Now().Unix()
	dataTimeStr := time.Unix(sr, 0).Format(timeLayout)
	return fix + getRandomSring(userId, 8) + dataTimeStr + getRandom(num)
}

//获取用户36进制
func getRandomSring(userId int64, MinLen int) string {
	s := "0123456789abcdefghijklmnopqrstuvwxyz"
	lenst := int64(len(s))
	var c, m int64
	var str string
	for {
		c = userId % 36
		m = userId / 36
		userId = m
		if c >= 0 && c < lenst {
			str = string(s[c:c+1]) + str
		}
		if m == 0 {
			break
		}
	}
	lens := len(str)
	for i := 1; i <= MinLen-lens; i++ {
		str = "0" + str
	}
	return strings.ToUpper(str)
}
