package main

import (
	"bilibili/monster-go/gre/client"
	"github.com/phuhao00/sugar"
)

func main() {
	//client.Broadcast()
	//return
	c := client.NewClient()
	c.InputHandlerRegister()
	c.MessageHandlerRegister()
	c.Run()
	sugar.WaitSignal(c.OnSystemSignal)
}
