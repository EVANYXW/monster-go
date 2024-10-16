package module

import (
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/server"
)

// 客户端消息接受体
type ClientNet struct {
	kernel       module.IKernel
	curStartNode *configs.ServerNode
	netType      module.NetType
	id           int32
}

func NewClientNet(id int32, maxConnNum uint32, msgHandler module.MsgHandler, netType module.NetType,
	packerFactory network.PackerFactory) *ClientNet {
	c := &ClientNet{
		id:     id,
		kernel: module.NewNetKernel(maxConnNum, msgHandler, packerFactory, module.WithNetType(netType)),
	}
	//module.NewBaseModule(id, c) // todo
	network.NetPointManager = c.kernel.GetNPManager()

	return c
}

func (c *ClientNet) Init(baseModule module.IBaseModule) bool {
	c.kernel.Init(baseModule)
	return true
}

// DoRun BaseModule 调用
func (c *ClientNet) DoRun() {
	c.kernel.DoRun()
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
		return module.ModuleOk()
	}
	return module.ModuleWait()
}

func (c *ClientNet) OnCloseCheck() int {
	return c.kernel.OnCloseCheck()
}

func (c *ClientNet) GetID() int32 {
	return c.id
}

func (c *ClientNet) GetKernel() module.IKernel {
	return c.kernel
}

func (c *ClientNet) Update() {

}

func (c *ClientNet) DoRegister() {
	c.kernel.DoRegister()
}

func (c *ClientNet) OnNetError(np *network.NetPoint) {
	logger.Debug("center onNetError")
	//c.nodeManager.OnNodeLost(np.ID, np.SID.Type)
	module.NodeManager.OnNodeLost(np.ID, np.SID.Type)
}

func (c *ClientNet) OnServerOk() {

}

func (c *ClientNet) OnNPAdd(np *network.NetPoint) {

}
