package gate

import (
	"github.com/evanyxw/monster-go/internal/servers"
	centerModule "github.com/evanyxw/monster-go/internal/servers/center/module"
	commonModule "github.com/evanyxw/monster-go/internal/servers/common/module"
	gateHandler "github.com/evanyxw/monster-go/internal/servers/gate/handler"
	"github.com/evanyxw/monster-go/internal/servers/gate/manager"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/server/engine"
	"sync/atomic"
)

type Gate struct {
	*engine.BaseEngine
	*centerModule.CenterConnector // 中心服务器连接器
	*commonModule.ClientNet       // 用户端网络模块
	*manager.ConnectorManager     // 其他服务器链接管理器
}

type gateServerInfo struct {
	mailID    atomic.Uint32
	managerID atomic.Uint32
}

func New() engine.IServerKernel {
	w := &Gate{
		engine.NewServer(servers.Gate),
		// center 服务连接器
		centerModule.NewCenterConnector(
			module.ModuleID_CenterConnector,
			gateHandler.NewGateServerInfoHandler(),
		),
		// 网关客户端tcp服务端口
		commonModule.NewClientNet(
			module.ModuleID_Client,
			5000,
			gateHandler.NewGateMsg(),
			module.Outer,
			new(network.DefaultPackerFactory),
		),
		// 其他内部服务的连接器
		manager.NewConnectorManager(module.ModuleID_ConnectorManager),
	}
	return w
}
