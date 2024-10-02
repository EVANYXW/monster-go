// Package connector @Author evan_yxw
// @Date 2024/9/25 23:33:00
// @Desc
package connector

import (
	"github.com/evanyxw/monster-go/pkg/module"
)

type IManager interface {
	CreateConnector(id uint32, ip string, port uint32) *module.ConnectorKernel
	GetConnector(ep uint32, id uint32) module.IModuleKernel
}

type ManagerFactory interface {
	Create(id int32) *Manager
	CreateConnector(handler module.GateAcceptorHandler, id uint32, ip string, port uint32) *module.ConnectorKernel
}
