package center

import (
	"github.com/evanyxw/monster-go/internal/servers"
	centerModule "github.com/evanyxw/monster-go/internal/servers/center/module"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/server"
	"github.com/evanyxw/monster-go/pkg/server/engine"
)

type Center struct {
	*engine.BaseEngine

	*centerModule.CenterNet
}

func New(info server.Info) engine.Kernel {
	w := &Center{
		engine.NewEngine(servers.Center),
		centerModule.NewCenterNet(module.ModuleID_SM, 10000, info),
	}

	return w
}

// Run 外部通知开启Module
//func (w *Center) Run() {
//	module.Run()
//
//}
//
//// Destroy 注销服务
//func (w *Center) Destroy() {
//	module.Close()
//}
//
//// OnSystemSignal 监听退出信道
//func (w *Center) OnSystemSignal(signal os.Signal) bool {
//	return engine.BaseSystemSignal(signal, "Center")
//}
