package server

import (
	"context"
	"fmt"
	"github.com/evanyxw/game_proto/msg"
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/internal/rpc/core"
	"google.golang.org/grpc"
)

// WorldServer 的实现
type WorldServer struct {
	msg.UnimplementedWorldServer
	core.RpcServer
	// 可以添加一些字段
	Name string
}

func NewWorldServer() *WorldServer {
	return &WorldServer{
		Name: "world",
	}
}

func (s *WorldServer) RegisterService(grpcServer *grpc.Server) {
	msg.RegisterWorldServer(grpcServer, &WorldServer{})
}

func (s *WorldServer) GetServiceName() string {
	return s.Name
}

// Broadcast
func (s *WorldServer) Broadcast(ctx context.Context, req *msg.Req) (*msg.Res, error) {
	fmt.Println("Broadcasting")
	// 实现 Test 方法的逻辑
	return &msg.Res{}, nil
}

func (s *WorldServer) Run() {
	config := configs.Get()
	s.RpcServer.Start(config.Rpc.Address, &WorldServer{})
}
