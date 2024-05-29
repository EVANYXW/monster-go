package handler

import (
	"github.com/evanyxw/monster-go/internal/servers"
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/server"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
)

type MsgNetHandler struct {
	isHandle bool
	owner    module.IModule
}

func NewNetCenter(isHandle bool) *MsgNetHandler {

	return &MsgNetHandler{
		isHandle: isHandle,
	}
}

func (m *MsgNetHandler) Start() {

}

func (m *MsgNetHandler) GetIsHandle() bool {
	return m.isHandle
}

func (m *MsgNetHandler) MsgRegister(processor *network.Processor) {
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_Cc_C_Handshake), m.Cc_C_Handshake)
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_Cc_C_Heartbeat), m.Cc_C_Heartbeat)
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_Cc_C_ServerOk), m.Cc_C_ServerOk)
}

func (m *MsgNetHandler) HandleMsg(pack *network.Packet) {
	// isHandle 为 true 消息会来这里处理
}

func (m *MsgNetHandler) Cc_C_Handshake(message *network.Packet) {
	localMsg := &xsf_pb.Cc_C_Handshake{}
	proto.Unmarshal(message.Msg.Data, localMsg)

	si := servers.NodeManager.AddNode(localMsg.ServerId, message.NetPoint.RemoteIP, localMsg.Ports)
	if si == nil {
		message.NetPoint.Close()
		return
	}

	message.NetPoint.SetID(si.ID)

	np := message.NetPoint
	if servers.NetPointManager.OnHandshake(np) {
		//c.NetKernel.OnNPAdd(np)
		m.OnNPAdd(np)
		message.NetPoint.OnHeartbeat()
		// 同步本地已经有的服务器列表信息到这个节点
		servers.NodeManager.Send(np, si)

		// 再回一个握手消息
		pb := &xsf_pb.C_Cc_Handshake{}
		pb.ServerId = server.ID
		pb.NewId = si.ID
		pb.Ports = si.Ports[:]
		np.SendMessage(uint64(xsf_pb.SMSGID_C_Cc_Handshake), pb)

		// 把该节点信息广播给其他所有服务器
		servers.NodeManager.Broadcast(si)
	}
}

func (m *MsgNetHandler) OnNetError(np *network.NetPoint) {
	m.OnNPDel(np)
}

func (m *MsgNetHandler) OnServerOk() {

}

func (m *MsgNetHandler) OnNPAdd(np *network.NetPoint) {
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

func (m *MsgNetHandler) OnNPDel(np *network.NetPoint) {
	servers.NodeManager.OnNodeLost(np.ID, np.SID.Type)
}

func (m *MsgNetHandler) Cc_C_Heartbeat(message *network.Packet) {
	message.NetPoint.OnHeartbeat()
}

func (m *MsgNetHandler) Cc_C_ServerOk(message *network.Packet) {
	np := message.NetPoint
	logger.Info("SMSGID_Cc_C_ServerOk", zap.Uint32("id", np.ID))
	servers.NodeManager.OnNodeOK(np.ID)
}
