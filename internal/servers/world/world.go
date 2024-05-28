package world

import (
	"fmt"
	"github.com/evanyxw/monster-go/cmd/factory"
	"github.com/evanyxw/monster-go/internal/servers/core"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/server"
	"go.uber.org/zap"
	"os"
	"syscall"
)

type World struct {
	centerConnector *core.CenterConnector
	clientNet       *core.ClientNet
}

func New(info server.Info) factory.CmdServer {
	w := &World{
		centerConnector: core.NewCenterConnector(module.ModuleID_CenterConnector, core.NewServerInfoHandler()),
		clientNet:       core.NewClientNet(module.ModuleID_Client, 10000, info, module.Inner),
	}

	return w
}

// Run 外部通知开启Module
func (w *World) Run() {
	//w.CenterConnector.Run()
	module.Run()
	//worldRpcServer := rpcServer.NewWorldServer()
	//go worldRpcServer.Run()
}

// Destroy 注销服务
func (w *World) Destroy() {
	//w.CenterConnector.Release()
	module.Close()
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
		logger.Debug("【 world 】 收到信号准备退出", zap.String("signal", signal.String()))
		tag = false
	}
	return tag
}
