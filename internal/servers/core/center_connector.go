package core

import (
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"time"
)

type CenterConnector struct {
	*module.BaseModule
	*module.ConnectorKernel
	serverInfoHandler module.IServerInfoHandler
	nodes             map[uint32]*network.ServerInfo
	ID                int32
}

func NewCenterConnector(id int32, serverInfoHandler module.IServerInfoHandler) *CenterConnector {
	c := &CenterConnector{
		ID: id,

		nodes:             make(map[uint32]*network.ServerInfo),
		serverInfoHandler: serverInfoHandler,
	}
	c.ConnectorKernel = module.NewConnectorKernel(c, "", 8023)

	baseModule := module.NewBaseModule(c)
	baseModule.NoWaitStart = true
	baseModule.Init()

	c.BaseModule = baseModule

	c.ConnectorKernel.Owner = c
	return c
}

func (c *CenterConnector) Init() {
	c.ConnectorKernel.Init()
}

func (c *CenterConnector) DoRun() {
	c.DoRegister()
	c.ConnectorKernel.Start()
	c.OnHandshake()
}

func (c *CenterConnector) DoStart() {

}

func (c *CenterConnector) DoRelease() {
	c.ConnectorKernel.Release()
}

func (c *CenterConnector) GetID() int32 {
	return c.ID
}

func (c *CenterConnector) OnStartCheck() int {
	return c.ConnectorKernel.OnStartCheck()
}

func (c *CenterConnector) OnCloseCheck() int {
	return c.ConnectorKernel.OnCloseCheck()
}

func (c *CenterConnector) OnServerOk() {

}

func (c *CenterConnector) Update() {

}

func (c *CenterConnector) OnNetError(np *network.NetPoint) {

}

func (c *CenterConnector) OnOK() {
	messageID := uint64(xsf_pb.SMSGID_Cc_C_ServerOk)
	msg, _ := rpc.GetMessage(messageID)
	sendMsg := msg.(*xsf_pb.Cc_C_ServerOk)
	pack, _ := c.Client.Pack(messageID, sendMsg)
	c.ConnectorKernel.NetPoint.SetSignal(pack)
}

func (c *CenterConnector) OnNPAdd(np *network.NetPoint) {

}

func (c *CenterConnector) OnNodeOk(id uint32) {
	var SID network.ServerID
	network.ID2Sid(id, &SID)

	s, ok := c.nodes[id]
	if !ok {
		logger.Error("centerConnectorHandler OnCCServerOk server not found", zap.Uint32("id", id), zap.Uint16("server", SID.ID),
			zap.String("type", network.EP2Name(SID.Type)),
			zap.Uint8("server", SID.Index))
	} else {
		s.Status = network.ServerInfo_Ok

		logger.Info("【中心服连接器】收到服务器已准备好", zap.Uint32("id", s.ID), zap.Uint16("server", SID.ID),
			zap.String("type", network.EP2Name(SID.Type)),
			zap.Uint8("index", SID.Index))

		c.serverInfoHandler.OnServerOk(s)

		isAllOK := true
		for _, node := range c.nodes {
			if node.Status != network.ServerInfo_Ok {
				isAllOK = false
			}
		}

		//xsf_log.Info("服务器节点信息", xsf_log.Bool("is all ok", isAllOK), xsf_log.Int("cc.nodes", len(cc.nodes)), xsf_log.Int("node list", len(xsf_config.NodeList)))
		// todo len(xsf_config.NodeList)
		if isAllOK && len(c.nodes)+1 >= 6 {
			logger.Info("服务器所有节点已全部开启！！！")
			// todo
			c.serverInfoHandler.OnServerOpenComplete()
		}
	}
}

func (c *CenterConnector) DoRegister() {
	c.ConnectorKernel.DoRegist()
	c.ConnectorKernel.RegisterMsg(uint16(xsf_pb.SMSGID_C_Cc_Handshake), c.C_Cc_Handshake)
	c.ConnectorKernel.RegisterMsg(uint16(xsf_pb.SMSGID_C_Cc_ServerInfo), c.C_Cc_ServerInfo)
	c.ConnectorKernel.RegisterMsg(uint16(xsf_pb.SMSGID_C_Cc_ServerOk), c.C_Cc_ServerOk)
	c.ConnectorKernel.RegisterMsg(uint16(xsf_pb.SMSGID_C_Cc_ServerLost), c.C_Cc_ServerLost)
}

