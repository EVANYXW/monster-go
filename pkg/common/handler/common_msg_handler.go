package handler

import (
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
)

type commonMsgHandler struct {
	owner *module.BaseModule
}

func NewCommonMsgHandler() *commonMsgHandler {
	return &commonMsgHandler{}
}

func (m *commonMsgHandler) Init(owner *module.BaseModule) {
	m.owner = owner
}

func (m *commonMsgHandler) OnInit(baseModule module.IBaseModule) {
	baseM := baseModule.(*module.BaseModule)
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

func (m *commonMsgHandler) SendHandshake(ck *module.ConnectorKernel) {

}
