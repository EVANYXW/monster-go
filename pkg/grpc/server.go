// Package grpcpool @Author evan_yxw
// @Date 2024/9/23 23:17:00
// @Desc
package grpc

import (
	"github.com/evanyxw/monster-go/pkg/etcdv3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

const (
	_minPort = 18000
	_maxPort = 19000
)

type Server struct {
	*grpc.Server
	name       string
	listener   net.Listener
	tcpEtcd    *etcdv3.Service
	etcdClient *etcdv3.Etcd
	Log        *zap.Logger

	serviceAddr string
	minPort     int
	maxPort     int
	isTracer    bool
	addr        string
	mdw         []grpc.UnaryServerInterceptor
}

func WithAddr(addr string) GrpcOptions {
	return func(options *option) {
		options.addr = addr
	}
}

func NewServer(service string) *Server {
	//op := &option{
	//	minPort: _minPort,
	//	maxPort: _maxPort,
	//}
	//for _, fn := range options {
	//	fn(op)
	//}

	s := &Server{
		name: service,
		//Log:      op.Log,
		//minPort:  op.minPort,
		//maxPort:  op.maxPort,
		//isTracer: op.isTracer,
		//mdw:      op.interceptor,
	}

	//s.Connect()
	return s
}

func (s *Server) Connect(options ...GrpcOptions) {
	opt := &option{
		minPort: _minPort,
		maxPort: _maxPort,
	}
	for _, fn := range options {
		fn(opt)
	}

	s.addr = opt.addr
	s.Log = opt.Log
	s.minPort = opt.minPort
	s.maxPort = opt.maxPort
	s.isTracer = opt.isTracer
	s.mdw = opt.interceptor

	//s.tcpEtcd = s.registerEtcd(s.etcdClient, s.name, s.addr)

	defer func() {
		if err := recover(); err != nil {
			if s.Log != nil {
				s.Log.Error("err: ", zap.Any("error", err))
			}
		}
	}()

	listen, err := net.Listen("tcp", s.addr)
	if err != nil {
		panic(err)
	}

	// 组合多个拦截器
	interceptors := []grpc.UnaryServerInterceptor{
		//grpc_opentracing.UnaryServerInterceptor(), // 添加 opentracing 的 Unary Server Interceptor
		//middleware.UnaryServerInterceptor,         // 添加你自定义的中间件拦截器
	}
	if s.isTracer {
		//tracer.NewTracer(s.name)
		//interceptors = append(interceptors, grpc_opentracing.UnaryServerInterceptor()) // tracer 追踪中间件
	}

	interceptors = append(interceptors, s.mdw...)

	s.Server = grpc.NewServer()
	s.listener = listen
}

func (s *Server) Run() {
	s.Server.Serve(s.listener)
}
