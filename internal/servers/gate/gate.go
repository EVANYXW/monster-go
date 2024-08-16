package gate

import (
	"github.com/evanyxw/monster-go/internal/servers"
	centerModule "github.com/evanyxw/monster-go/internal/servers/center/module"
	commonModule "github.com/evanyxw/monster-go/internal/servers/common/module"
	gateHandler "github.com/evanyxw/monster-go/internal/servers/gate/handler"
	"github.com/evanyxw/monster-go/internal/servers/gate/manager"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/server"
	"github.com/evanyxw/monster-go/pkg/server/engine"
)

type Gate struct {
	*engine.BaseEngine

	*centerModule.CenterConnector // 中心服务器连接器
	*commonModule.ClientNet       // 用户端网络模块
	*manager.ConnectorManager     // 其他服务器链接管理器
}

func New(info server.Info) engine.Kernel {
	gateMsgHandler := gateHandler.New()
	w := &Gate{
		engine.NewEngine(servers.Gate),
		// center 服务连接器
		centerModule.NewCenterConnector(
			module.ModuleID_CenterConnector,
			gateHandler.NewGateServerInfoHandler(),
		),
		// 网关客户端tcp服务端口
		commonModule.NewClientNet(
			module.ModuleID_Client,
			5000,
			gateMsgHandler,
			info,
			module.Outer,
			new(network.DefaultPackerFactory),
		),
		// 其他服务的连接器
		manager.NewConnectorManager(module.ModuleID_ConnectorManager),
	}
	//servers.ClientManager = gateMsgHandler.ClientManager
	return w
}
