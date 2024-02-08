package utils

import (
	"bufio"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"math/big"
	"os"
)

func FenToYuanStr(amount int64) string {
	amountObj := decimal.NewFromBigInt(big.NewInt(amount), 0)
	quantity := decimal.NewFromFloat(100)
	amountYuan := amountObj.Div(quantity).StringFixed(2)
	return amountYuan
}

func FenToYuanFloat(amount int64) float64 {
	amountObj := decimal.NewFromBigInt(big.NewInt(amount), 0)
	quantity := decimal.NewFromFloat(100)
	amountYuan, _ := amountObj.Div(quantity).Float64()
	return amountYuan
}

func YuanToFenStr(amount string) (int64, error) {
	amountObj, err := decimal.NewFromString(amount)
	if err != nil {
		return 0, err
	}
	quantity := decimal.NewFromFloat(100)
	amountFen := amountObj.Mul(quantity).IntPart()
	return amountFen, nil
}

func YuanToFenFloat(amount float64) int64 {
	amountObj := decimal.NewFromFloat(amount)
	quantity := decimal.NewFromFloat(100)
	amountFen := amountObj.Mul(quantity).IntPart()
	return amountFen
}

func FloatMul100Str(str string) string {
	amountObj, _ := decimal.NewFromString(str)
	quantity := decimal.NewFromFloat(100)
	return amountObj.Mul(quantity).StringFixed(2)
}

func FloatMul100(amount float64) float64 {
	amountObj := decimal.NewFromFloat(amount)
	quantity := decimal.NewFromFloat(100)
	floatMul100, _ := amountObj.Mul(quantity).Float64()
	return floatMul100
}

func FloatDiv100(amount float64) float64 {
	amountObj := decimal.NewFromFloat(amount)
	quantity := decimal.NewFromFloat(100)
	floatDiv100, _ := amountObj.Div(quantity).Float64()
	return floatDiv100
}

func FloatDiv100Str(amount float64) string {
	amountObj := decimal.NewFromFloat(amount)
	quantity := decimal.NewFromFloat(100)
	return amountObj.Div(quantity).StringScaled(2)
}

// 文件读取
func ReadFile(fileName string) ([]byte, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	result, err := ioutil.ReadAll(bufio.NewReader(f))
	if err != nil {
		return nil, err
	}
	return result, nil
}
