package gate

import (
	"fmt"
	"github.com/evanyxw/monster-go/cmd/factory"
	centerModule "github.com/evanyxw/monster-go/internal/servers/center/module"
	gateHandle "github.com/evanyxw/monster-go/internal/servers/gate/handler"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/server"
	"go.uber.org/zap"
	"os"
	"syscall"
)

type Gate struct {
	centerConnector  *centerModule.CenterConnector // 中心服务器连接器
	clientNet        *module.ClientNet             // 用户端网络模块
	connectorManager *module.ConnectorManager      // 其他服务器链接管理器
}

func New(info server.Info) factory.CmdServer {
	w := &Gate{
		centerConnector: centerModule.NewCenterConnector(module.ModuleID_CenterConnector,
			gateHandle.NewGateServerInfoHandler()),
		clientNet: module.NewClientNet(module.ModuleID_Client, 10000,
			gateHandle.New(true), info, module.Outer),
		connectorManager: module.NewConnectorManager(module.ModuleID_ConnectorManager),
	}

	return w
}

// Run 外部告诉内部服务器启动
func (w *Gate) Run() {
	module.Run()
}

// Destroy 外部通知内部注销关闭，信号量
func (w *Gate) Destroy() {
	module.Close()
}

// OnSystemSignal 外部通知内部,监听退出信道
func (w *Gate) OnSystemSignal(signal os.Signal) bool {
	tag := true
	switch signal {
	case syscall.SIGHUP:
		//todo
		fmt.Println("SIGHUP")
	case syscall.SIGPIPE:
		fmt.Println("SIGPIPE")
	default:
		logger.Debug("【 gate 】 收到信号准备退出", zap.String("signal", signal.String()))
		tag = false
	}
	return tag
}
