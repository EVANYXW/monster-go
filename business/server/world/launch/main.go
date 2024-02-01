package main

import (
	"bilibili/monster-go/business/server/world"
	"github.com/phuhao00/sugar"
)

func main() {
	wd := world.NewWorld()
	wd.Start()
	//logger.Logger.InfoF("server start !!")
	sugar.WaitSignal(world.Oasis.OnSystemSignal)
	//ch := make(chan os.Signal, 1)
	//signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGPIPE)
	select {}
}
