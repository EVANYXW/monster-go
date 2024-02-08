package server

import (
	"context"
	"fmt"
	"github.com/evanyxw/game_proto/msg"
	"github.com/evanyxw/monster-go/pkg/middleware"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// LogicServiceServer 的实现
type MyLogicServiceServer struct {
	msg.UnimplementedLogicServiceServer
	// 可以添加一些字段
}

// 确保实现了 Test 方法
func (s *MyLogicServiceServer) Test(ctx context.Context, req *msg.Req) (*msg.Res, error) {
	// 实现 Test 方法的逻辑
	return &msg.Res{}, nil
}

func (s *MyLogicServiceServer) Run() {
	listen, err := net.Listen("tcp", ":8024")
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	st := grpc.NewServer(grpc.UnaryInterceptor(middleware.UnaryServerInterceptor))
	msg.RegisterLogicServiceServer(st, &MyLogicServiceServer{})
	fmt.Println("【 logic rpc 】 server is started")

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
