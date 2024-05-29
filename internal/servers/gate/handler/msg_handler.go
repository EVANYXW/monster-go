package handler

import (
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"
	"go.uber.org/zap"
)

type MsgHandler struct {
	isHandle bool
}

func New(isHandle bool) *MsgHandler {
	return &MsgHandler{
		isHandle: isHandle,
	}
}

func (m *MsgHandler) Start() {

}

func (m *MsgHandler) GetIsHandle() bool {
	return m.isHandle
}

func (m *MsgHandler) HandleMsg(pack *network.Packet) {
	ep := rpc.GetMsgEp(pack.Msg.ID)
	switch ep {
	case server.EP_Game: // 发往游戏服
		fallthrough
	case server.EP_Login: // 发往登录服
		fallthrough
	case server.EP_Mail: // 发往邮件服
		fallthrough
	case server.EP_Manager:
	default:
		logger.Error("ClientKernel OnNetData epDest error", zap.Int("epDest", ep))
		//c.Disconnect(int32(xsf_pb.DisconnectReason_MsgInvalid), true)
	}
}

func (m *MsgHandler) OnNetError(np *network.NetPoint) {

}

func (m *MsgHandler) OnServerOk() {

}

func (m *MsgHandler) OnNPAdd(np *network.NetPoint) {

}

func (m *MsgHandler) MsgRegister(processor *network.Processor) {
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_Cc_C_Handshake), m.Cc_C_Handshake)
}

func (m *MsgHandler) Cc_C_Handshake(message *network.Packet) {

}
