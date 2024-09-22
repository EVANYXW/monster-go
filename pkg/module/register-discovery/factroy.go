// Package center @Author evan_yxw
// @Date 2024/9/18 22:48:00
// @Desc
package register_discovery

import "github.com/evanyxw/monster-go/pkg/module"

type Connector interface {
	module.IModule
}

type ConnectorFactory interface {
	CreateConnector() Connector
	IsConnectorServer() bool
	CreateConnectorManager() Connector
}

type NetFactory interface {
	CreateNet() Connector
}
