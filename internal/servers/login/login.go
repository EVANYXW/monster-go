package login

import (
	"github.com/evanyxw/monster-go/internal/servers"
	commonModule "github.com/evanyxw/monster-go/internal/servers/common/module"
	accHandler "github.com/evanyxw/monster-go/internal/servers/gate/handler"
	loginModule "github.com/evanyxw/monster-go/internal/servers/login/module"
	"github.com/evanyxw/monster-go/pkg/module"
	register_discovery "github.com/evanyxw/monster-go/pkg/module/register-discovery/center"
	centerHandler "github.com/evanyxw/monster-go/pkg/module/register-discovery/center/handler"
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
	baseEngine := engine.NewServer(
		servers.Login,
		register_discovery.NewFactor(),
	).WithOutput(&output.Config{
		Name: servers.Login,
		Addr: "",
		Url:  "http://",
	}).
		WithModule(register_discovery.NewCenterConnector(
			module.ModuleID_CenterConnector,
			centerHandler.NewServerInfoHandler(),
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
