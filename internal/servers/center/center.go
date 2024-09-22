package center

import (
	"github.com/evanyxw/monster-go/internal/servers"
	"github.com/evanyxw/monster-go/pkg/module"
	register_discovery "github.com/evanyxw/monster-go/pkg/module/register-discovery/center"
	"github.com/evanyxw/monster-go/pkg/output"
	"github.com/evanyxw/monster-go/pkg/server/engine"
)

type Center struct {
	*engine.BaseEngine

	//*centerModule.CenterNet
}

func New() engine.IServerKernel {
	baseEngine := engine.NewServer(
		servers.Center,
		register_discovery.NewFactor(),
	).
		WithOutput(&output.Config{
			Name: servers.Center,
			Addr: "",
			Url:  "http://",
		}).
		WithModule(register_discovery.NewCenterNet(module.ModuleID_SM, 10000))

	w := &Center{
		baseEngine,
	}

	return w
}
