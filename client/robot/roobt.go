package robot

import (
	"fmt"
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"
)

type Robot struct {
	*module.BaseModule
	connectorKernel   *module.ConnectorKernel
	serverInfoHandler module.IServerInfoHandler
	ID                int32
}

func NewRobot() *Robot {
	robot := &Robot{}
	robot.connectorKernel = module.NewConnectorKernel(robot, "", 30000)
	return robot
}

func (r *Robot) Start() {
	r.connectorKernel.Start()
	r.Test()
}

func (r *Robot) Test() {
	messageID := uint64(xsf_pb.SMSGID_Cc_C_Handshake)
	msg, _ := rpc.GetMessage(messageID)
	localMsg := msg.(*xsf_pb.Cc_C_Handshake)
	localMsg.ServerId = server.ID
	localMsg.Ports = server.Ports[:]

	pack, _ := r.connectorKernel.Client.Pack(messageID, localMsg)

	r.connectorKernel.NetPoint.SetSignal(pack)
	fmt.Println("Test")
}

func (r *Robot) OnNetError(np *network.NetPoint) {

}

func (r *Robot) OnServerOk() {

}

func (r *Robot) OnNPAdd(np *network.NetPoint) {

}
