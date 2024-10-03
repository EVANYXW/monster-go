package gate_grpc

import (
	"github.com/evanyxw/monster-go/internal/servers"
	"github.com/evanyxw/monster-go/pkg/module/register-discovery/center"
	"github.com/evanyxw/monster-go/pkg/output"
	"github.com/evanyxw/monster-go/pkg/server/engine"
	"sync/atomic"
)

type Gate struct {
	*engine.BaseEngine
}

type gateServerInfo struct {
	mailID    atomic.Uint32
	managerID atomic.Uint32
}

func New() engine.IServerKernel {

	baseEngine := engine.NewGrpcServer(
		servers.Gate,
		center.NewFactor(center.WithServerConnectorManager()),
		[]engine.GrpcServer{},
	).WithOutput(&output.Config{
		Name: servers.Gate,
		Addr: "",
		Url:  "http://",
	})

	return &Gate{
		baseEngine,
	}
}
