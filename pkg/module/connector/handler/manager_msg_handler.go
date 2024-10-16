package handler

import (
	"fmt"
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

type managerMsgHandler struct {
}

func NewManagerMsg() *managerMsgHandler {
	return &managerMsgHandler{}
}

func (m *managerMsgHandler) OnInit(baseModule *module.BaseModule) {

}

func (m *managerMsgHandler) Start() {

}

func (m *managerMsgHandler) OnNetMessage(pack *network.Packet) {

}

func (m *managerMsgHandler) OnNetConnected(np *network.NetPoint) {
	fmt.Println("Gate ManagerMsgHandler OnNetConnected..............")
	messageID := uint64(xsf_pb.SMSGID_Gt_GtA_Handshake)
	msg, _ := rpc.GetMessage(messageID)
	localMsg := msg.(*xsf_pb.Gt_GtA_Handshake)
	localMsg.ServerId = server.ID
	//pack, _ := servers.ConnectorKernel.Client.Pack(messageID, localMsg)
	//servers.ConnectorKernel.NetPoint.SetSignal(pack)

	np.SendMessage(localMsg)
}

func (m *managerMsgHandler) OnNetAccept(np *network.NetPoint, acceptor *network.Acceptor) {
	np.Connect()
	conn := np.Conn.(*net.TCPConn)
	acceptor.RemoveConn(conn, np)
}

func (m *managerMsgHandler) OnNetError(np *network.NetPoint, acceptor *network.Acceptor) {

}

func (m *managerMsgHandler) OnServerOk() {

}

func (m *managerMsgHandler) OnOk() {

}

func (m *managerMsgHandler) OnUpdate() {

}

func (m *managerMsgHandler) MsgRegister(processor *network.Processor) {
	fmt.Println("GtA_Gt_Handshake MsgRegister laile")
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_GtA_Gt_Handshake), m.GtA_Gt_Handshake)
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_GtA_Gt_ClientMessage), m.GtA_Gt_ClientMessage)
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_GtA_Gt_ClientDisconnect), m.GtA_Gt_ClientDisconnect)
	processor.RegisterMsg(uint16(xsf_pb.SMSGID_GtA_Gt_SetServerID), m.GtA_Gt_SetServerID)
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
		cnf := configs.All()
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
	rpc.Import(message.Msg.Data, clientMessage)
	fmt.Println("GtA_Gt_ClientMessage")

	for i := 0; i < len(clientMessage.ClientId); i++ {
		clt := module.ClientManager.GetClient(clientMessage.ClientId[i])
		if clt != nil && clt.GetID() > 0 {
			clt.SetSignal(clientMessage.GetClientMessage())
		}
	}
}

func (m *managerMsgHandler) GtA_Gt_ClientDisconnect(message *network.Packet) {
	clientDisconnect := &xsf_pb.GtA_Gt_ClientDisconnect{}
	rpc.Import(message.Msg.Data, clientDisconnect)
	fmt.Println("GtA_Gt_ClientDisconnect")

	clt := module.ClientManager.GetClient(clientDisconnect.ClientId)
	if clt != nil && clt.GetID() > 0 {
		//GoDisconnect(int32(localMsg.PB.Reason), true)
		clt.GoDisconnect(message.Msg.RawID)
	}

}

func (m *managerMsgHandler) GtA_Gt_SetServerID(message *network.Packet) {
	msg := &xsf_pb.GtA_Gt_SetServerID{}
	rpc.Import(message.Msg.Data, msg)

	clt := module.ClientManager.GetClient(msg.ClientId)
	if clt != nil && clt.GetID() > 0 {
		args := []interface{}{
			msg.Ep,
			msg.ServerId,
		}
		clt.SetServerID(args)
	}
}
