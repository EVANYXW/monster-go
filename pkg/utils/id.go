/**
 * @api post utils.
 *
 * User: yunshengzhu
 * Date: 2019-04-25
 * Time: 17:22
 */
package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var (
	ErrFormatInvalid      = errors.New("格式错误")
	ErrAddressInvalid     = errors.New("地址码错误")
	ErrBirthFormatInvalid = errors.New("出生日期格式错误")
	ErrBirthRangeInvalid  = errors.New("出生日期范围错误")
	ErrSumInvalid         = errors.New("校验和错误")

	reg      = regexp.MustCompile("^(\\d{6})(18|19|20)?(\\d{2})(0\\d|10|11|12)([012]\\d|3[01])(\\d{3})(\\d|X)?$")
	area     = map[string]string{"11": "北京", "12": "天津", "13": "河北", "14": "山西", "15": "内蒙", "21": "辽宁", "22": "吉林", "23": "黑龙", " 31": "上海", "32": "江苏", "33": "浙江", "34": "安徽", "35": "福建", "36": "江西", "37": "山东", "41": "河南", "42": "湖北", "43": "湖南", "44": "广东", "45": "广西", "46": "海南", "50": "重庆", "51": "四川", "52": "贵州", "53": "云南", "54": "西藏", "61": "陕西", "62": "甘肃", "63": "青海", "64": "宁夏", "65": "新疆", "71": "台湾", "81": "香港", "82": "澳门", "91": "国外"}
	min_date = time.Date(1890, 0, 0, 0, 0, 0, 0, time.Local)
	max_date = time.Now()
	weight   = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2} //十七位数字本体码权重
	code     = []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
)

// 整体校验格式
func ValidateReg(id string) error {
	if reg.MatchString(id) {
		return nil
	}
	return ErrFormatInvalid
}

// 校验地区码
func ValidateArea(id string) error {
	if _, ok := area[id[0:2]]; ok {
		return nil
	}
	return ErrAddressInvalid
}

// 校验生日,包括格式和范围
func ValidateBirth(id string) error {
	birth := id[6:14]
	if date, err := time.Parse("20060102", birth); err != nil {
		return ErrBirthFormatInvalid
	} else if date.After(max_date) && date.Before(min_date) {
		return ErrBirthRangeInvalid
	}
	return nil
}

// 校验 校验和
func ValidateSum(id string) error {
	sum := 0
	for i, char := range id[:len(id)-1] {
		char_f, _ := strconv.ParseFloat(string(char), 64)
		sum += int(char_f) * weight[i]
	}
	if code[sum%11] == id[len(id)-1] {
		return nil
	}
	return ErrSumInvalid
}

// 校验
func Validate(id string) (flag bool, err error) {
	if err = ValidateReg(id); err != nil {
		return false, err
	}
	if err = ValidateArea(id); err != nil {
		return false, err
	}
	if err = ValidateBirth(id); err != nil {
		return false, err
	}
	if err = ValidateSum(id); err != nil {
		return false, err
	}
	return true, nil
}

// 验证生日,返回生日string
func CheckBirthday(id string) int64 {
	birthday := year(id) + "-" + month(id) + "-" + day(id)
	t, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	if t.Unix() > time.Now().Unix() {
		fmt.Println("sorry:不能识别来自未来的人!")
		return 0
	}
	return t.Unix()
}

// 获取年龄.
func Age(id string) int {
	birthdayTime := time.Unix(CheckBirthday(id), 0)
	return timeSub(time.Now(), birthdayTime) / 365
}

// 2个时间相差天数
func timeSub(t1, t2 time.Time) int {
	t1 = t1.UTC().Truncate(24 * time.Hour)
	t2 = t2.UTC().Truncate(24 * time.Hour)
	return int(t1.Sub(t2).Hours() / 24)
}

// 获取性别.
func Sex(id string) string {
	v, _ := strconv.Atoi(id[16:17])
	if v%2 != 0 {
		return "男"
	}
	return "女"
}

// 获取年.
func year(id string) string {
	return id[6:10]
}

// 获取月.
func month(id string) string {
	return id[10:12]
}

// 获取日.
func day(id string) string {
	return id[12:14]
}
