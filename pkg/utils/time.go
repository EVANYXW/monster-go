package utils

import (
	"time"
)

//------------------------------------------------ 返回当天的0点
func PM12() int64 {
	t := time.Now().Unix()
	return t - t%86400
}

func UnixPM12() int64 {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	timeNumber := t.Unix() - 3600*8
	return timeNumber
}

//------------------------------------------------ 返回当前的整点时间
func SharpClock() int64 {
	t := time.Now().Unix()
	return t - t%3600
}

func GetDataMonth(timeObj time.Time, monthNum int) int64 {
	year, month, _ := timeObj.Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	start := thisMonth.AddDate(0, monthNum, 0)
	return start.Unix()
}
