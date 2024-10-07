// Package grpc @Author evan_yxw
// @Date 2024/10/5 11:55:00
// @Desc
package grpc

import (
	"context"
	"github.com/evanyxw/monster-go/pkg/plugins"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"time"
)

const (
	_name = "grpc"
)

var (
	_factory = &Factory{}
)

type Factory struct {
}

func NewFactory() *Factory {
	_factory = &Factory{}
	return _factory
}

func (f Factory) Type() plugins.Type {
	return plugins.Grpc
}

func (f Factory) Name() string {
	return _name
}

func (f Factory) Setup(ctx context.Context) (plugins.PluginIns, error) {
	//TODO implement me
	panic("implement me")
}

func (f Factory) Destroy(ins plugins.PluginIns) error {
	//TODO implement me
	panic("implement me")
}

func (f Factory) Reload(ins plugins.PluginIns, m map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (f Factory) CanUnload(ins plugins.PluginIns) bool {
	//TODO implement me
	panic("implement me")
}

func (f Factory) ClientConn(target string, opts ...grpc.DialOption) (grpc.ClientConnInterface, error) { // 获取Client对象
	//ttarget := fmt.Sprintf("router://%s", target)
	opts = append(opts,
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"Balance_m3g"}`),
		//grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(transport.Instance().ClientInterceptors()...)),
		grpc.WithTimeout(time.Second*10),
	)
	if conn, err := grpc.Dial(target, opts...); err != nil {
		return nil, errors.Wrapf(err, "Dial Target %s", target)
	} else {
		return conn, err
	}
}
