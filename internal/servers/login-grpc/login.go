package login_grpc

import (
	"github.com/evanyxw/monster-go/internal/servers"
	"github.com/evanyxw/monster-go/internal/servers/login-grpc/loginser"
	"github.com/evanyxw/monster-go/pkg/output"
	"github.com/evanyxw/monster-go/pkg/server"
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
	//baseEngine := engine.NewGrpcServer(
	//	servers.LoginGrpc,
	//	[]server.GrpcServer{
	//		loginser.NewLogin(),
	//	},
	//).WithOutput(&output.Config{
	//	Name: servers.LoginGrpc,
	//	Addr: "",
	//	Url:  "http://",
	//})

	baseEngine := engine.NewGrpcServer(
		servers.LoginGrpc,
		[]server.GrpcServer{
			loginser.NewLogin(),
		},
	).WithOutput(&output.Config{
		Name: servers.LoginGrpc,
		Addr: "",
		Url:  "http://",
	})

	return &Gate{
		baseEngine,
	}
}
