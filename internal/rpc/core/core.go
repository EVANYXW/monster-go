package core

import (
	"fmt"
	"github.com/evanyxw/monster-go/pkg/middleware"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

type RpcServer struct {
}

func (r *RpcServer) Start(address string, service IServer) {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	st := grpc.NewServer(grpc.UnaryInterceptor(middleware.UnaryServerInterceptor))

	RegisterGRPCService(st, service)

	serverName := service.GetServiceName()
	fmt.Println(fmt.Sprintf("【 %s rpc 】 server is started", serverName))
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

// 封装注册 gRPC 服务的方法
func RegisterGRPCService(server *grpc.Server, serviceType interface{}) {
	registerMethod := reflect.ValueOf(serviceType).MethodByName("RegisterService")
	registerMethod.Call([]reflect.Value{reflect.ValueOf(server)})
}

//func RegisterPb[T, V, F any](serviceRegistrar T, service V, handler F) {
//	//handler(serviceRegistrar, service)
//	handlerValue := reflect.ValueOf(handler)
//	if handlerValue.Kind() == reflect.Func {
//		//if handlerValue.Type().NumIn() == 2 && handlerValue.Type().In(0) == reflect.TypeOf(serviceRegistrar) && handlerValue.Type().In(1) == reflect.TypeOf(service) {
//		//	handlerValue.Call([]reflect.Value{reflect.ValueOf(serviceRegistrar), reflect.ValueOf(service)})
//		//} else {
//		//	fmt.Println("Handler function does not match the expected signature")
//		//}
//		handlerValue.Call([]reflect.Value{reflect.ValueOf(serviceRegistrar), reflect.ValueOf(service)})
//	} else {
//		fmt.Println("Handler is not a function")
//	}
//}
