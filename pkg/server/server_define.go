package server

import "google.golang.org/grpc"

type GrpcServer interface {
	TransportRegister() func(grpc.ServiceRegistrar) error
}

// center 启动其他服务器的状态
const (
	CN_RunStep_None = iota
	CN_RunStep_StartServer
	CN_RunStep_WaitHandshake
	CN_RunStep_HandshakeDone
	CN_RunStep_Done
)
