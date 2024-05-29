package center

import (
	"fmt"
	"github.com/evanyxw/monster-go/cmd/factory"
	centerModule "github.com/evanyxw/monster-go/internal/servers/center/module"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/server"
	"go.uber.org/zap"
	"os"
	"syscall"
)

type Center struct {
	centerNet *centerModule.CenterNet
}

func New(info server.Info) factory.CmdServer {
	w := &Center{
		centerNet: centerModule.NewCenterNet(module.ModuleID_SM, 10000, info),
	}

	return w
}

// Run 外部通知开启Module
func (w *Center) Run() {
	module.Run()

}

// Destroy 注销服务
func (w *Center) Destroy() {
	module.Close()
}

// OnSystemSignal 监听退出信道
func (w *Center) OnSystemSignal(signal os.Signal) bool {
	tag := true
	switch signal {
	case syscall.SIGHUP:
		//todo
		fmt.Println("SIGHUP")
	case syscall.SIGPIPE:
		fmt.Println("SIGPIPE")
	default:
		logger.Info("【 center 】 收到信号准备退出", zap.String("signal", signal.String()))
		tag = false
	}
	return tag
}
