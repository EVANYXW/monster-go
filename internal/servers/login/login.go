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
		WithModule(module.ModuleID_CenterConnector, register_discovery.NewCenterConnector(
			centerHandler.NewServerInfoHandler(),
		)).WithModule(module.ModuleID_GateAcceptor, commonModule.NewClientNet(
		10000,
		accHandler.NewAcceptor(),
		module.Inner,
		new(network.ClientPackerFactory),
	)).
		WithModule(module.ModuleID_LoginManager, loginModule.NewLoginManager()).
		WithModule(module.ModuleID_LoginConfig, loginModule.NewLoginConfig()).
		WithModule(module.ModuleID_Redis, commonModule.NewRedisClient())

	return &login{
		baseEngine,
	}
}
