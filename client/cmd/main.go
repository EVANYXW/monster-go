package main

import (
	"fmt"
	"github.com/evanyxw/monster-go/client/robot"
	"time"
)

func main() {
	fmt.Println(111)
	c := robot.NewRobot()
	c.Start()
	time.Sleep(10 * time.Second)
}
