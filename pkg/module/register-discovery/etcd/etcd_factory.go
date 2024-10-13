// Package center @Author evan_yxw
// @Date 2024/9/19 12:24:00
// @Desc
package etcd

import (
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/module/connector"
	"github.com/evanyxw/monster-go/pkg/module/connector/factory"
	register_discovery "github.com/evanyxw/monster-go/pkg/module/register-discovery"
)

type Factor struct {
	options
}

type options struct {
	isServerConnector bool
	isGateway         bool
}

type Options func(opt *options)

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
	return false
}

func (f *Factor) SetGateWay() {
	f.isGateway = true
}

func (f *Factor) CreateConnector(servername string, isWatch bool, netType module.NetType) register_discovery.Connector {
	return NewEtcdConnector(module.ModuleID_Etcd, servername, isWatch, netType)
}

func (f *Factor) GetType() register_discovery.Type {
	return register_discovery.TypeEtcd
}

func (f *Factor) CreateConnectorManager() connector.TcpConnectorManager {
	//return connector.NewManager(module.ModuleID_ConnectorManager, managerFactory)
	c := factory.CenterManagerFactory{}
	manager := c.CreateManager(module.ModuleID_ConnectorManager)
	return manager
}

func (f *Factor) CreateNet() register_discovery.Connector {
	return nil
}
