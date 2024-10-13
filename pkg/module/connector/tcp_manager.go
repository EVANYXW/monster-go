// Package tcp_manager @Author evan_yxw
// @Date 2024/10/12 18:25:00
// @Desc
package connector

import "github.com/evanyxw/monster-go/pkg/module"

type TcpConnectorManager interface {
	CreateConnector(id uint32, ip string, port uint32) *module.ConnectorKernel
	GetConnector(ep uint32, id uint32) module.IKernel
}
