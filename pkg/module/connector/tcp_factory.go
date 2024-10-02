// Package connector @Author evan_yxw
// @Date 2024/9/28 14:12:00
// @Desc
package connector

import (
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/module/connector/handler"
	"github.com/evanyxw/monster-go/pkg/network"
)

// TcpManagerFactory tcp 管理器
type TcpManagerFactory struct {
}

func (t *TcpManagerFactory) CreateConnector(handler module.GateAcceptorHandler, id uint32, ip string, port uint32) *module.ConnectorKernel {
	//msgHandler := handler.NewManager()
	msgHandler := handler.(module.MsgHandler)
	ck := module.NewConnectorKernel(ip, port, msgHandler, new(network.ClientPackerFactory))
	//ck := module.NewConnectorKernel(ip, port, msgHandler, new(network.DefaultPackerFactory))
	ck.SetID(id)
	return ck
}

func (t *TcpManagerFactory) Create(id int32) *Manager {
	c := &Manager{
		id:      id,
		factory: t,
	}
	hdler := handler.NewManagerMsg()
	c.handler = hdler
	c.kernel = module.NewKernel(hdler, network.NetPointManager.GetRpcAcceptor(),
		network.NetPointManager.GetProcessor())

	return c
}
