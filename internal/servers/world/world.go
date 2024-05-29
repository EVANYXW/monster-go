package world

import (
	"fmt"
	"github.com/evanyxw/monster-go/cmd/factory"
	module2 "github.com/evanyxw/monster-go/internal/servers/center/module"
	"github.com/evanyxw/monster-go/internal/servers/core/handler"
	gateHandle "github.com/evanyxw/monster-go/internal/servers/gate/handler"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/server"
	"go.uber.org/zap"
	"os"
	"syscall"
)

type World struct {
	centerConnector *module2.CenterConnector
	clientNet       *module.ClientNet
}

func New(info server.Info) factory.CmdServer {
	w := &World{
		centerConnector: module2.NewCenterConnector(module.ModuleID_CenterConnector, handler.NewServerInfoHandler()),
		clientNet: module.NewClientNet(module.ModuleID_Client, 10000,
			gateHandle.New(false), info, module.Inner),
	}

	return w
}

// Run 外部通知开启Module
func (w *World) Run() {
	module.Run()
	//worldRpcServer := rpcServer.NewWorldServer()
	//go worldRpcServer.Run()
}

// Destroy 注销服务
func (w *World) Destroy() {
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
