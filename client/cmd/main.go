package main

import (
	"fmt"
	"github.com/evanyxw/monster-go/client/robot"
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/timeutil"
	"time"
)

func tcpClient() {
	logger.NewLogger(
		logger.WithDisableConsole(),
		logger.WithField("domain", fmt.Sprintf("%s[%s]", "client", "dev")),
		logger.WithTimeLayout(timeutil.CSTLayout),
		logger.WithFileP(configs.LogFile, "client"),
	)
	module.Init(module.ModuleID_Max)
	c := robot.NewRobot()
	module.Run()
	c.Start()
	time.Sleep(10 * time.Hour)
}

func main() {
	tcpClient()
}
