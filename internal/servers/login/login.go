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
	"github.com/evanyxw/monster-go/pkg/server"
	"github.com/evanyxw/monster-go/pkg/server/engine"
)

type login struct {
	*engine.BaseEngine

	*centerModule.CenterConnector
	*commonModule.ClientNet
	*loginModule.LoginManager
}

func New(info server.Info) engine.Kernel {
	w := &login{
		engine.NewEngine(servers.Login),
		centerModule.NewCenterConnector(
			module.ModuleID_CenterConnector,
			handler.NewServerInfoHandler(),
		),
		commonModule.NewClientNet(
			module.ModuleID_GateAcceptor,
			10000,
			accHandler.NewAcceptor(),
			info,
			module.Inner,
			network.NewDefaultPacker(),
		),
		loginModule.NewLoginManager(module.ModuleID_LoginManager),
	}

	return w
}
