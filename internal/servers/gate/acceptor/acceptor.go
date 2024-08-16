// Package acceptor @Author evan_yxw
// @Date 2024/8/6 16:51:00
// @Desc
package acceptor

import (
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"
	"github.com/golang/protobuf/proto"
)

type acceptor struct {
	ID int64
}

func NewAcceptor() *acceptor {
	return &acceptor{}
}

func (m *acceptor) SendMessage2Client(clientID uint32) {
	message, _ := rpc.GetMessage(uint64(xsf_pb.SMSGID_GtA_Gt_ClientMessage))
	localMsg := message.(*xsf_pb.GtA_Gt_ClientMessage)

	localMsg.ClientId = make([]uint32, 1)
	localMsg.ClientId[0] = clientID
	// fixMe
	//localMsg.ClientMessage = xsf_net.GetClientPacker().Pack(msg)

	var CID server.ClientID
	server.ID2Cid(clientID, &CID)

	m.SendMessage2Agent(uint32(CID.Gate), uint64(xsf_pb.SMSGID_GtA_Gt_ClientMessage), localMsg)

}

func (m *acceptor) SendMessage2Agent(agentID uint32, msgId uint64, message proto.Message) {
	manager := module.GetManager(module.ModuleID_SM)
	np := manager.Get(agentID)
	if np != nil {
		np.SendMessage(msgId, message)
	}
	//manager := GetManager(int(a.id))
	//np := manager.Get(agentID)
	//if np != nil {
	//	np.SendMessage(msg)
	//}
}
