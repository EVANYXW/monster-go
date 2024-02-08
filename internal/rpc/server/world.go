package server

import (
	"context"
	"fmt"
	"github.com/evanyxw/game_proto/msg"
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/pkg/middleware"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// WorldServer 的实现
type WorldServer struct {
	msg.UnimplementedWorldServer
	// 可以添加一些字段
}

// Broadcast
func (s *WorldServer) Broadcast(ctx context.Context, req *msg.Req) (*msg.Res, error) {
	fmt.Println("Broadcasting")
	// 实现 Test 方法的逻辑
	return &msg.Res{}, nil
}

func (s *WorldServer) Run() {
	config := configs.Get()

	listen, err := net.Listen("tcp", config.Rpc.Address)
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	st := grpc.NewServer(grpc.UnaryInterceptor(middleware.UnaryServerInterceptor))
	msg.RegisterWorldServer(st, &WorldServer{})
	fmt.Println("【 world rpc 】 server is started")
	go func() {
		err = st.Serve(listen)
		if err != nil {
			panic(err)
		}
	}()

	// 监听操作系统的终止信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号
	<-sigCh

	fmt.Println("Shutting down server...")

	// 优雅关闭 gRPC 服务器
	st.GracefulStop()

	fmt.Println("Server is gracefully stopped")

}
