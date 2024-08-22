package module

import (
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/internal/servers"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/server"
)

// 客户端消息接受体

type ClientNet struct {
	*module.BaseModule
	kernel       module.IModuleKernel
	nodeManager  module.NodeManager
	curStartNode *configs.ServerNode

	ID         int32
	startIndex int

	netType module.NetType
}

func NewClientNet(id int32, maxConnNum uint32, msgHandler module.MsgHandler, info server.Info, netType module.NetType,
	packerFactory network.PackerFactory) *ClientNet {
	c := &ClientNet{
		ID:          id,
		nodeManager: module.NewNodeManager(),
	}
	c.kernel = module.NewNetKernel(maxConnNum, info, msgHandler, packerFactory, module.WithNetType(netType))
	baseModule := module.NewBaseModule(c)

	c.BaseModule = baseModule
	//servers.NetPointManager = c.kernel.NPManager
	servers.NetPointManager = c.kernel.GetNPManager()

	return c
}

func (c *ClientNet) Init() bool {
	c.kernel.Init()
	return true
}

// DoRun BaseModule 调用
func (c *ClientNet) DoRun() {
	//c.DoRegister()
	c.nodeManager.Start()
	c.kernel.DoRun()

	c.startIndex = 0
}

func (c *ClientNet) DoWaitStart() {
	c.kernel.DoWaitStart()
}

func (c *ClientNet) DoRelease() {
	c.kernel.DoRelease()
}

func (c *ClientNet) OnOk() {

}

func (c *ClientNet) OnStartCheck() int {
	// TCP链接准备好
	if c.kernel.GetStatus() == server.Net_RunStep_Done {
		return module.ModuleRunCode_Ok
	}
	return module.ModuleRunCode_Wait
}

func (c *ClientNet) OnCloseCheck() int {
	return c.kernel.OnCloseCheck()
}

func (c *ClientNet) GetKernel() module.IModuleKernel {
	return c.kernel
}

func (c *ClientNet) Update() {

}

func (c *ClientNet) GetID() int32 {
	return c.ID
}

func (c *ClientNet) DoRegister() {
	c.kernel.DoRegister()
}

func (c *ClientNet) OnNetError(np *network.NetPoint) {
	logger.Debug("center onNetError")
	c.nodeManager.OnNodeLost(np.ID, np.SID.Type)
}

func (c *ClientNet) OnServerOk() {

}

func (c *ClientNet) OnNPAdd(np *network.NetPoint) {

}
