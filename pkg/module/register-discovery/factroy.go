// Package center @Author evan_yxw
// @Date 2024/9/18 22:48:00
// @Desc
package register_discovery

import (
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/module/connector"
)

type Connector interface {
	module.IModule
}

type Type int

const (
	TypeCenter Type = iota
	TypeEtcd
)

type ConnectorFactory interface {
	CreateConnector(servername string) Connector
	//IsConnectorServer() bool
	CreateConnectorManager(managerFactory connector.ManagerFactory) Connector
	GetType() Type
	SetGateWay()
}

type NetFactory interface {
	CreateNet() Connector
}
