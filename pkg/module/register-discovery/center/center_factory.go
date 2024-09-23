// Package center @Author evan_yxw
// @Date 2024/9/19 12:24:00
// @Desc
package center

import (
	"github.com/evanyxw/monster-go/pkg/module"
	register_discovery "github.com/evanyxw/monster-go/pkg/module/register-discovery"
	"github.com/evanyxw/monster-go/pkg/module/register-discovery/center/handler"
	"github.com/evanyxw/monster-go/pkg/module/register-discovery/center/manager"
)

type Factor struct {
	options
}

func NewFactor(opts ...Options) *Factor {
	opt := &options{}
	for _, fun := range opts {
		fun(opt)
	}

	return &Factor{
		options: *opt,
	}
}

func (f *Factor) IsConnectorServer() bool {
	return f.options.isConnectorServer
}

func (f *Factor) CreateConnector() register_discovery.Connector {
	if f.isConnectorServer {
		return NewCenterConnector(module.ModuleID_CenterConnector, handler.NewGateServerInfoHandler())
	}
	return NewCenterConnector(module.ModuleID_CenterConnector, handler.NewServerInfoHandler())
}

func (f *Factor) CreateConnectorManager() register_discovery.Connector {
	return manager.NewConnectorManager(module.ModuleID_ConnectorManager)
}

func (f *Factor) CreateNet() register_discovery.Connector {
	return NewCenterNet(module.ModuleID_SM, 10000)
}

type options struct {
	isConnectorServer bool
}

type Options func(opt *options)

func WithServerConnectorManager() Options {
	return func(opt *options) {
		opt.isConnectorServer = true
	}
}
