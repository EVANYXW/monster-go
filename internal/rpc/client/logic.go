package client

import (
	"bilibili/monster-go/internal/rpc"
	"context"
	"fmt"
	"github.com/evanyxw/game_proto/msg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	c "hl.hexinchain.com/welfare-center/basic/context"
	"hl.hexinchain.com/welfare-center/basic/etcdv3"
	"hl.hexinchain.com/welfare-center/basic/middleware"
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
