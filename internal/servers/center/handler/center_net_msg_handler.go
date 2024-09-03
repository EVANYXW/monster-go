package handler

import (
	"fmt"
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"
	"go.uber.org/zap"
	"net"
)

type centerNetMsgHandler struct {
	owner module.IModule
}

func NewCenterNetMsg() *centerNetMsgHandler {
	return &centerNetMsgHandler{}
}

func (m *centerNetMsgHandler) Start() {

}

func (m *centerNetMsgHandler) MsgRegister(processor *network.Processor) {
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_Cc_C_Handshake), m.Cc_C_Handshake)
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_Cc_C_Heartbeat), m.Cc_C_Heartbeat)
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_Cc_C_ServerOk), m.Cc_C_ServerOk)
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_Cc_C_ServerClose), m.Cc_C_ServerClose)
}

func (m *centerNetMsgHandler) OnNetMessage(pack *network.Packet) {
	// isHandle 为 true 消息会来这里处理
}

func (m *centerNetMsgHandler) Cc_C_Handshake(message *network.Packet) {
	localMsg := &xsf_pb.Cc_C_Handshake{}
	rpc.Import(message.Msg.Data, localMsg)

	si := module.NodeManager.AddNode(localMsg.ServerId, message.NetPoint.RemoteIP, localMsg.Ports)
	if si == nil {
		fmt.Println("Cc_C_Handshake AddNode si is nil, net point close!")
		message.NetPoint.Close()
		return
	}

	message.NetPoint.SetID(si.ID)

	np := message.NetPoint
	if network.NetPointManager.OnHandshake(np) {
		//c.NetKernel.OnNPAdd(np)
		m.OnNPAdd(np)
		message.NetPoint.OnHeartbeat()
		// 同步本地已经有的服务器列表信息到这个节点
		module.NodeManager.Send(np, si)

		// 再回一个握手消息
		pb := &xsf_pb.C_Cc_Handshake{}
		pb.ServerId = server.ID
		pb.NewId = si.ID
		pb.Ports = si.Ports[:]
		np.SendMessage(pb)

		// 把该节点信息广播给其他所有服务器
		module.NodeManager.Broadcast(si)
	}
}

func (m *centerNetMsgHandler) OnNetConnected(np *network.NetPoint) {

}

func (m *centerNetMsgHandler) OnRpcNetAccept(np *network.NetPoint, acceptor *network.Acceptor) {
	async.Go(func() {
		np.Connect()
	})
	//conn := np.Conn.(*net.TCPConn)
	//acceptor.RemoveConn(conn, np)
}

func (m *centerNetMsgHandler) OnNetError(np *network.NetPoint, acceptor *network.Acceptor) {
	m.OnNPDel(np)
	conn := np.Conn.(*net.TCPConn)
	acceptor.RemoveConn(conn, np)
}

func (m *centerNetMsgHandler) OnServerOk() {

}

func (m *centerNetMsgHandler) OnOk() {

}

func (m *centerNetMsgHandler) OnUpdate() {

}

func (m *centerNetMsgHandler) OnNPAdd(np *network.NetPoint) {
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

func (m *centerNetMsgHandler) OnNPDel(np *network.NetPoint) {
	module.NodeManager.OnNodeLost(np.ID, np.SID.Type)
}

func (m *centerNetMsgHandler) Cc_C_Heartbeat(message *network.Packet) {
	message.NetPoint.OnHeartbeat()
}

func (m *centerNetMsgHandler) Cc_C_ServerOk(message *network.Packet) {
	np := message.NetPoint
	logger.Info("SMSGID_Cc_C_ServerOk", zap.Uint32("id", np.ID))
	module.NodeManager.OnNodeOK(np.ID)
}

func (m *centerNetMsgHandler) Cc_C_ServerClose(message *network.Packet) {
	manager := module.GetManager(module.ModuleID_SM)
	msg := &xsf_pb.C_Cc_ServerClose{}
	manager.Broadcast(msg, 0)
}