func (c *CenterConnector) OnHandshake() {
	messageID := uint64(xsf_pb.SMSGID_Cc_C_Handshake)
	msg, _ := rpc.GetMessage(messageID)
	localMsg := msg.(*xsf_pb.Cc_C_Handshake)
	localMsg.ServerId = network.ID
	localMsg.Ports = network.Ports[:]

	//localMsg := &xsf_pb.Cc_C_Handshake{}
	//localMsg.ServerId = network.ID
	//localMsg.Ports = network.Ports[:]

	pack, _ := c.Client.Pack(messageID, localMsg)

	//message, _ := c.Client.UnPack(pack)
	//fmt.Println("unpack message id:", message.ID)
	//fmt.Println("unpack message data:", message.Data)

	//buffer := network.BufferPacker{}
	//byteData := buffer.TestPack(messageID, localMsg)
	//fmt.Println(byteData)
	//message, _ := buffer.TestRead(byteData)
	//var is xsf_pb.Cc_C_Handshake
	//err := proto.Unmarshal(message.Data, &is)
	//fmt.Println(is, err)

	c.ConnectorKernel.NetPoint.SetSignal(pack)
}

func (c *CenterConnector) AddNode(message *xsf_pb.C_Cc_ServerInfo) {
	pb := message
	for _, info := range pb.Infos {
		node, ok := c.nodes[info.ServerId]

		isNewAdd := false
		if ok {
			node.IP = info.Ip
			node.Ports = [network.EP_Max]uint32(info.Ports)
			node.Status = info.Status
		} else {
			isNewAdd = true
			node = new(network.ServerInfo)
			node.ID = info.ServerId
			node.IP = info.Ip
			node.Ports = [network.EP_Max]uint32(info.Ports)
			node.Status = info.Status
			c.nodes[node.ID] = node
		}

		var SID network.ServerID
		network.ID2Sid(node.ID, &SID)

		//logger.Info("收到服务器信息", zapcore.Field{Key: "id", Integer: int64(node.ID)},
		//	zapcore.Field{Key: "type", Integer: int64(SID.Type)},
		//	zapcore.Field{Key: "index", Integer: int64(SID.Index)},
		//	zapcore.Field{Key: "status", Integer: int64(node.Status)})
		logger.Info("收到服务器信息")
		if isNewAdd {
			logger.Info("新增结点")
			c.serverInfoHandler.OnServerNew(node)
			if node.Status == network.ServerInfo_Ok {
				c.OnNodeOk(node.ID)
			}
		}
	}
}

func (c *CenterConnector) OnNodeLost(id uint32) {
	delete(c.nodes, id)

	// todo
	//c.handler.OnServerLost(id)

	var SID network.ServerID
	network.ID2Sid(id, &SID)
}

func (c *CenterConnector) handshakeTicker() {
	async.Go(func() {
		ticker := time.NewTicker(6 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			//msg := &xsf_pb.Cc_C_Heartbeat{}
			messageID := uint64(xsf_pb.SMSGID_Cc_C_Heartbeat)
			msg, _ := rpc.GetMessage(messageID)
			pack, _ := c.Client.Pack(messageID, msg)
			c.NetPoint.SetSignal(pack)
		}
	})
}

func (c *CenterConnector) C_Cc_Handshake(message *network.Packet) {
	messageID := uint64(xsf_pb.SMSGID_C_Cc_Handshake)
	msg, _ := rpc.GetMessage(messageID)
	localMsg := msg.(*xsf_pb.C_Cc_Handshake)
	proto.Unmarshal(message.Msg.Data, localMsg)
	c.SetID(localMsg.ServerId)

	network.ID = localMsg.NewId
	network.UpdateSID()
	network.Ports = [network.EP_Max]uint32(localMsg.Ports)

	// 握手定时器
	c.handshakeTicker()
	module.DoStart() // 去开启gate对外的net

	//if xsf_server.Status.Get() == xsf_server.ServerStatus_Running {
	//	xsf_log.Info("服务器本身已启动完毕，直接同步数据")
	//	cc.OnOK(sc)
	//}
	c.OnOK()
}

func (c *CenterConnector) C_Cc_ServerInfo(message *network.Packet) {
	localMsg := &xsf_pb.C_Cc_ServerInfo{}
	proto.Unmarshal(message.Msg.Data, localMsg)
	c.AddNode(localMsg)
	logger.Info("C_Cc_ServerInfo center connector nodes:", zap.Int("node length:", len(c.nodes)))
}

func (c *CenterConnector) C_Cc_ServerOk(message *network.Packet) {
	localMsg := &xsf_pb.C_Cc_ServerOk{}
	proto.Unmarshal(message.Msg.Data, localMsg)
	c.OnNodeOk(localMsg.ServerId)
	logger.Info("C_Cc_ServerOk center connector nodes:", zap.Int("node length:", len(c.nodes)))
}

func (c *CenterConnector) C_Cc_ServerLost(message *network.Packet) {

	localMsg := &xsf_pb.C_Cc_ServerLost{}
	proto.Unmarshal(message.Msg.Data, localMsg)
	c.OnNodeLost(localMsg.ServerId)
	logger.Info("C_Cc_ServerLost center connector nodes:", zap.Int("node length:", len(c.nodes)))
}
