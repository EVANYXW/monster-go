package module

import (
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/internal/servers"
	"github.com/evanyxw/monster-go/internal/servers/center/handler"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/server"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"strings"
)

type CenterNet struct {
	*module.BaseModule
	netKernel   *module.NetKernel
	nodeManager module.NodeManager

	ID           int32
	status       int
	startIndex   int
	curStartNode *configs.ServerNode
}

func NewCenterNet(id int32, maxConnNum uint32, info server.Info) *CenterNet {
	centerNet := &CenterNet{
		ID:          id,
		nodeManager: module.NewNodeManager(),
		netKernel: module.NewNetKernel(maxConnNum, info, handler.NewNetCenter(false),
			module.WithNoWaitStart(true)),
	}

	centerNet.BaseModule = module.NewBaseModule(centerNet)

	servers.NetPointManager = centerNet.netKernel.GetNPManager()
	servers.NodeManager = centerNet.nodeManager

	return centerNet
}

func (c *CenterNet) Init() {
	c.netKernel.Init()
}

func (c *CenterNet) DoRun() {
	c.DoRegister()
	c.nodeManager.Start()
	c.netKernel.Start()

	c.status = server.CN_RunStep_StartServer
	c.startIndex = 0
}

func (c *CenterNet) DoStart() {

}

func (c *CenterNet) DoRelease() {
	c.netKernel.Release()
}

func (c *CenterNet) OnStartCheck() int {
	serverCnf := configs.Get().Server
	if !serverCnf.AutoStart {
		return module.ModuleRunCode_Ok
	}

	serverList := configs.Get().ServerList
	switch c.status {
	case server.CN_RunStep_StartServer:
		c.curStartNode = &(serverList[c.startIndex])
		dir, _ := os.Getwd()

		// 兼容开发时的直接运行
		binDir := dir + "/bin"
		cmdStr := "./bin/nld_server run --server_name " + c.curStartNode.EPName

		_, err := os.Stat(binDir)
		if os.IsNotExist(err) {
			logger.Info("找不到bin文件夹，执行当前目录sh文件")
			cmdStr = "./single_start.sh " + c.curStartNode.EPName
		}

		cmdFields := strings.Fields(cmdStr)
		cmd := exec.Command(cmdFields[0], cmdFields[1:]...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Start()
		if err != nil {
			logger.Error("Error running command:", zap.Error(err))
		}

		// fixMe 恢复
		c.status = server.CN_RunStep_WaitHandshake
		//c.status = server.CN_RunStep_HandshakeDone
	case server.CN_RunStep_HandshakeDone:
		c.startIndex++
		if c.startIndex >= len(serverList) {
			return module.ModuleRunCode_Ok
		} else {
			c.status = server.CN_RunStep_StartServer
		}
	}

	return module.ModuleRunCode_Wait
}

func (c *CenterNet) OnCloseCheck() int {
	return c.netKernel.OnCloseCheck()
}

func (c *CenterNet) Update() {

}

func (c *CenterNet) GetID() int32 {
	return c.ID
}

func (c *CenterNet) DoRegister() {
	c.netKernel.DoRegist()
}

func (c *CenterNet) Release() {
	c.netKernel.Release()
}

func (c *CenterNet) GetKernel() module.IModuleKernel {
	return c.netKernel
}

func (c *CenterNet) OnServerOk() {

}

//func (c *CenterNet) OnNetError(np *network.NetPoint) {
//	c.OnNPDel(np)
//}
//
//func (c *CenterNet) OnNPAdd(np *network.NetPoint) {
//	//if c.curStartNode == nil {
//	//	return
//	//}
//	//
//	//// fixMe 恢复
//	////if np.SID.Type == network.Name2EP(c.curStartNode.EPName) {
//	////	logger.Info("centerNetHandler OnNPAdd", zap.Uint16("server", np.SID.ID),
//	////		zap.String("type", network.EP2Name(np.SID.Type)), zap.Uint8("index", np.SID.Index))
//	////	c.status = server.CN_RunStep_HandshakeDone
//	////}
//}
//
//func (c *CenterNet) OnNPDel(np *network.NetPoint) {
//	c.nodeManager.OnNodeLost(np.ID, np.SID.Type)
//}

func (c *CenterNet) GetNodeManager() module.NodeManager {
	return c.nodeManager
}

func (c *CenterNet) GetNetKernel() *module.NetKernel {
	return c.netKernel
}

//func (c *CenterNet) Cc_C_Handshake(message *network.Packet) {
//	localMsg := &xsf_pb.Cc_C_Handshake{}
//	proto.Unmarshal(message.Msg.Data, localMsg)
//	si := c.nodeManager.AddNode(localMsg.ServerId, message.NetPoint.RemoteIP, localMsg.Ports)
//	if si == nil {
//		message.NetPoint.Close()
//		return
//	}
//
//	message.NetPoint.SetID(si.ID)
//
//	np := message.NetPoint
//	if c.netKernel.GetNPManager().OnHandshake(np) {
//		//c.NetKernel.OnNPAdd(np)
//		c.OnNPAdd(np)
//		message.NetPoint.OnHeartbeat()
//		// 同步本地已经有的服务器列表信息到这个节点
//		c.nodeManager.Send(np, si)
//
//		// 再回一个握手消息
//		pb := &xsf_pb.C_Cc_Handshake{}
//		pb.ServerId = server.ID
//		pb.NewId = si.ID
//		pb.Ports = si.Ports[:]
//		np.SendMessage(uint64(xsf_pb.SMSGID_C_Cc_Handshake), pb)
//
//		// 把该节点信息广播给其他所有服务器
//		c.nodeManager.Broadcast(si)
//	}
//}
//
//func (c *CenterNet) Cc_C_Heartbeat(message *network.Packet) {
//	message.NetPoint.OnHeartbeat()
//}
//
//func (c *CenterNet) Cc_C_ServerOk(message *network.Packet) {
//	np := message.NetPoint
//	logger.Info("SMSGID_Cc_C_ServerOk", zap.Uint32("id", np.ID))
//	c.nodeManager.OnNodeOK(np.ID)
//}
