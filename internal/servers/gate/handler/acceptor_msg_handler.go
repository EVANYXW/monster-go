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

/***
其他服务器被gate连接 就要用该acceptor
*/

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

func (m *AcceptorMsgHandler) OnRpcNetAccept(np *network.NetPoint, acceptor *network.Acceptor) {
	async.Go(func() {
		np.Connect()
	})
}

func (m *AcceptorMsgHandler) OnNetError(np *network.NetPoint, acceptor *network.Acceptor) {
	conn := np.Conn.(*net.TCPConn)
	acceptor.RemoveConn(conn, np)
}

func (m *AcceptorMsgHandler) OnServerOk() {

}

func (m *AcceptorMsgHandler) OnOk() {

}

func (m *AcceptorMsgHandler) OnNPAdd(np *network.NetPoint) {

}

func (m *AcceptorMsgHandler) OnUpdate() {

}

func (m *AcceptorMsgHandler) MsgRegister(processor *network.Processor) {
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_Gt_GtA_Handshake), m.Gt_GtA_Handshake)
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_Gt_GtA_Heartbeat), m.Gt_GtA_Heartbeat)
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_Gt_GtA_ClientClose), m.Gt_GtA_ClientClose)

}

func (m *AcceptorMsgHandler) Gt_GtA_Handshake(message *network.Packet) {
	localMsg := &xsf_pb.Gt_GtA_Handshake{}
	rpc.Import(message.Msg.Data, localMsg)
	message.NetPoint.SetID(localMsg.ServerId)

	logger.Info("acceptor OnNetMessage", zap.Uint32("ServerID", localMsg.ServerId))

	// 回一个握手消息
	if network.NetPointManager.OnHandshake(message.NetPoint) {
		logger.Info("acceptor OnNetMessage OnHandshake done", zap.Uint32("ServerID", localMsg.ServerId))
		msgId := uint64(xsf_pb.SMSGID_GtA_Gt_Handshake)
		msg, _ := rpc.GetMessage(msgId)
		localMsg := msg.(*xsf_pb.GtA_Gt_Handshake)

		localMsg.ServerId = server.ID
		message.NetPoint.SendMessage(localMsg)
	}
}

func (m *AcceptorMsgHandler) Gt_GtA_Heartbeat(message *network.Packet) {
	message.NetPoint.OnHeartbeat()
}

func (m *AcceptorMsgHandler) Gt_GtA_ClientClose(message *network.Packet) {
	fmt.Println("我收到gate的消息，关闭client")
	fmt.Println(message.Msg.RawID)
	module.GtAClientManager.CloseClient(message.Msg.RawID)
}
