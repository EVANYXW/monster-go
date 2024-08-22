package handler

import (
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"
	"go.uber.org/zap"
	"net"
	"time"
)

type centerConnectorMsgHandler struct {
	isHandle          bool
	owner             module.IModule
	nodes             map[uint32]*network.ServerInfo
	serverInfoHandler module.IServerInfoHandler
}

func NewCenterConnector(serverInfoHandler module.IServerInfoHandler) *centerConnectorMsgHandler {
	return &centerConnectorMsgHandler{
		nodes:             make(map[uint32]*network.ServerInfo),
		serverInfoHandler: serverInfoHandler,
	}
}

func (m *centerConnectorMsgHandler) Start() {

}

func (m *centerConnectorMsgHandler) MsgRegister(processor *network.Processor) {
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_C_Cc_Handshake), m.C_Cc_Handshake)
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_C_Cc_ServerInfo), m.C_Cc_ServerInfo)
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_C_Cc_ServerOk), m.C_Cc_ServerOk)
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_C_Cc_ServerLost), m.C_Cc_ServerLost)
}

//func (m *centerConnectorMsgHandler) GetIsHandle() bool {
//	return m.isHandle
//}

func (m *centerConnectorMsgHandler) OnNetMessage(pack *network.Packet) {
	// isHandle 为 true 消息会来这里处理
}

func (m *centerConnectorMsgHandler) OnNetConnected(np *network.NetPoint) {
	m.SendHandshake()
}

func (m *centerConnectorMsgHandler) OnRpcNetAccept(np *network.NetPoint, acceptor *network.Acceptor) {
	np.Connect()
	conn := np.Conn.(*net.TCPConn)
	acceptor.RemoveConn(conn, np)
}

func (m *centerConnectorMsgHandler) OnNetError(np *network.NetPoint, acceptor *network.Acceptor) {
	m.OnNPDel(np)
}

func (m *centerConnectorMsgHandler) OnServerOk() {

}

func (m *centerConnectorMsgHandler) OnNPAdd(np *network.NetPoint) {
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

func (m *centerConnectorMsgHandler) OnOk() {
	messageID := uint64(xsf_pb.SMSGID_Cc_C_ServerOk)
	msg, _ := rpc.GetMessage(messageID)
	sendMsg := msg.(*xsf_pb.Cc_C_ServerOk)
	sendMsg.ServerId = server.ID
	m.SendMessage(messageID, sendMsg)
}

func (m *centerConnectorMsgHandler) OnNodeOk(id uint32) {
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

func (m *centerConnectorMsgHandler) SendHandshake() {
	messageID := uint64(xsf_pb.SMSGID_Cc_C_Handshake)
	msg, _ := rpc.GetMessage(messageID)
	localMsg := msg.(*xsf_pb.Cc_C_Handshake)
	localMsg.ServerId = server.ID
	localMsg.Ports = server.Ports[:]

	//localMsg := &xsf_pb.Cc_C_Handshake{}
	//localMsg.ServerId = network.ID
	//localMsg.Ports = network.Ports[:]

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
	//pack, _ := servers.ConnectorKernel.Client.Pack(messageID, localMsg)
	//servers.ConnectorKernel.NetPoint.SetSignal(pack)

	m.SendMessage(messageID, localMsg)
}

func (m *centerConnectorMsgHandler) AddNode(message *xsf_pb.C_Cc_ServerInfo) {
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

func (m *centerConnectorMsgHandler) OnNodeLost(id uint32) {
	delete(m.nodes, id)

	// todo
	//c.handler.OnServerLost(id)

	var SID server.ServerID
	server.ID2Sid(id, &SID)
}

func (m *centerConnectorMsgHandler) OnNPDel(np *network.NetPoint) {

	//servers.NodeManager.OnNodeLost(np.ID, np.SID.Type)
}

func (m *centerConnectorMsgHandler) OnHandshakeTicker() {
	async.Go(func() {
		//ticker := time.NewTicker(6 * time.Second)
		//defer ticker.Stop()
		//for range ticker.C {
		//	//msg := &xsf_pb.Cc_C_Heartbeat{}
		//	messageID := uint64(xsf_pb.SMSGID_Cc_C_Heartbeat)
		//	msg, _ := rpc.GetMessage(messageID)
		//	m.SendMessage(messageID, msg)
		//}

		timer := time.NewTimer(6 * time.Second)
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				//msg := &xsf_pb.Cc_C_Heartbeat{}
				messageID := uint64(xsf_pb.SMSGID_Cc_C_Heartbeat)
				msg, _ := rpc.GetMessage(messageID)
				m.SendMessage(messageID, msg)
			}
			_ = timer.Reset(time.Duration(6 * time.Second)) //重制心跳上报时间间隔
		}
	})
}

func (m *centerConnectorMsgHandler) SendMessage(msgId uint64, message interface{}) {
	pack, _ := module.ConnKernel.Client.Pack(message)
	module.ConnKernel.NetPoint.SetSignal(pack)
}

func (m *centerConnectorMsgHandler) C_Cc_Handshake(message *network.Packet) {
	messageID := uint64(xsf_pb.SMSGID_C_Cc_Handshake)
	msg, _ := rpc.GetMessage(messageID)
	localMsg := msg.(*xsf_pb.C_Cc_Handshake)
	rpc.Import(message.Msg.Data, localMsg)
	module.ConnKernel.SetID(localMsg.ServerId)

	server.ID = localMsg.NewId
	server.UpdateSID()
	server.Ports = [server.EP_Max]uint32(localMsg.Ports)

	// 握手定时器
	m.OnHandshakeTicker()
	module.DoWaitStart() // 去开启gate对外的net

	//if xsf_server.Status.Get() == xsf_server.ServerStatus_Running {
	//	xsf_log.Info("服务器本身已启动完毕，直接同步数据")
	//	cc.OnOK(sc)
	//}
	//fixMe gate 链接world ，比world创建tcp链接更快
	if module.Status.Load() == int32(module.ModuleRunStatus_Running) {
		logger.Info("服务器本身已启动完毕，直接同步数据")
		m.OnOk()
	}

}

func (m *centerConnectorMsgHandler) C_Cc_ServerInfo(message *network.Packet) {
	localMsg := &xsf_pb.C_Cc_ServerInfo{}
	rpc.Import(message.Msg.Data, localMsg)
	m.AddNode(localMsg)
	logger.Info("C_Cc_ServerInfo center connector nodes:", zap.Int("node length:", len(m.nodes)))
}

func (m *centerConnectorMsgHandler) C_Cc_ServerOk(message *network.Packet) {
	localMsg := &xsf_pb.C_Cc_ServerOk{}
	rpc.Import(message.Msg.Data, localMsg)
	m.OnNodeOk(localMsg.ServerId)
	logger.Info("C_Cc_ServerOk center connector nodes:", zap.Int("node length:", len(m.nodes)))
}

func (m *centerConnectorMsgHandler) C_Cc_ServerLost(message *network.Packet) {

	localMsg := &xsf_pb.C_Cc_ServerLost{}
	rpc.Import(message.Msg.Data, localMsg)
	m.OnNodeLost(localMsg.ServerId)
	logger.Info("C_Cc_ServerLost center connector nodes:", zap.Int("node length:", len(m.nodes)), zap.Uint32("server id", localMsg.ServerId))
}
