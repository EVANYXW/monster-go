package core

import "google.golang.org/grpc"

type IServer interface {
	RegisterService(grpcServer *grpc.Server)
	GetServiceName() string
	Run()
}
