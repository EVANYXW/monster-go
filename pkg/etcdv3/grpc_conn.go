/**
 * @api post etcdv3.
 *
 * User: yunshengzhu
 * Date: 2021/12/19
 * Time: 上午11:08
 */
package etcdv3

import (
	"context"
	"fmt"
	"github.com/evanyxw/monster-go/pkg/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

var connServiceMap = make(map[string]*grpc.ClientConn, 0)
var connServiceMx sync.RWMutex

func NewRpcConn(ctx context.Context, etcd *Etcd, service string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	connServiceMx.RLock()
	if res, ok := connServiceMap[service]; ok {
		connServiceMx.RUnlock()
		return res, nil
	}
	connServiceMx.RUnlock()
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
	opts = append(opts, grpc.WithUnaryInterceptor(middleware.Interceptor))
	rs := NewResolver(etcd, service)
	opts = append(opts, grpc.WithResolvers(rs))
	conn, err := grpc.DialContext(ctx, rs.Scheme()+"://authority/"+service, opts...)
	if err != nil {
		return nil, err
	}
	connServiceMx.Lock()
	connServiceMap[service] = conn
	connServiceMx.Unlock()
	return conn, nil
}
