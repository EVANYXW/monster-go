package handler

import (
	"github.com/evanyxw/monster-go/internal/servers"
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"time"
)

type MsgHandler struct {
	isHandle          bool
	owner             module.IModule
	nodes             map[uint32]*network.ServerInfo
	serverInfoHandler module.IServerInfoHandler
}

func NewConnectorCenter(isHandle bool, serverInfoHandler module.IServerInfoHandler) *MsgHandler {
	return &MsgHandler{
		isHandle:          isHandle,
		nodes:             make(map[uint32]*network.ServerInfo),
		serverInfoHandler: serverInfoHandler,
	}
}

func (m *MsgHandler) Start() {
	m.OnHandshake()
}

func (m *MsgHandler) GetIsHandle() bool {
	return m.isHandle
}

func (m *MsgHandler) MsgRegister(processor *network.Processor) {
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_C_Cc_Handshake), m.C_Cc_Handshake)
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_C_Cc_ServerInfo), m.C_Cc_ServerInfo)
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_C_Cc_ServerOk), m.C_Cc_ServerOk)
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_C_Cc_ServerLost), m.C_Cc_ServerLost)
}

func (m *MsgHandler) HandleMsg(pack *network.Packet) {
	// isHandle 为 true 消息会来这里处理
}

func (m *MsgHandler) OnNetError(np *network.NetPoint) {
	m.OnNPDel(np)
}

func (m *MsgHandler) OnServerOk() {

}

func (m *MsgHandler) OnNPAdd(np *network.NetPoint) {
	//if m.curStartNode == nil {
	//	return
	//}

	// fixMe 恢复
	//if np.SID.Type == network.Name2EP(c.curStartNode.EPName) {
	//	logger.Info("centerNetHandler OnNPAdd", zap.Uint16("server", np.SID.ID),
	//		zap.String("type", network.EP2Name(np.SID.Type)), zap.Uint8("index", np.SID.Index))
	//	c.status = server.CN_RunStep_HandshakeDone
	//}
}

func (m *MsgHandler) OnOK() {
	messageID := uint64(xsf_pb.SMSGID_Cc_C_ServerOk)
	msg, _ := rpc.GetMessage(messageID)
	sendMsg := msg.(*xsf_pb.Cc_C_ServerOk)
	pack, _ := servers.ConnectorKernel.Client.Pack(messageID, sendMsg)
	servers.ConnectorKernel.NetPoint.SetSignal(pack)
}

func (m *MsgHandler) OnNodeOk(id uint32) {
	var SID server.ServerID
	server.ID2Sid(id, &SID)

	s, ok := m.nodes[id]
	if !ok {
		logger.Error("centerConnectorHandler OnCCServerOk server not found", zap.Uint32("id", id), zap.Uint16("server", SID.ID),
			zap.String("type", server.EP2Name(SID.Type)),
			zap.Uint8("server", SID.Index))
	} else {
		s.Status = network.ServerInfo_Ok

		logger.Info("【中心服连接器】收到服务器已准备好", zap.Uint32("id", s.ID), zap.Uint16("server", SID.ID),
			zap.String("type", server.EP2Name(SID.Type)),
			zap.Uint8("index", SID.Index))

		m.serverInfoHandler.OnServerOk(s)

		isAllOK := true
		for _, node := range m.nodes {
			if node.Status != network.ServerInfo_Ok {
				isAllOK = false
			}
		}

		//xsf_log.Info("服务器节点信息", xsf_log.Bool("is all ok", isAllOK), xsf_log.Int("cc.nodes", len(cc.nodes)), xsf_log.Int("node list", len(xsf_config.NodeList)))
		// todo len(xsf_config.NodeList)
		if isAllOK && len(m.nodes)+1 >= 6 {
			logger.Info("服务器所有节点已全部开启！！！")
			// todo
			m.serverInfoHandler.OnServerOpenComplete()
		}
	}
}

func (m *MsgHandler) OnHandshake() {
	messageID := uint64(xsf_pb.SMSGID_Cc_C_Handshake)
	msg, _ := rpc.GetMessage(messageID)
	localMsg := msg.(*xsf_pb.Cc_C_Handshake)
	localMsg.ServerId = server.ID
	localMsg.Ports = server.Ports[:]

	//localMsg := &xsf_pb.Cc_C_Handshake{}
	//localMsg.ServerId = network.ID
	//localMsg.Ports = network.Ports[:]

	pack, _ := servers.ConnectorKernel.Client.Pack(messageID, localMsg)

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

	servers.ConnectorKernel.NetPoint.SetSignal(pack)
}

