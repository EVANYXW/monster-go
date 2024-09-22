package module

import (
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/internal/servers/center/handler"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
)

type CenterConnector struct {
	kernel module.IModuleKernel
	id     int32
}

func NewCenterConnector(id int32, serverInfoHandler module.IServerInfoHandler) *CenterConnector {
	centerCnf := configs.Get().Center
	c := &CenterConnector{
		id: id,
		kernel: module.NewConnectorKernel(centerCnf.Ip, centerCnf.Port,
			handler.NewCenterConnectorMsg(serverInfoHandler),
			new(network.DefaultPackerFactory),
			module.WithCNoWaitStart(true)),
	}

	module.ConnKernel = c.kernel.(*module.ConnectorKernel)
	module.NewBaseModule(id, c)

	return c
}

func (c *CenterConnector) Init(baseModule *module.BaseModule) bool {
	c.kernel.Init(baseModule)
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

func (c *CenterConnector) OnOk() {
	c.kernel.OnOk()
}

func (c *CenterConnector) OnStartCheck() int {
	return c.kernel.OnStartCheck()
}

func (c *CenterConnector) OnCloseCheck() int {
	return c.kernel.OnCloseCheck()
}

func (c *CenterConnector) GetID() int32 {
	return c.id
}

func (c *CenterConnector) GetKernel() module.IModuleKernel {
	return c.kernel
}

func (c *CenterConnector) Update() {

}

func (c *CenterConnector) DoRegister() {
	c.kernel.DoRegister()
}
