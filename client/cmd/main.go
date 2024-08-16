package main

import (
	"github.com/evanyxw/monster-go/client/robot"
	"time"
)

func main() {
	c := robot.NewRobot()
	c.Start()
	time.Sleep(10 * time.Second)
}
