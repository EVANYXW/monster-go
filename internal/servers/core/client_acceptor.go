package core

import (
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/server"
)

type ClientNet struct {
	*module.BaseModule
	*module.NetKernel
	module.NodeManager
	curStartNode *configs.ServerNode
	ID           int32
	status       int
	startIndex   int

	netType module.NetType
}

func NewClientNet(id int32, maxConnNum uint32, info server.Info, netType module.NetType) *ClientNet {
	acceptor := &ClientNet{
		ID:          id,
		NodeManager: module.NewNodeManager(),
	}
	acceptor.NetKernel = module.NewNetKernel(maxConnNum, info, acceptor, false, netType)
	baseModule := module.NewBaseModule(acceptor)
	baseModule.Init()
	acceptor.BaseModule = baseModule
	return acceptor
}

// 外部通知开启Module
func (c *ClientNet) Run() {
	c.BaseModule.Run()
}

func (c *ClientNet) Init() {
	c.NetKernel.Init()
}

func (c *ClientNet) DoRun() {
	c.DoRegister()
	c.NodeManager.Start()
	c.NetKernel.Start()

	c.status = server.CN_RunStep_StartServer
	c.startIndex = 0
}

func (c *ClientNet) DoStart() {
	c.NetKernel.DoStart()
}

func (c *ClientNet) DoRelease() {
	c.NetKernel.Release()
}

func (c *ClientNet) OnStartCheck() int {
	return module.ModuleRunCode_Ok
}

func (c *ClientNet) OnCloseCheck() int {
	return c.NetKernel.OnCloseCheck()
}

func (c *ClientNet) Update() {

}

func (c *ClientNet) GetID() int32 {
	return c.ID
}

func (c *ClientNet) DoRegister() {
	c.NetKernel.DoRegist()
	//c.NetKernel.RegisterMsg(uint16(xsf_pb.SMSGID_Cc_C_Handshake), c.Cc_C_Handshake)
	//c.NetKernel.RegisterMsg(uint16(xsf_pb.SMSGID_Cc_C_Heartbeat), c.Cc_C_Heartbeat)
	//c.NetKernel.RegisterMsg(uint16(xsf_pb.SMSGID_Cc_C_ServerOk), c.Cc_C_ServerOk)
}

func (c *ClientNet) Release() {
	c.NetKernel.Release()
}

func (c *ClientNet) OnNetError(np *network.NetPoint) {
	logger.Debug("center onNetError")
	c.NodeManager.OnNodeLost(np.ID, np.SID.Type)
}

func (c *ClientNet) OnServerOk() {

}

func (c *ClientNet) OnNPAdd(np *network.NetPoint) {

}
