// Package grpcpool @Author evan_yxw
// @Date 2024/7/5 19:19:00
// @Desc
package grpc

import (
	"context"
	"google.golang.org/grpc/metadata"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	connector *Connector
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

type Connector struct {
	logger     *zap.Logger
	serverName string
	connOnce   sync.Once
	//connOnceMap sync.Map
	connMap  map[string]*grpc.ClientConn
	isTracer bool
}

// NewServerConnector 创建grpc服务器连接器
func NewServerConnector(serverName string, logger *zap.Logger, options ...GrpcOptions) *Connector {
	op := &option{}
	for _, fn := range options {
		fn(op)
	}

	c := &Connector{
		serverName: serverName,
		logger:     logger,
		connMap:    make(map[string]*grpc.ClientConn),
		isTracer:   op.isTracer,
	}
	connector = c
	return c
}

func GetConnectorInstance() *Connector {
	return connector
}

func (c *Connector) Dial(serverName string, addr string) *grpc.ClientConn {
	if _, ok := c.connMap[serverName]; !ok {
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, addr)
		if err != nil {
			panic(err)
		}
		c.connMap[serverName] = conn
	}

	return c.connMap[serverName]
}

//func (c *Connector) GetConn() *etcdv3.Etcd {
//	return c.Etcd
//}

//func (c *Connector) Dial(service string, mdw ...grpc.UnaryClientInterceptor) *grpc.ClientConn {
//	if _, ok := c.connMap[service]; !ok {
//		ctx := context.Background()
//		opts := make([]grpc.DialOption, 0)
//		//rs := etcdv3.NewResolver(grpc_client.EtcdCnf, grpc_client.UserService)
//		rs := etcdv3.NewResolver(c.GetConn(), service)
//		opts = append(opts, grpc.WithResolvers(rs))
//		opts = append(opts, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
//		opts = append(opts, grpc.WithChainUnaryInterceptor(mdw...))
//
//		if c.isTracer {
//			mdw = append(mdw, grpc_opentracing.UnaryClientInterceptor()) // 添加 opentracing 的 UnaryClientInterceptor
//		}
//
//		//opts = append(opts, grpc.WithUnaryInterceptor(grpc_client.Connector.Interceptor))
//		//opts = append(opts, grpc.WithUnaryInterceptor(grpc_opentracing.UnaryClientInterceptor())) // 添加 opentracing 的 UnaryClientInterceptor
//		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
//		conn, err := grpc.DialContext(ctx, rs.Scheme()+"://authority/"+service, opts...)
//		if err != nil {
//			panic(err)
//		}
//		c.connMap[service] = conn
//	}
//
//	return c.connMap[service]
//}

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
		c.logger.With(zap.String("reqId", reqId))
		c.logger.With(zap.String("duration", time.Since(start).String()))
		c.logger.With(zap.String("method", info.FullMethod))
		c.logger.With(zap.Any("req", req))
		c.logger.With(zap.Any("resp", resp))
		if c.serverName != "" {
			c.logger.With(zap.String("serverName", c.serverName))
		}
		if err != nil {
			c.logger.With(zap.String("err", err.Error()))
		}
		c.logger.Info("server after handling.")
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
		with := c.logger.With(zap.String("reqId", reqId))
		with = c.logger.With(zap.String("duration", time.Since(start).String()))
		with = c.logger.With(zap.String("method", method))
		with = c.logger.With(zap.Any("req", req))

		if c.serverName != "" {
			with = c.logger.With(zap.String("client server:", c.serverName))
		}
		with = c.logger.With(zap.Any("res", reply))
		if err != nil {
			with = c.logger.With(zap.Any("err", err.Error()))
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
