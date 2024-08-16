package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func GetVerifyCode(num int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	var vcode string
	switch num {
	case 1:
		vcode = fmt.Sprintf("%01v", rnd.Int31n(10))
	case 2:
		vcode = fmt.Sprintf("%02v", rnd.Int31n(100))
	case 3:
		vcode = fmt.Sprintf("%03v", rnd.Int31n(1000))
	case 4:
		vcode = fmt.Sprintf("%04v", rnd.Int31n(10000))
	case 5:
		vcode = fmt.Sprintf("%05v", rnd.Int31n(100000))
	case 6:
		vcode = fmt.Sprintf("%06v", rnd.Int31n(1000000))
	default:
		vcode = fmt.Sprintf("%06v", rnd.Int31n(1000000))
	}
	return vcode
}

const DATE_FORMAT = "2006-01-02 15:04:05"
const YYMMDDHHII_DATE_FORMAT = "2006-01-02 15:04"
const YYMMDDHHII_FORMAT = "0601021504"
const DATE_FORMAT_YMD = "2006-01-02"
const DATE_FORMAT_YMD1 = "20060102"
const DATE_FORMAT_YM = "2006-01"
const DATE_FORMAT_MDY = "01-02-06"
const DATE_FORMAT_YDM = "01/02/2006"
const DATE_FORMAT_MDYHIS = "01/02/2006 15:04:05"

func DateString(s int64, format string) string {
	if s <= 0 {
		return ""
	}
	tm := time.Unix(s, 0)
	return tm.Format(format)
}

func DataUnix(str, format string) int64 {
	if str == "" {
		return 0
	}
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(format, str, loc)
	tm := theTime.Unix()
	if tm <= 0 {
		return 0
	}
	return tm
}

func StrToDate(str, format string) time.Time {
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(format, str, loc)
	return theTime
}

func CreateOrderNumber(num int) string {
	var str string
	if num > 14 {
		str = GetRandomStr(num - 14)
	}
	tm := time.Unix(time.Now().Unix(), 0)
	return fmt.Sprintf("%s%s", tm.Format("20060102150405"), str)
}

func GetRandomStr(num int) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < num; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GetRandomStrToken(num int) string {
	str := "0123456789qwertyuiopasdfghjklzxcvbnm"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < num; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

const DATE_FORMAT_YMD01 = "2006-1-02"
const DATE_FORMAT_YMD02 = "2006-1-2"
const DATE_FORMAT_YMD03 = "2006-01-2"

const DATE_FORMAT_HIS = "150405"
const DATE_FORMAT_MD = "0102"

func ImportDataUnix(str string) int64 {
	if str == "" {
		return 0
	}
	return dataToInt(str)
}

var dataFormMap map[int]string

func init() {
	dataFormMap = make(map[int]string)
	dataFormMap[1] = "2006-01-02"
	dataFormMap[2] = "2006-01-2"
	dataFormMap[3] = "2006-1-02"
	dataFormMap[4] = "2006-1-2"
	dataFormMap[5] = "06-01-02"
	dataFormMap[6] = "06-01-2"
	dataFormMap[7] = "06-1-02"
	dataFormMap[8] = "06-1-2"

	dataFormMap[9] = "2006/01/02"
	dataFormMap[10] = "2006/01/2"
	dataFormMap[11] = "2006/1/02"
	dataFormMap[12] = "2006/1/2"
	dataFormMap[13] = "06/01/02"
	dataFormMap[14] = "06/01/2"
	dataFormMap[15] = "06/1/02"
	dataFormMap[16] = "06/1/2"

	dataFormMap[17] = "2006.01.02"
	dataFormMap[18] = "2006.01.2"
	dataFormMap[19] = "2006.1.02"
	dataFormMap[20] = "2006.1.2"
	dataFormMap[21] = "06.01.02"
	dataFormMap[22] = "06.01.2"
	dataFormMap[23] = "06.1.02"
	dataFormMap[24] = "06.1.2"

	dataFormMap[25] = "01-02-2006"
	dataFormMap[26] = "01-2-2006"
	dataFormMap[27] = "1-02-2006"
	dataFormMap[28] = "1-2-2006"
	dataFormMap[29] = "01-02-06"
	dataFormMap[30] = "01-2-06"
	dataFormMap[31] = "1-02-06"
	dataFormMap[32] = "1-2-06"

	dataFormMap[33] = "01/02/2006"
	dataFormMap[34] = "01/2/2006"
	dataFormMap[35] = "1/02/2006"
	dataFormMap[36] = "1/2/2006"
	dataFormMap[37] = "01/02/06"
	dataFormMap[38] = "01/2/06"
	dataFormMap[39] = "1/02/06"
	dataFormMap[40] = "1/2/06"

	dataFormMap[41] = "01.02.2006"
	dataFormMap[42] = "01.2.2006"
	dataFormMap[43] = "1.02.2006"
	dataFormMap[44] = "1.2.2006"
	dataFormMap[45] = "01.02.06"
	dataFormMap[46] = "01.2.06"
	dataFormMap[47] = "1.02.06"
	dataFormMap[48] = "1.2.06"

	dataFormMap[49] = "2006年01月02日"
	dataFormMap[50] = "2006年01月2日"
	dataFormMap[51] = "2006年1月02日"
	dataFormMap[52] = "2006年1月2日"
	dataFormMap[53] = "06年01月02日"
	dataFormMap[54] = "06年01月2日"
	dataFormMap[55] = "06年1月02日"
	dataFormMap[56] = "06年1月2日"

}

func dateToInt(str string, num int) int64 {
	if num > 49 {
		return 0
	}
	d := DataUnix(str, dataFormMap[num])
	if d == -62135596800 {
		return dateToInt(str, num+1)
	}
	return d
}

func dataToInt(str string) int64 {
	d := dateToInt(str, 1)
	if d == 0 {
		d = dateToInt(convertToFormatDay(str), 1)
		fmt.Println("d:", d, convertToFormatDay(str))
	} else {
		fmt.Println("d:", d)
	}
	return d
}

func convertToFormatDay(excelDaysString string) string {
	baseDiffDay := 38719
	curDiffDay := excelDaysString
	b, _ := strconv.Atoi(curDiffDay)
	realDiffDay := b - baseDiffDay
	realDiffSecond := realDiffDay * 24 * 3600
	baseOriginSecond := 1136185445
	resultTime := time.Unix(int64(baseOriginSecond+realDiffSecond), 0).Format("2006-01-02")
	return resultTime
}

func StrArrToInt64(list []string) (oList []int64, err error) {
	oList = make([]int64, 0)
	for _, v := range list {
		int64, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return oList, err
		}
		oList = append(oList, int64)

	}
	return oList, nil
}

func ToJsonString(i interface{}) string {
	return string(ToJsonBytes(i))
}

func ToJsonStringIndent(i interface{}) string {
	if IsNil(i) {
		return ""
	}
	bs, er := json.MarshalIndent(i, "", " ")
	if er != nil {
		return ""
	}

	return string(bs)
}

func IsNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}

func ToJsonBytes(i interface{}) []byte {
	if IsNil(i) {
		return nil
	}
	bs, er := json.Marshal(i)
	if er != nil {
		return nil
	}
	return bs
}

