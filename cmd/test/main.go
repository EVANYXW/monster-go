// Package test @Author evan_yxw
// @Date 2024/9/13 22:49:00
// @Desc
package main

import (
	"fmt"
	"time"
)

func pri() {
	time.Sleep(time.Second * 1)
	fmt.Println("11111")
}

func test() {
	go pri()
	//time.Sleep(time.Second * 8)
}

func main() {
	test()
	fmt.Println(2222)
	time.Sleep(time.Second * 10)
}
