package main

import (
	"github.com/evanyxw/monster-go/gre/client"
	"github.com/phuhao00/sugar"
)

func main() {
	//logger.NewLogger(
	//	logger.WithDisableConsole(),
	//	logger.WithField("domain", fmt.Sprintf("%s[%s]", configs.ProjectName, env.Active().Value())),
	//	logger.WithTimeLayout(timeutil.CSTLayout),
	//	logger.WithFileP(configs.LogFile, "client"),
	//)
	//client.Broadcast()
	//return

	c := client.NewClient()
	c.InputHandlerRegister()
	c.MessageHandlerRegister()
	c.Run()
	sugar.WaitSignal(c.OnSystemSignal)
}
