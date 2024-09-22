package gate

import (
	"github.com/evanyxw/monster-go/internal/servers"
	commonModule "github.com/evanyxw/monster-go/internal/servers/common/module"
	gateHandler "github.com/evanyxw/monster-go/internal/servers/gate/handler"
	"github.com/evanyxw/monster-go/pkg/module"
	register_discovery "github.com/evanyxw/monster-go/pkg/module/register-discovery/center"
	centerHandler "github.com/evanyxw/monster-go/pkg/module/register-discovery/center/handler"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/output"
	"github.com/evanyxw/monster-go/pkg/server/engine"
	"sync/atomic"
)

type Gate struct {
	*engine.BaseEngine
}

type gateServerInfo struct {
	mailID    atomic.Uint32
	managerID atomic.Uint32
}

func New() engine.IServerKernel {
	//baseEngine := engine.NewServer(
	//	servers.Gate,
	//	register_discovery.NewFactor(register_discovery.WithServerConnectorManager()),
	//).WithOutput(&output.Config{
	//	Name: servers.Gate,
	//	Addr: "",
	//	Url:  "http://",
	//}).WithModule(register_discovery.NewCenterConnector(
	//	module.ModuleID_CenterConnector,
	//	centerHandler.NewGateServerInfoHandler(),
	//)).WithModule(commonModule.NewClientNet(
	//	module.ModuleID_Client,
	//	5000,
	//	gateHandler.NewGateMsg(),
	//	module.Outer,
	//	new(network.DefaultPackerFactory),
	//)).WithModule(manager.NewConnectorManager(module.ModuleID_ConnectorManager))
	//
	//return &Gate{
	//	baseEngine,
	//}

	baseEngine := engine.NewServer(
		servers.Gate,
		register_discovery.NewFactor(register_discovery.WithServerConnectorManager()),
	).WithOutput(&output.Config{
		Name: servers.Gate,
		Addr: "",
		Url:  "http://",
	}).WithModule(register_discovery.NewCenterConnector(
		module.ModuleID_CenterConnector,
		centerHandler.NewGateServerInfoHandler(),
	)).WithModule(commonModule.NewClientNet(
		module.ModuleID_Client,
		5000,
		gateHandler.NewGateMsg(),
		module.Outer,
		new(network.DefaultPackerFactory),
	))

	return &Gate{
		baseEngine,
	}
}
