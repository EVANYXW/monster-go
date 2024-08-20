package handler

import (
	"fmt"
	"github.com/evanyxw/monster-go/internal/servers/gate/acceptor"
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/network"
)

type loginMsgHandler struct {
	acceptor acceptor.IAcceptor
}

func NewLoginMsgHandler() *loginMsgHandler {
	return &loginMsgHandler{
		acceptor: acceptor.NewAcceptor(),
	}
}

func (m *loginMsgHandler) Start() {

}

func (m *loginMsgHandler) OnNetMessage(pack *network.Packet) {

}

func (m *loginMsgHandler) OnNetError(np *network.NetPoint, acceptor *network.Acceptor) {

}

func (m *loginMsgHandler) OnNetConnected(np *network.NetPoint) {

}

func (m *loginMsgHandler) OnRpcNetAccept(np *network.NetPoint, acceptor *network.Acceptor) {

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

	pb := &xsf_pb.L_Clt_LoginResult{}
	pb.Result = uint32(xsf_pb.LoginResult_LoginParamError)
	m.acceptor.SendMessage2Client(message, pb)
}
