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
	//TODO implement me
	panic("implement me")
}

func (f *Factor) CreateConnector(servername string) register_discovery.Connector {
	if f.isGateway {
		return NewCenterConnector(module.ModuleID_CenterConnector, handler.NewGateServerInfoHandler())
	}
	return NewCenterConnector(module.ModuleID_CenterConnector, handler.NewServerInfoHandler())
}

func (f *Factor) CreateConnectorManager(managerFactory connector.ManagerFactory) register_discovery.Connector {
	return connector.NewManager(module.ModuleID_ConnectorManager, managerFactory)
}

func (f *Factor) SetGateWay() {
	f.isGateway = true
}

func (f *Factor) CreateNet() register_discovery.Connector {
	return NewCenterNet(module.ModuleID_SM, 10000)
}

func (f *Factor) GetType() register_discovery.Type {
	return register_discovery.TypeCenter
}

type options struct {
	isConnectorManager bool
	isGateway          bool
}

type Options func(opt *options)
