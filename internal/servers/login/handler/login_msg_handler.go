package handler

import (
	"fmt"
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
)

type loginMsgHandler struct {
	isHandle bool
}

func NewLoginMsgHandler(isHandle bool) *loginMsgHandler {
	return &loginMsgHandler{
		isHandle: isHandle,
	}
}

func (m *loginMsgHandler) Start() {

}

//func (m *loginMsgHandler) GetIsHandle() bool {
//	return m.isHandle
//}

func (m *loginMsgHandler) OnNetMessage(pack *network.Packet) {

}

func (m *loginMsgHandler) OnNetError(np *network.NetPoint) {

}

func (m *loginMsgHandler) OnNetConnected(np *network.NetPoint) {

}

func (m *loginMsgHandler) OnRpcNetAccept(np *network.NetPoint) {

}

func (m *loginMsgHandler) OnServerOk() {

}

func (m *loginMsgHandler) OnOk() {

}

func (m *loginMsgHandler) OnNPAdd(np *network.NetPoint) {

}

func (m *loginMsgHandler) MsgRegister(processor *network.Processor) {
	processor.RegisterMsg(uint16(xsf_pb.MSGID_Clt_L_Login), m.Clt_L_Login)
}

func (m *loginMsgHandler) Clt_L_Login(message *network.Packet) {
	fmt.Println("我收到登录消息啦～")

	msg, _ := rpc.GetMessage(uint64(xsf_pb.SMSGID_GtA_Gt_ClientMessage))
	localMsg := msg.(*xsf_pb.GtA_Gt_ClientMessage)
	localMsg.ClientId = append(localMsg.ClientId, message.NetPoint.ID)
	message.NetPoint.SendMessage(localMsg)
}
