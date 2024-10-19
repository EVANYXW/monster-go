package handler

import (
	"github.com/evanyxw/monster-go/pkg/module/module_def"
	"github.com/evanyxw/monster-go/pkg/network"
)

type commonMsgHandler struct {
	owner *module_def.BaseModule
}

func NewCommonMsgHandler() *commonMsgHandler {
	return &commonMsgHandler{}
}

func (m *commonMsgHandler) Init(owner *module_def.BaseModule) {
	m.owner = owner
}

func (m *commonMsgHandler) OnInit(baseModule module_def.IBaseModule) {
	baseM := baseModule.(*module_def.BaseModule)
	m.owner = baseM
}

func (m *commonMsgHandler) Start() {

}

func (m *commonMsgHandler) OnNetMessage(pack *network.Packet) {

}

func (m *commonMsgHandler) MsgRegister(processor *network.Processor) {

}

func (m *commonMsgHandler) OnServerOk() {

}

func (m *commonMsgHandler) OnOk() {

}

func (m *commonMsgHandler) OnUpdate() {

}

func (m *commonMsgHandler) OnNPAdd(np *network.NetPoint) {

}

func (m *commonMsgHandler) SendHandshake(iconn network.IConn) {

}
