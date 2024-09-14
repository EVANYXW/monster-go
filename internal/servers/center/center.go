package center

import (
	"github.com/evanyxw/monster-go/internal/servers"
	centerModule "github.com/evanyxw/monster-go/internal/servers/center/module"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/output"
	"github.com/evanyxw/monster-go/pkg/server"
	"github.com/evanyxw/monster-go/pkg/server/engine"
)

type Center struct {
	*engine.BaseEngine

	*centerModule.CenterNet
}

func New(info server.Info) engine.IServerKernel {
	w := &Center{
		engine.NewEngine(
			servers.Center,
			engine.WithOutput(&output.Config{Name: "center", Addr: "", Url: "http://"}),
			//engine.WithModule(module.ModuleID_SM, centerModule.NewCenterNet(module.ModuleID_SM, 10000)),
		),
		centerModule.NewCenterNet(module.ModuleID_SM, 10000),
	}

	return w
}
