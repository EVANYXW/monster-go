package client

import (
	"context"
	"fmt"
	"github.com/evanyxw/game_proto/msg"
	"github.com/evanyxw/monster-go/internal/rpc"
	c "github.com/evanyxw/monster-go/pkg/context"
	"github.com/evanyxw/monster-go/pkg/etcdv3"
	"github.com/evanyxw/monster-go/pkg/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

var merchantQueryConn *grpc.ClientConn
var merchantQueryConnOnce sync.Once

func NewLogicRpcClient() (msg.WorldClient, error) {
	merchantQueryConnOnce.Do(func() {
		ctx := c.NewContext(context.Background())
		opts := make([]grpc.DialOption, 0)
		rs := etcdv3.NewResolver(etcdConfig, rpc.LogicRpc)
		opts = append(opts, grpc.WithResolvers(rs))
		opts = append(opts, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
		opts = append(opts, grpc.WithUnaryInterceptor(middleware.Interceptor))
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		conn, err := grpc.DialContext(ctx, rs.Scheme()+"://authority/"+rpc.LogicRpc, opts...)
		if err != nil {
			panic(err)
		}
		if err != nil {
			ctx.Log.Error("err:", err)
			panic(err)
			return
		}
		merchantQueryConn = conn
	})
	return msg.NewWorldClient(merchantQueryConn), nil
}
