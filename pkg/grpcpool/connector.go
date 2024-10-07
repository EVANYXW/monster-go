// Package grpcpool @Author evan_yxw
// @Date 2024/7/5 19:19:00
// @Desc
package grpcpool

import (
	"context"
	"fmt"
	"github.com/evanyxw/monster-go/pkg/etcdv3"
	"github.com/evanyxw/monster-go/pkg/logger"
	//grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	//grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Connector struct {
	*etcdv3.Etcd
	logger     *zap.Logger
	servername string
	connOnce   sync.Once
	//connOnceMap sync.Map
	connMap  map[string]*grpc.ClientConn
	isTracer bool
}

var (
	connector *Connector
	once      sync.Once
)

// NewConnector 创建grpc服务器连接器
func NewConnector(servername string, e *etcdv3.Etcd, logger *zap.Logger, options ...GrpcOptions) *Connector {
	op := &option{}
	for _, fn := range options {
		fn(op)
	}

	once.Do(func() {
		c := &Connector{
			Etcd:       e,
			servername: servername,
			logger:     logger,
			connMap:    make(map[string]*grpc.ClientConn),
			isTracer:   op.isTracer,
		}
		connector = c
	})

	return connector
}

func Instance() *Connector {
	return connector
}

func (c *Connector) DialAddr(servername string, addr string, opts ...grpc.DialOption) *grpc.ClientConn {
	if _, ok := c.connMap[servername]; !ok {
		//ttarget := fmt.Sprintf("router://%s", target)
		opts = append(opts,
			grpc.WithInsecure(),
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"Balance_m3g"}`),
			//grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient()),
			grpc.WithTimeout(time.Second*10),
		)
		conn, err := grpc.Dial(addr, opts...)
		if err != nil {
			panic(err)
		}
		c.connMap[servername] = conn
	}

	return c.connMap[servername]
}

func (c *Connector) GetConn() *etcdv3.Etcd {
	return c.Etcd
}

func (c *Connector) Dial(service string, mdw ...grpc.UnaryClientInterceptor) *grpc.ClientConn {
	if _, ok := c.connMap[service]; !ok {
		ctx := context.Background()
		opts := make([]grpc.DialOption, 0)
		//rs := etcdv3.NewResolver(grpc_client.EtcdCnf, grpc_client.UserService)
		rs := etcdv3.NewResolver(c.GetConn(), service)
		opts = append(opts, grpc.WithResolvers(rs))
		opts = append(opts, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
		opts = append(opts, grpc.WithChainUnaryInterceptor(mdw...))

		if c.isTracer {
			//mdw = append(mdw, grpc_opentracing.UnaryClientInterceptor()) // 添加 opentracing 的 UnaryClientInterceptor
		}

		//opts = append(opts, grpc.WithUnaryInterceptor(grpc_client.Connector.Interceptor))
		//opts = append(opts, grpc.WithUnaryInterceptor(grpc_opentracing.UnaryClientInterceptor())) // 添加 opentracing 的 UnaryClientInterceptor
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		conn, err := grpc.DialContext(ctx, rs.Scheme()+"://authority/"+service, opts...)
		if err != nil {
			panic(err)
		}
		c.connMap[service] = conn
	}

	return c.connMap[service]
}

func (c *Connector) UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	reqIds := md["reqid"]
	reqId := ""
	if len(reqIds) > 0 {
		reqId = reqIds[0]
	}

	logLevel := false
	loglevel := md["middlewareloglevel"]
	if len(loglevel) > 0 {
		if loglevel[0] == "true" {
			logLevel = true
		}
	}

	resp, err := handler(ctx, req)
	if logLevel {
		logger.With(zap.String("reqId", reqId))
		logger.With(zap.String("duration", time.Since(start).String()))
		logger.With(zap.String("method", info.FullMethod))
		logger.With(zap.Any("req", req))
		logger.With(zap.Any("resp", resp))
		if c.servername != "" {
			logger.With(zap.String("servername", c.servername))
		}
		if err != nil {
			logger.With(zap.String("err", err.Error()))
		}
		logger.Info("server after handling.")
	}
	return resp, err
}

// Interceptor 客户端拦截器
func (c *Connector) Interceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	reqIds := md["reqid"]
	reqId := ""
	if len(reqIds) > 0 {
		reqId = reqIds[0]
	}

	logLevel := false
	loglevel := md["middlewareloglevel"]
	if len(loglevel) > 0 {
		if loglevel[0] == "true" {
			logLevel = true
		}
	}

	err := invoker(ctx, method, req, reply, cc, opts...)
	if logLevel {
		with := logger.With(zap.String("reqId", reqId))
		with = logger.With(zap.String("duration", time.Since(start).String()))
		with = logger.With(zap.String("method", method))
		with = logger.With(zap.Any("req", req))

		if c.servername != "" {
			with = logger.With(zap.String("client server:", c.servername))
		}
		with = logger.With(zap.Any("res", reply))
		if err != nil {
			with = logger.With(zap.Any("err", err.Error()))
		}
		with.Info("client after handling.")
	}
	return err
}

func (c *Connector) Release() {
	for _, client := range c.connMap {
		client.Close()
	}
}
