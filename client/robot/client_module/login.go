// Package module @Author evan_yxw
// @Date 2024/8/19 13:56:00
// @Desc
package client_module

import "C"
import (
	"github.com/evanyxw/monster-go/client/robot/handler"
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/kernel"
	"github.com/evanyxw/monster-go/pkg/module/module_def"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
)

var _ module_def.IModule = &Login{}

type Login struct {
	*kernel.ConnectorKernel
	id int32
}

func (l *Login) GetID() int32 {
	return l.id
}

func (l *Login) GetKernel() module_def.IKernel {
	return l.ConnectorKernel
}

func (l *Login) Init(baseModule module_def.IBaseModule) bool {
	return true
}

func (l *Login) DoRegister() {
	l.ConnectorKernel.DoRegister()
}

func (l *Login) DoRun() {
	l.ConnectorKernel.DoRun()

	l.Handshake()
	l.Login()
}

func (l *Login) DoWaitStart() {

}

func (l *Login) DoRelease() {

}

func (l *Login) OnOk() {

}

func (l *Login) Update() {

}

func New(id int32) *Login {
	l := &Login{
		id: id,
		ConnectorKernel: kernel.NewConnectorKernel("192.168.101.8", 30000,
			handler.NewLoginHandler(),
			new(network.DefaultPackerFactory),
			kernel.WithCNoWaitStart(true),
		),
	}
	module_def.NewBaseModule(id, l)
	return l
}

func (l *Login) Handshake() {
	msg := &xsf_pb.Clt_Gt_Handshake{}
	l.SendMessage(msg)
}

func (l *Login) Login() {
	messageID := uint64(xsf_pb.MSGID_Clt_L_Login)
	msg, _ := rpc.GetMessage(messageID)
	localMsg := msg.(*xsf_pb.Clt_L_Login)
	localMsg.LoginType = uint32(xsf_pb.LoginType_PHXH)
	localMsg.LoginDatas = []string{
		"yxw",
		"123456",
	}
	l.SendMessage(localMsg)
}
