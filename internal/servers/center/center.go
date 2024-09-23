package center

import (
	"github.com/evanyxw/monster-go/internal/servers"
	"github.com/evanyxw/monster-go/pkg/module/register-discovery/center"
	"github.com/evanyxw/monster-go/pkg/output"
	"github.com/evanyxw/monster-go/pkg/server/engine"
)

type Center struct {
	*engine.BaseEngine

	//*centerModule.CenterNet
}

func New() engine.IServerKernel {
	//w := &Center{
	//	engine.NewServer(servers.Center),
	//	centerModule.NewCenterNet(module.ModuleID_SM, 10000),
	//}

	//baseEngine := engine.NewServer(
	//	servers.Center,
	//).WithModule(centerModule.NewCenterNet(module.ModuleID_SM, 10000))
	//
	//w := &Center{
	//	baseEngine,
	//}
	//return w
	baseEngine := engine.NewCenterServer(
		servers.Center,
		center.NewFactor(),
	).WithOutput(&output.Config{
		Name: servers.Center,
		Addr: "",
		Url:  "http://",
	})
	//.WithModule(centerModule.NewCenterNet(module.ModuleID_SM, 10000))

	w := &Center{
		baseEngine,
	}
	return w
}
