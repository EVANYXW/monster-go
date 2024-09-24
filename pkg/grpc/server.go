// Package grpcpool @Author evan_yxw
// @Date 2024/9/23 23:17:00
// @Desc
package grpc

import (
	"fmt"
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
	ip          string
	port        int
}

func NewServer(service string, ip string, port int, options ...GrpcOptions) *Server {
	op := &option{
		minPort: _minPort,
		maxPort: _maxPort,
	}
	for _, fn := range options {
		fn(op)
	}

	s := &Server{
		name:     service,
		Log:      op.Log,
		minPort:  op.minPort,
		maxPort:  op.maxPort,
		isTracer: op.isTracer,
		ip:       ip,
		port:     port,
	}

	s.init(op.interceptor...)
	return s
}

func (s *Server) SetPort(minPort, maxPort int) {
	s.minPort = minPort
	s.maxPort = maxPort
}

func (s *Server) init(mdw ...grpc.UnaryServerInterceptor) {
	s.serviceAddr = fmt.Sprintf("%s:%d", s.ip, s.port)
	//s.tcpEtcd = s.registerEtcd(s.etcdClient, s.name, s.serviceAddr)

	defer func() {
		if err := recover(); err != nil {
			if s.Log != nil {
				s.Log.Error("err: ", zap.Any("error", err))
			}
		}
	}()

	listen, err := net.Listen("tcp", s.serviceAddr)
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

	interceptors = append(interceptors, mdw...)

	s.Server = grpc.NewServer()
	s.listener = listen
}
