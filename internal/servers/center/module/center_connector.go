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
	kernel module.IModuleKernel
	ID     int32
}

func NewCenterConnector(id int32, serverInfoHandler module.IServerInfoHandler) *CenterConnector {
	c := &CenterConnector{
		ID: id,
	}

	centerCnf := configs.Get().Center
	c.kernel = module.NewConnectorKernel(centerCnf.Ip, centerCnf.Port,
		handler.NewCenterConnector(serverInfoHandler),
		new(network.DefaultPackerFactory),
		module.WithCNoWaitStart(true))

	c.BaseModule = module.NewBaseModule(c)
	servers.ConnectorKernel = c.kernel.(*module.ConnectorKernel)

	return c
}

func (c *CenterConnector) Init() bool {
	c.kernel.Init()
	return true
}

func (c *CenterConnector) DoRun() {
	c.kernel.DoRun()
}

func (c *CenterConnector) DoWaitStart() {

}

func (c *CenterConnector) DoRelease() {
	c.kernel.DoRelease()
}

func (c *CenterConnector) GetID() int32 {
	return c.ID
}

func (c *CenterConnector) OnOk() {
	c.kernel.OnOk()
}

func (c *CenterConnector) OnStartCheck() int {
	return c.kernel.OnStartCheck()
}

func (c *CenterConnector) OnCloseCheck() int {
	return c.kernel.OnCloseCheck()
}

func (c *CenterConnector) GetKernel() module.IModuleKernel {
	return c.kernel
}

func (c *CenterConnector) Update() {

}

func (c *CenterConnector) DoRegister() {
	c.kernel.DoRegister()
}
