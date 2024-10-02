// Package center @Author evan_yxw
// @Date 2024/9/19 12:24:00
// @Desc
package center

import (
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/module/connector"
	register_discovery "github.com/evanyxw/monster-go/pkg/module/register-discovery"
	"github.com/evanyxw/monster-go/pkg/module/register-discovery/center/handler"
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

func (f *Factor) CreateConnectorManager(managerFactory connector.ManagerFactory) register_discovery.Connector {
	return connector.NewManager(module.ModuleID_ConnectorManager, managerFactory)
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
