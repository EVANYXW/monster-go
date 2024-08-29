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
