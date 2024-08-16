package module

import (
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/internal/servers"
	"github.com/evanyxw/monster-go/internal/servers/center/handler"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
)

type CenterConnector struct {
	*module.BaseModule
	connectorKernel *module.ConnectorKernel
	ID              int32
}

func NewCenterConnector(id int32, serverInfoHandler module.IServerInfoHandler) *CenterConnector {
	c := &CenterConnector{
		ID: id,
	}

	centerCnf := configs.Get().Center
	c.connectorKernel = module.NewConnectorKernel(centerCnf.Ip, centerCnf.Port,
		handler.NewCenterConnector(serverInfoHandler),
		network.NewDefaultPacker(),
		module.WithCNoWaitStart(true))

	c.BaseModule = module.NewBaseModule(c)
	servers.ConnectorKernel = c.connectorKernel

	return c
}

func (c *CenterConnector) Init() {
	c.connectorKernel.Init()
}

func (c *CenterConnector) DoRun() {
	//c.DoRegister()
	c.connectorKernel.Start()
	//c.OnHandshake() // handler
}

func (c *CenterConnector) DoWaitStart() {

}

func (c *CenterConnector) DoRelease() {
	c.connectorKernel.Release()
}

func (c *CenterConnector) GetID() int32 {
	return c.ID
}

func (c *CenterConnector) OnOk() {
	c.connectorKernel.OnOK()
}

func (c *CenterConnector) OnStartCheck() int {
	return c.connectorKernel.OnStartCheck()
}

func (c *CenterConnector) OnCloseCheck() int {
	return c.connectorKernel.OnCloseCheck()
}

func (c *CenterConnector) GetKernel() module.IModuleKernel {
	return c.connectorKernel
}

func (c *CenterConnector) Update() {

}

func (c *CenterConnector) DoRegister() {
	c.connectorKernel.DoRegist()
}
