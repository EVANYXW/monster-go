// Package center @Author evan_yxw
// @Date 2024/9/19 12:24:00
// @Desc
package etcd

import (
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/module/connector"
	register_discovery "github.com/evanyxw/monster-go/pkg/module/register-discovery"
)

type Factor struct {
	options
}

type options struct {
	isServerConnector bool
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

func (f *Factor) CreateConnector(servername string) register_discovery.Connector {
	return NewEtcdConnector(module.ModuleID_Etcd, servername)
}

func (f *Factor) GetType() register_discovery.Type {
	return register_discovery.TypeEtcd
}

func (f *Factor) CreateConnectorManager(managerFactory connector.ManagerFactory) register_discovery.Connector {
	return nil
}

func (f *Factor) CreateNet() register_discovery.Connector {
	return nil
}
