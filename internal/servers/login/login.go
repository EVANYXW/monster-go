package login

import (
	"github.com/evanyxw/monster-go/internal/servers"
	centerModule "github.com/evanyxw/monster-go/internal/servers/center/module"
	"github.com/evanyxw/monster-go/internal/servers/common/handler"
	commonModule "github.com/evanyxw/monster-go/internal/servers/common/module"
	accHandler "github.com/evanyxw/monster-go/internal/servers/gate/handler"
	loginModule "github.com/evanyxw/monster-go/internal/servers/login/module"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
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

	baseEngine := engine.NewServer(
		servers.Login,
	).WithOutput(&output.Config{
		Name: servers.Login,
		Addr: "",
		Url:  "http://",
	}).
		WithModule(centerModule.NewCenterConnector(
			module.ModuleID_CenterConnector,
			handler.NewServerInfoHandler(),
		)).WithModule(commonModule.NewClientNet(
		module.ModuleID_GateAcceptor,
		10000,
		accHandler.NewAcceptor(),
		module.Inner,
		new(network.ClientPackerFactory),
	)).
		WithModule(loginModule.NewLoginManager(module.ModuleID_LoginManager)).
		WithModule(loginModule.NewLoginConfig(module.ModuleID_LoginConfig)).
		WithModule(commonModule.NewRedisClient(module.ModuleID_Redis))

	return &login{
		baseEngine,
	}
}
