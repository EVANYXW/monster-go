// Package center @Author evan_yxw
// @Date 2024/9/19 12:24:00
// @Desc
package center

import (
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

	return &Factor{}
}

func (c *Factor) CreateConnector() register_discovery.Connector {
	if c.isConnectorServer {
		return NewCenterConnector(handler.NewGateServerInfoHandler())
	}
	return NewCenterConnector(handler.NewServerInfoHandler())
}

func (c *Factor) CreateNet() register_discovery.Connector {
	return NewCenterNet(10000)
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
