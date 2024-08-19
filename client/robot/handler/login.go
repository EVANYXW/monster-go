// Package handler @Author evan_yxw
// @Date 2024/8/19 14:13:00
// @Desc
package handler

import (
	"fmt"
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"
	"github.com/golang/protobuf/proto"
	"time"
)

type loginMsgHandler struct {
}

func NewLoginHandler() *loginMsgHandler {
	return &loginMsgHandler{}
}

func (m *loginMsgHandler) Start() {

}

func (m *loginMsgHandler) OnNetMessage(pack *network.Packet) {

}

func (m *loginMsgHandler) OnNetConnected(np *network.NetPoint) {

}

func (m *loginMsgHandler) OnRpcNetAccept(np *network.NetPoint) {
	np.Connect()
}

func (m *loginMsgHandler) OnNetError(np *network.NetPoint) {

}

func (m *loginMsgHandler) OnServerOk() {

}

func (m *loginMsgHandler) OnOk() {

}

func (m *loginMsgHandler) OnNPAdd(np *network.NetPoint) {

}

func (m *loginMsgHandler) MsgRegister(processor *network.Processor) {
	processor.RegisterMsg(uint16(xsf_pb.MSGID_Gt_Clt_Handshake), m.Gt_Clt_Handshake)
	processor.RegisterMsg(uint16(xsf_pb.MSGID_L_Clt_LoginResult), m.L_Clt_LoginResult)
}

func (m *loginMsgHandler) SendHandshake(ck *module.ConnectorKernel) {
	messageID := uint64(xsf_pb.SMSGID_Gt_GtA_Handshake)
	msg, _ := rpc.GetMessage(messageID)
	localMsg := msg.(*xsf_pb.Gt_GtA_Handshake)
	localMsg.ServerId = server.ID
	//pack, _ := servers.ConnectorKernel.Client.Pack(messageID, localMsg)
	//servers.ConnectorKernel.NetPoint.SetSignal(pack)

	ck.NetPoint.SendMessage(localMsg)
}

func (m *loginMsgHandler) OnHandshakeTicker(np *network.NetPoint) {
	async.Go(func() {
		cnf := configs.Get()
		ticker := time.NewTicker(time.Second * time.Duration(cnf.HtCheck))
		defer ticker.Stop()

		for range ticker.C {
			np.SendMessage(&xsf_pb.Gt_GtA_Heartbeat{})
		}
	})
}

func (m *loginMsgHandler) SendMessage(msgId uint64, message proto.Message) {

}

func (m *loginMsgHandler) Gt_Clt_Handshake(message *network.Packet) {
	fmt.Println("GtA_Gt_Handshake laile")
	logger.Info("GtA_Gt_Handshake 收到")
	//m.OnHandshakeTicker(message.NetPoint)
}

func (m *loginMsgHandler) L_Clt_LoginResult(message *network.Packet) {
	resMsg := &xsf_pb.L_Clt_LoginResult{}
	proto.Unmarshal(message.Msg.Data, resMsg)

	fmt.Println(resMsg)
}
