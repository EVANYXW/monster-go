package handler

import (
	"fmt"
	"github.com/evanyxw/monster-go/internal/servers/gate/acceptor"
	"github.com/evanyxw/monster-go/internal/servers/login/client"
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
)

type loginMsgHandler struct {
	gateAcceptor  acceptor.IAcceptor
	clientManager *client.Manager
	owner         *module.BaseModule
}

func NewLoginMsgHandler() *loginMsgHandler {
	clientManager := client.NewClientManager()
	module.GtAClientManager = clientManager
	return &loginMsgHandler{
		gateAcceptor:  acceptor.NewGate(),
		clientManager: clientManager,
	}
}

func (m *loginMsgHandler) OnInit(baseModule module.IBaseModule) {
	baseM := baseModule.(*module.BaseModule)
	m.owner = baseM
	m.clientManager.Init(m.owner.RpcAcceptor)
}

func (m *loginMsgHandler) Start() {

}

func (m *loginMsgHandler) OnNetMessage(pack *network.Packet) {

}

func (m *loginMsgHandler) MsgRegister(processor *network.Processor) {
	processor.RegisterMsg(uint16(xsf_pb.MSGID_Clt_L_Login), m.Clt_L_Login)
}

func (m *loginMsgHandler) OnServerOk() {

}

func (m *loginMsgHandler) OnOk() {

}

func (m *loginMsgHandler) OnUpdate() {
	m.clientManager.OnUpdate()
}

func (m *loginMsgHandler) OnNPAdd(np *network.NetPoint) {

}

func (m *loginMsgHandler) Clt_L_Login(message *network.Packet) {
	fmt.Println("我收到登录消息啦～")
	msg := &xsf_pb.Clt_L_Login{}
	rpc.Import(message.Msg.Data, msg)

	isOK := true
	switch msg.LoginType {
	case uint32(xsf_pb.LoginType_PHXH):
		isOK = len(msg.LoginDatas) == int(xsf_pb.PHXHLoginData_PHXHLD_Max)
	default:
		isOK = false
	}

	if !isOK {
		sendMsg := &xsf_pb.L_Clt_LoginResult{}
		sendMsg.Result = uint32(xsf_pb.LoginResult_LoginParamError)
		m.gateAcceptor.SendMessage2Client(message, sendMsg)

		m.gateAcceptor.DisconnectClient(message, uint32(xsf_pb.DisconnectReason_LoginError))
		return
	}

	clt := m.clientManager.NewClient(message.Msg.RawID, msg.LoginType, msg.LoginDatas)
	if clt != nil {
		clt.Start()
	}
}
