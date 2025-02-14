// Package center @Author evan_yxw
// @Date 2024/9/18 22:48:00
// @Desc
package register_discovery

import (
	"github.com/evanyxw/monster-go/pkg/kernel"
	"github.com/evanyxw/monster-go/pkg/module/connector"
	"github.com/evanyxw/monster-go/pkg/module/module_def"
)

type Connector interface {
	module_def.IModule
}

type Type int

const (
	TypeCenter Type = iota
	TypeEtcd
)

type options struct {
	servername string
	netType    kernel.NetType
	isWatch    bool
}

type Option func(options *options)

type ConnectorFactory interface {
	CreateConnector(options ...Option) Connector
	CreateConnectorManager() connector.TcpConnectorManager
	GetType() Type
	SetGateWay()
}

type NetFactory interface {
	CreateNet() Connector
}

func NewOption() *options {
	return &options{}
}

func (o *options) GetServername() string {
	return o.servername
}

func (o *options) GetIsWatch() bool {
	return o.isWatch
}

func (o *options) GetNetType() kernel.NetType {
	return o.netType
}

func WithServername(name string) Option {
	return func(options *options) {
		options.servername = name
	}
}

func WithNetType(netType kernel.NetType) Option {
	return func(options *options) {
		options.netType = netType
	}
}

func WithWatch() Option {
	return func(options *options) {
		options.isWatch = true
	}
}
