package world

import (
	"fmt"
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/internal/network"
	rpcServer "github.com/evanyxw/monster-go/internal/rpc/server"
	"github.com/evanyxw/monster-go/internal/server/core"
	"github.com/evanyxw/monster-go/internal/server/factory"
	"github.com/evanyxw/monster-go/pkg/env"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/timeutil"
	"go.uber.org/zap"
	"os"
	"syscall"
)

var (
	Logger *zap.Logger
)

type World struct {
	*core.Server
}

var Oasis *World

func initLog() {
	log, err := logger.NewJSONLogger(
		logger.WithDisableConsole(),
		logger.WithField("domain", fmt.Sprintf("%s[%s]", configs.ProjectName, env.Active().Value())),
		logger.WithTimeLayout(timeutil.CSTLayout),
		logger.WithFileP(configs.ProjectCronLogFile),
	)
	if err != nil {
		panic(err)
	}
	Logger = log
}

func New(info network.Info) factory.Server {
	initLog()
	w := &World{
		core.NewServer(info, Logger),
	}

	return w
}

// Run 启动服务
func (w *World) Run() {
	w.HandlerRegister()
	w.Server.Run()

	worldRpcServer := rpcServer.NewWorldServer()
	go worldRpcServer.Run()
}

// Destroy 注销服务
func (w *World) Destroy() {
	w.Server.Destroy()
}

// OnSystemSignal 监听退出信道
func (w *World) OnSystemSignal(signal os.Signal) bool {
	tag := true
	switch signal {
	case syscall.SIGHUP:
		//todo
		fmt.Println("SIGHUP")
	case syscall.SIGPIPE:
		fmt.Println("SIGPIPE")
	default:
		Logger.Debug("[World] 收到信号准备退出 %v \n", zap.String("signal", signal.String()))
		tag = false
	}
	return tag
}
