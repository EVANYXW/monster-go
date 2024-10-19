package gate

import (
	"github.com/evanyxw/monster-go/internal/servers"
	"github.com/evanyxw/monster-go/pkg/module/register-discovery/center"
	"github.com/evanyxw/monster-go/pkg/output"
	"github.com/evanyxw/monster-go/pkg/server/engine"
)

type Gate struct {
	*engine.BaseEngine
}

func New() engine.IServerKernel {
	baseEngine := engine.NewGateTcpServer(
		servers.Gate,
		center.NewFactor(),
		//etcd.NewFactor(),
	).WithOutput(&output.Config{
		Name: servers.Gate,
		Addr: "",
		Url:  "http://",
	})

	return &Gate{
		baseEngine,
	}
}
