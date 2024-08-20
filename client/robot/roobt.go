package robot

import (
	"github.com/evanyxw/monster-go/client/robot/client_module"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
)

type Robot struct {
	serverInfoHandler module.IServerInfoHandler
	ID                int32

	loginModule *client_module.Login
}

func NewRobot() *Robot {

	robot := &Robot{
		loginModule: client_module.New(module.ModuleID_Client),
	}

	return robot
}

func (r *Robot) Start() {

	//r.loginModule.Start()
}

//func (r *Robot) Handshake() {
//	msg := &xsf_pb.Clt_Gt_Handshake{}
//	pack, _ := r.connectorKernel.Client.Pack(msg)
//	r.connectorKernel.NetPoint.SetSignal(pack)
//}
//
//func (r *Robot) Login() {
//	messageID := uint64(xsf_pb.MSGID_Clt_L_Login)
//	msg, _ := rpc.GetMessage(messageID)
//	localMsg := msg.(*xsf_pb.Clt_L_Login)
//
//	pack, _ := r.connectorKernel.Client.Pack(localMsg)
//
//	r.connectorKernel.NetPoint.SetSignal(pack)
//	fmt.Println("Test")
//}

//func (r *Robot) Test() {
//	messageID := uint64(xsf_pb.SMSGID_Cc_C_Handshake)
//	msg, _ := rpc.GetMessage(messageID)
//	localMsg := msg.(*xsf_pb.Cc_C_Handshake)
//	localMsg.ServerId = server.ID
//	localMsg.Ports = server.Ports[:]
//
//	pack, _ := r.connectorKernel.Client.Pack(localMsg)
//
//	r.connectorKernel.NetPoint.SetSignal(pack)
//	fmt.Println("Test")
//}

func (r *Robot) OnNetError(np *network.NetPoint) {

}

func (r *Robot) OnServerOk() {

}

func (r *Robot) OnNPAdd(np *network.NetPoint) {

}
