package handler

import (
	"github.com/evanyxw/monster-go/internal/servers"
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
)

/***
其他服务器连接 gate 就要用该acceptor
*/

//type IAcceptor interface {
//	//AddCloseEventRpc(rpc IClientCloseHandler)
//	//RegistClientMessage(msgID uint16, rpc *xsf_rpc.Acceptor)
//
//	// 断开一个客户端连接
//	DisconnectClient(id uint32, reason uint32)
//
//	// 发送消息到网关
//	//SendMessage2Agent(agentID uint32, msg xsf_net.IMessage)
//
//	// 发送消息到客户端
//	//SendMessage2Client(clientID uint32, msg xsf_net.IMessage)
//
//	// 广播一个消息给所有客户端
//	//Broadcast(msg xsf_net.IMessage)
//
//	// 设置一个客户端的服务器转播id
//	SetServerID(clientID uint32, ep uint8, serverID uint32)
//}
//
//type acceptor struct {
//}
//
//var (
//	acc *acceptor
//)
//
//func GetAcceptor() IAcceptor {
//	return &acceptor{}
//}

type AcceptorMsgHandler struct {
}

func NewAcceptor() *AcceptorMsgHandler {
	// 设置GlobalProcess后，net_kernel的Process将被GlobalProcess接管
	//network.GlobalProcess = network.NewProcessor()
	return &AcceptorMsgHandler{}
}

func (m *AcceptorMsgHandler) Start() {

}

//func (m *AcceptorMsgHandler) GetIsHandle() bool {
//	return m.isHandle
//}

func (m *AcceptorMsgHandler) OnNetMessage(pack *network.Packet) {

}

func (m *AcceptorMsgHandler) OnNetConnected(np *network.NetPoint) {

}

func (m *AcceptorMsgHandler) OnRpcNetAccept(np *network.NetPoint) {
	np.Connect()
}

func (m *AcceptorMsgHandler) OnNetError(np *network.NetPoint) {

}

func (m *AcceptorMsgHandler) OnServerOk() {

}

func (m *AcceptorMsgHandler) OnOk() {

}

func (m *AcceptorMsgHandler) OnNPAdd(np *network.NetPoint) {

}

func (m *AcceptorMsgHandler) MsgRegister(processor *network.Processor) {
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_Gt_GtA_Handshake), m.Gt_GtA_Handshake)
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_Gt_GtA_Heartbeat), m.Gt_GtA_Heartbeat)
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_Gt_GtA_ClientClose), m.Gt_GtA_ClientClose)
}

func (m *AcceptorMsgHandler) Gt_GtA_Handshake(message *network.Packet) {
	localMsg := &xsf_pb.Gt_GtA_Handshake{}
	proto.Unmarshal(message.Msg.Data, localMsg)
	message.NetPoint.SetID(localMsg.ServerId)

	logger.Info("acceptor OnNetMessage", zap.Uint32("ServerID", localMsg.ServerId))

	// 回一个握手消息
	if servers.NetPointManager.OnHandshake(message.NetPoint) {
		logger.Info("acceptor OnNetMessage OnHandshake done", zap.Uint32("ServerID", localMsg.ServerId))
		msgId := uint64(xsf_pb.SMSGID_GtA_Gt_Handshake)
		msg, _ := rpc.GetMessage(msgId)
		localMsg := msg.(*xsf_pb.GtA_Gt_Handshake)

		localMsg.ServerId = server.ID
		message.NetPoint.SendMessage(msgId, localMsg)
	}
}

func (m *AcceptorMsgHandler) Gt_GtA_Heartbeat(message *network.Packet) {
	message.NetPoint.OnHeartbeat()
}

func (m *AcceptorMsgHandler) Gt_GtA_ClientClose(message *network.Packet) {

}
