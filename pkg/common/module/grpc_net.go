package module

import (
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/server"
)

// GrpcNet Grpc 网络
type GrpcNet struct {
	kernel       module.IModuleKernel
	curStartNode *configs.ServerNode
	netType      module.NetType
	id           int32
}

func NewGrpcNet(id int32, serverName string, msgHandler module.MsgHandler) *ClientNet {
	c := &ClientNet{
		id:      id,
		netType: module.Inner,
		kernel:  module.NewGrpcNetKernel(serverName, msgHandler),
	}

	//network.NetPointManager = c.kernel.GetNPManager()
	return c
}

func (c *GrpcNet) Init(baseModule *module.BaseModule) bool {
	c.kernel.Init(baseModule)
	return true
}

// DoRun BaseModule 调用
func (c *GrpcNet) DoRun() {
	c.kernel.DoRun()
}

func (c *GrpcNet) DoWaitStart() {
	c.kernel.DoWaitStart()
}

func (c *GrpcNet) DoRelease() {
	c.kernel.DoRelease()
}

func (c *GrpcNet) OnOk() {

}

func (c *GrpcNet) OnStartCheck() int {
	// TCP链接准备好
	if c.kernel.GetStatus() == server.Net_RunStep_Done {
		return module.ModuleRunCode_Ok
	}
	return module.ModuleRunCode_Wait
}

func (c *GrpcNet) OnCloseCheck() int {
	return c.kernel.OnCloseCheck()
}

func (c *GrpcNet) GetID() int32 {
	return c.id
}

func (c *GrpcNet) GetKernel() module.IModuleKernel {
	return c.kernel
}

func (c *GrpcNet) Update() {

}

func (c *GrpcNet) DoRegister() {
	c.kernel.DoRegister()
}

func (c *GrpcNet) OnNetError(np *network.NetPoint) {
	logger.Debug("center onNetError")
	//c.nodeManager.OnNodeLost(np.ID, np.SID.Type)
	module.NodeManager.OnNodeLost(np.ID, np.SID.Type)
}

func (c *GrpcNet) OnServerOk() {

}

func (c *GrpcNet) OnNPAdd(np *network.NetPoint) {

}
