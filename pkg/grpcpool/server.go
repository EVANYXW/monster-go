// Package grpcpool @Author evan_yxw
// @Date 2024/7/10 16:22:00
// @Desc
package grpcpool

import (
	"fmt"
	"github.com/evanyxw/monster-go/pkg/etcdv3"
	"github.com/evanyxw/monster-go/pkg/ipPort"
	//"github.com/evanyxw/monster-go/pkg/tracer"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

var (
	logFile    = "./logs/"
	defaultEnv = "fat"
)

const (
	_minPort = 18000
	_maxPort = 19000
)

type GrpcOptions func(option *option)

type option struct {
	Log         *zap.Logger
	minPort     int
	maxPort     int
	interceptor []grpc.UnaryServerInterceptor
	isTracer    bool
	addr        string
}

// WithLogger Logger
func WithLogger(logger *zap.Logger) GrpcOptions {
	return func(opt *option) {
		opt.Log = logger
	}
}

// WithPorts port range
func WithPorts(min, max int) GrpcOptions {
	return func(opt *option) {
		opt.minPort = min
		opt.maxPort = max
	}
}

// WithInterceptor 拦截器
func WithInterceptor(interceptor ...grpc.UnaryServerInterceptor) GrpcOptions {
	return func(opt *option) {
		opt.interceptor = interceptor
	}
}

func WithTracer(isTracer bool) GrpcOptions {
	return func(opt *option) {
		opt.isTracer = isTracer
	}
}

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
	mdw         []grpc.UnaryServerInterceptor
}

func NewServer(service string, etcdClient *etcdv3.Etcd, options ...GrpcOptions) *Server {
	op := &option{
		minPort: _minPort,
		maxPort: _maxPort,
	}
	for _, fn := range options {
		fn(op)
	}

	s := &Server{
		name:       service,
		etcdClient: etcdClient,
		Log:        op.Log,
		minPort:    op.minPort,
		maxPort:    op.maxPort,
		isTracer:   op.isTracer,
		mdw:        op.interceptor,
	}

	//s.Connect(op.interceptor...)
	return s
}

func (s *Server) SetPort(minPort, maxPort int) {
	s.minPort = minPort
	s.maxPort = maxPort
}

func (s *Server) GetAddr() string {
	return s.serviceAddr
}

func (s *Server) Connect() {
	ip, port, err := ipPort.GetDynamicIpAndRangePort(s.minPort, s.maxPort)
	//ip = "host.docker.internal"
	if err != nil {
		panic(err)
	}
	s.serviceAddr = fmt.Sprintf("%s:%d", ip, port)
	s.tcpEtcd = s.registerEtcd(s.etcdClient, s.name, s.serviceAddr)

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
	//if s.isTracer {
	//	tracer.NewTracer(s.name)
	//	interceptors = append(interceptors, grpc_opentracing.UnaryServerInterceptor()) // tracer 追踪中间件
	//}

	interceptors = append(interceptors, s.mdw...)

	s.Server = grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(interceptors...)))
	s.listener = listen
}

func (s *Server) registerEtcd(etcd *etcdv3.Etcd, serverName, address string) *etcdv3.Service {
	if s.Log != nil {
		s.Log.Info(fmt.Sprintf("register server 【%s】:%s", serverName, address))
	}

	tcpEtcdServe, err := etcdv3.NewService(etcd, etcdv3.ServiceInfo{Name: serverName, Address: address})
	if err != nil {
		panic(err)
	}
	go func() {
		if err = tcpEtcdServe.Start(); err != nil {
			fmt.Println(err)
		}
	}()

	return tcpEtcdServe
}

func (s *Server) GetListen() net.Listener {
	return s.listener
}

func (s *Server) Release() {
	s.tcpEtcd.Stop()
	s.Stop()
	//tracer.Close()
}

func (s *Server) Run() {
	s.Server.Serve(s.listener)
}