func (m *MsgHandler) AddNode(message *xsf_pb.C_Cc_ServerInfo) {
	pb := message
	for _, info := range pb.Infos {
		node, ok := m.nodes[info.ServerId]

		isNewAdd := false
		if ok {
			node.IP = info.Ip
			node.Ports = [server.EP_Max]uint32(info.Ports)
			node.Status = info.Status
		} else {
			isNewAdd = true
			node = new(network.ServerInfo)
			node.ID = info.ServerId
			node.IP = info.Ip
			node.Ports = [server.EP_Max]uint32(info.Ports)
			node.Status = info.Status
			m.nodes[node.ID] = node
		}

		var SID server.ServerID
		server.ID2Sid(node.ID, &SID)

		//logger.Info("收到服务器信息", zapcore.Field{Key: "id", Integer: int64(node.ID)},
		//	zapcore.Field{Key: "type", Integer: int64(SID.Type)},
		//	zapcore.Field{Key: "index", Integer: int64(SID.Index)},
		//	zapcore.Field{Key: "status", Integer: int64(node.Status)})
		logger.Info("收到服务器信息")
		if isNewAdd {
			logger.Info("新增结点")
			m.serverInfoHandler.OnServerNew(node)
			if node.Status == network.ServerInfo_Ok {
				m.OnNodeOk(node.ID)
			}
		}
	}
}

func (m *MsgHandler) OnNodeLost(id uint32) {
	delete(m.nodes, id)

	// todo
	//c.handler.OnServerLost(id)

	var SID server.ServerID
	server.ID2Sid(id, &SID)
}

func (m *MsgHandler) OnNPDel(np *network.NetPoint) {
	servers.NodeManager.OnNodeLost(np.ID, np.SID.Type)
}

func (m *MsgHandler) handshakeTicker() {
	async.Go(func() {
		ticker := time.NewTicker(6 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			//msg := &xsf_pb.Cc_C_Heartbeat{}
			messageID := uint64(xsf_pb.SMSGID_Cc_C_Heartbeat)
			msg, _ := rpc.GetMessage(messageID)
			pack, _ := servers.ConnectorKernel.Client.Pack(messageID, msg)
			servers.ConnectorKernel.NetPoint.SetSignal(pack)
		}
	})
}

func (m *MsgHandler) C_Cc_Handshake(message *network.Packet) {
	messageID := uint64(xsf_pb.SMSGID_C_Cc_Handshake)
	msg, _ := rpc.GetMessage(messageID)
	localMsg := msg.(*xsf_pb.C_Cc_Handshake)
	proto.Unmarshal(message.Msg.Data, localMsg)
	servers.ConnectorKernel.SetID(localMsg.ServerId)

	server.ID = localMsg.NewId
	server.UpdateSID()
	server.Ports = [server.EP_Max]uint32(localMsg.Ports)

	// 握手定时器
	m.handshakeTicker()
	module.DoStart() // 去开启gate对外的net

	//if xsf_server.Status.Get() == xsf_server.ServerStatus_Running {
	//	xsf_log.Info("服务器本身已启动完毕，直接同步数据")
	//	cc.OnOK(sc)
	//}
	//fixMe gate 链接world ，比world创建tcp链接更快
	m.OnOK()
}

func (m *MsgHandler) C_Cc_ServerInfo(message *network.Packet) {
	localMsg := &xsf_pb.C_Cc_ServerInfo{}
	proto.Unmarshal(message.Msg.Data, localMsg)
	m.AddNode(localMsg)
	logger.Info("C_Cc_ServerInfo center connector nodes:", zap.Int("node length:", len(m.nodes)))
}

func (m *MsgHandler) C_Cc_ServerOk(message *network.Packet) {
	localMsg := &xsf_pb.C_Cc_ServerOk{}
	proto.Unmarshal(message.Msg.Data, localMsg)
	m.OnNodeOk(localMsg.ServerId)
	logger.Info("C_Cc_ServerOk center connector nodes:", zap.Int("node length:", len(m.nodes)))
}

func (m *MsgHandler) C_Cc_ServerLost(message *network.Packet) {

	localMsg := &xsf_pb.C_Cc_ServerLost{}
	proto.Unmarshal(message.Msg.Data, localMsg)
	m.OnNodeLost(localMsg.ServerId)
	logger.Info("C_Cc_ServerLost center connector nodes:", zap.Int("node length:", len(m.nodes)))
}
