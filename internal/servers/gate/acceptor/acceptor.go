// Package acceptor @Author evan_yxw
// @Date 2024/8/6 16:51:00
// @Desc
package acceptor

import (
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/golang/protobuf/proto"
)

type IAcceptor interface {
	//AddCloseEventRpc(rpc IClientCloseHandler)
	//RegistClientMessage(msgID uint16, rpc *xsf_rpc.Acceptor)

	// 断开一个客户端连接
	//DisconnectClient(id uint32, reason uint32)

	// 发送消息到网关
	//SendMessage2Agent(agentID uint32, msg xsf_net.IMessage)

	// 发送消息到客户端
	SendMessage2Client(packet *network.Packet, msg proto.Message)

	// 广播一个消息给所有客户端
	//Broadcast(msg xsf_net.IMessage)

	// 设置一个客户端的服务器转播id
	//SetServerID(clientID uint32, ep uint8, serverID uint32)
}

type acceptor struct {
	ID int64
}

func NewAcceptor() *acceptor {
	return &acceptor{}
}

func (m *acceptor) SendMessage2Client(packet *network.Packet, msg proto.Message) {
	message, _ := rpc.GetMessage(uint64(xsf_pb.SMSGID_GtA_Gt_ClientMessage))
	localMsg := message.(*xsf_pb.GtA_Gt_ClientMessage)

	localMsg.ClientId = make([]uint32, 1)
	localMsg.ClientId[0] = packet.Msg.RawID
	// fixMe
	packer := network.NewDefaultPacker()
	data, _ := packer.Pack(msg)
	//data, _ := proto.Marshal(msg)
	localMsg.ClientMessage = data

	//var CID server.ClientID
	//server.ID2Cid(packet.Msg.RawID, &CID)
	//m.SendMessage2Agent(uint32(CID.Gate), uint64(xsf_pb.SMSGID_GtA_Gt_ClientMessage), localMsg)
	m.SendMessage2Agent(packet.NetPoint.ID, uint64(xsf_pb.SMSGID_GtA_Gt_ClientMessage), localMsg)
}

func (m *acceptor) SendMessage2Agent(serverId uint32, msgId uint64, message proto.Message) {
	manager := module.GetManager(module.ModuleID_SM)
	np := manager.Get(serverId)
	if np != nil {
		np.SendMessage(message)
	}
	//manager := GetManager(int(a.id))
	//np := manager.Get(agentID)
	//if np != nil {
	//	np.SendMessage(msg)
	//}
}