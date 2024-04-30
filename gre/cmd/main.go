package main

import (
	"github.com/evanyxw/monster-go/gre/client"
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
