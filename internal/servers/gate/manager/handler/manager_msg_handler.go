package handler

import (
	"fmt"
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/internal/servers"
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"
	"github.com/golang/protobuf/proto"
	"time"
)

type managerMsgHandler struct {
}

func NewManager() *managerMsgHandler {
	return &managerMsgHandler{}
}

func (m *managerMsgHandler) Start() {

}

//func (m *managerMsgHandler) GetIsHandle() bool {
//	return m.isHandle
//}

func (m *managerMsgHandler) OnNetMessage(pack *network.Packet) {

}

func (m *managerMsgHandler) OnNetConnected(np *network.NetPoint) {

}

func (m *managerMsgHandler) OnRpcNetAccept(np *network.NetPoint) {
	np.Connect()
}

func (m *managerMsgHandler) OnNetError(np *network.NetPoint) {

}

func (m *managerMsgHandler) OnServerOk() {

}

func (m *managerMsgHandler) OnOk() {

}

func (m *managerMsgHandler) OnNPAdd(np *network.NetPoint) {

}

func (m *managerMsgHandler) MsgRegister(processor *network.Processor) {
	fmt.Println("GtA_Gt_Handshake MsgRegister laile")
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_GtA_Gt_Handshake), m.GtA_Gt_Handshake)
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_GtA_Gt_ClientMessage), m.GtA_Gt_ClientMessage)
}

func (m *managerMsgHandler) SendHandshake(ck *module.ConnectorKernel) {
	messageID := uint64(xsf_pb.SMSGID_Gt_GtA_Handshake)
	msg, _ := rpc.GetMessage(messageID)
	localMsg := msg.(*xsf_pb.Gt_GtA_Handshake)
	localMsg.ServerId = server.ID
	//pack, _ := servers.ConnectorKernel.Client.Pack(messageID, localMsg)
	//servers.ConnectorKernel.NetPoint.SetSignal(pack)

	ck.NetPoint.SendMessage(localMsg)
}

func (m *managerMsgHandler) OnHandshakeTicker(np *network.NetPoint) {
	async.Go(func() {
		cnf := configs.Get()
		ticker := time.NewTicker(time.Second * time.Duration(cnf.HtCheck))
		defer ticker.Stop()

		for range ticker.C {
			np.SendMessage(&xsf_pb.Gt_GtA_Heartbeat{})
		}
	})
}

func (m *managerMsgHandler) SendMessage(msgId uint64, message proto.Message) {

}

func (m *managerMsgHandler) GtA_Gt_Handshake(message *network.Packet) {
	fmt.Println("GtA_Gt_Handshake laile")
	logger.Info("GtA_Gt_Handshake 收到")
	m.OnHandshakeTicker(message.NetPoint)
}

// GtA_Gt_ClientMessage gate accepts push client messages
func (m *managerMsgHandler) GtA_Gt_ClientMessage(message *network.Packet) {
	clientMessage := &xsf_pb.GtA_Gt_ClientMessage{}
	proto.Unmarshal(message.Msg.Data, clientMessage)
	fmt.Println("GtA_Gt_ClientMessage")

	for i := 0; i < len(clientMessage.ClientId); i++ {
		clt := servers.ClientManager.GetClient(clientMessage.ClientId[i])
		if clt != nil && clt.GetID() > 0 {
			clt.SetSignal(clientMessage.GetClientMessage())
		}
	}
}
