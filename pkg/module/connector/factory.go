// Package connector @Author evan_yxw
// @Date 2024/9/25 23:33:00
// @Desc
package connector

import (
	"github.com/evanyxw/monster-go/pkg/server/tcp_manager"
)

//type ManagerFactory interface {
//	Create(id int32) *Manager
//	CreateConnector(handler module.GateAcceptorHandler, id uint32, ip string, port uint32) *module.ConnectorKernel
//}

type AbstractFactory interface {
	CreateManager(id int32) tcp_manager.TcpConnectorManager
}
