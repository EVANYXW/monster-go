// Package tcp_manager @Author evan_yxw
// @Date 2024/10/12 18:25:00
// @Desc
package connector

import (
	"github.com/evanyxw/monster-go/pkg/module/module_def"
	"github.com/evanyxw/monster-go/pkg/network"
)

type TcpConnectorManager interface {
	//CreateConnector(id uint32, ip string, port uint32) *kernel.ConnectorKernel
	CreateConnector(id uint32, ip string, port uint32) network.IConn
	GetConnector(ep uint32, id uint32) module_def.IKernel
	DelConnector(id uint32)
}
