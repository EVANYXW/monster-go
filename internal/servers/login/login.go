package login

import (
	"github.com/evanyxw/monster-go/internal/servers"
	accHandler "github.com/evanyxw/monster-go/internal/servers/gate/handler"
	loginModule "github.com/evanyxw/monster-go/internal/servers/login/module"
	commonModule "github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/module/register-discovery/etcd"
	"github.com/evanyxw/monster-go/pkg/output"
	"github.com/evanyxw/monster-go/pkg/server/engine"
)

type login struct {
	*engine.BaseEngine

	//*centerModule.CenterConnector
	//*commonModule.ClientNet
	//*loginModule.LoginManager
	//*loginModule.LoginConfig
	//*commonModule.RedisClient
}

func New() engine.IServerKernel {
	//w := &login{
	//	engine.NewServer(servers.Login),
	//	centerModule.NewCenterConnector(
	//		module.ModuleID_CenterConnector,
	//		handler.NewServerInfoHandler(),
	//	),
	//	commonModule.NewClientNet(
	//		module.ModuleID_GateAcceptor,
	//		10000,
	//		accHandler.NewAcceptor(),
	//		module.Inner,
	//		new(network.ClientPackerFactory),
	//	),
	//	loginModule.NewLoginManager(module.ModuleID_LoginManager),
	//	loginModule.NewLoginConfig(module.ModuleID_LoginConfig),
	//	commonModule.NewRedisClient(module.ModuleID_Redis),
	//}
	//
	//return w

	//baseEngine := engine.NewTcpServer(
	//	servers.Login,
	//	center.NewFactor(),
	//).WithOutput(&output.Config{
	//	Name: servers.Login,
	//	Addr: "",
	//	Url:  "http://",
	//}).
	//	//.WithModule(centerModule.NewCenterConnector(
	//	//		module.ModuleID_CenterConnector,
	//	//		handler.NewServerInfoHandler(),
	//	//	))
	//	WithModule(commonModule.NewClientNet(
	//		module.ModuleID_GateAcceptor,
	//		10000,
	//		accHandler.NewAcceptor(),
	//		module.Inner,
	//		new(network.ClientPackerFactory),
	//	)).
	//	WithModule(loginModule.NewLoginManager(module.ModuleID_LoginManager)).
	//	WithModule(loginModule.NewLoginConfig(module.ModuleID_LoginConfig)).
	//	WithModule(commonModule.NewRedisClient(module.ModuleID_Redis))
	//
	//return &login{
	//	baseEngine,
	//}

	baseEngine := engine.NewTcpServer(
		servers.Login,
		accHandler.NewAcceptor(),
		//center.NewFactor(),
		etcd.NewFactor(),
	).WithOutput(&output.Config{
		Name: servers.Login,
		Addr: "",
		Url:  "http://",
	}).
		WithModule(loginModule.NewLoginManager()).
		WithModule(loginModule.NewLoginConfig()).
		WithModule(commonModule.NewRedisClient())

	return &login{
		baseEngine,
	}
}
