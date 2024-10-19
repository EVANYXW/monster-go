package handler

import (
	"fmt"
	"github.com/evanyxw/monster-go/internal/servers/gate/client"
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/kernel"
	"github.com/evanyxw/monster-go/pkg/module/module_def"
	"github.com/evanyxw/monster-go/pkg/network"
	"net"
	"time"
)

type GateMsgHandler struct {
	ClientManager module_def.IGtClientManager
}

func NewGateMsg() *GateMsgHandler {
	kernel.ClientManager = client.NewClientManager()
	return &GateMsgHandler{
		ClientManager: kernel.ClientManager,
	}
}

func (m *GateMsgHandler) OnInit(baseModule module_def.IBaseModule) {

}

func (m *GateMsgHandler) Start() {

}

func (m *GateMsgHandler) OnNetMessage(pack *network.Packet) {

}

func (m *GateMsgHandler) OnNetConnected(np *network.NetPoint) {

}

func (m *GateMsgHandler) OnNetAccept(np *network.NetPoint, acceptor *network.Acceptor) {
	newClient, isNew := m.ClientManager.NewClient(np)
	if newClient != nil {
		if isNew {

		} else {
			fmt.Println(111)
		}
		async.Go(func() {
			np.Connect()
		})
		lastHeartbeat := newClient.GetLastHeartbeat()

		network.ServerHeartbeat(time.Duration(6), time.Duration(30), lastHeartbeat, func() {
			conn := np.Conn.(*net.TCPConn)
			clt := newClient.(*client.Client)
			m.ClientClose(clt)
			acceptor.RemoveConn(conn, np)
			newClient.Close()
		})
	} else {
		//Error("Client Login is full")
		fmt.Println("OnRpcNetAccept NewClient newClient is nil")
		np.Close()
	}
}

func (m *GateMsgHandler) OnNetError(np *network.NetPoint, acceptor *network.Acceptor) {
	conn := np.Conn.(*net.TCPConn)
	acceptor.RemoveConn(conn, np)
	// 连接发生问题后,没有主动去关闭 client.Close()、心跳没有回会断掉
}

func (m *GateMsgHandler) OnServerOk() {

}

func (m *GateMsgHandler) OnOk() {

}

func (m *GateMsgHandler) OnUpdate() {

}

func (m *GateMsgHandler) OnNPAdd(np *network.NetPoint) {

}

func (m *GateMsgHandler) MsgRegister(processor *network.Processor) {
}

func (m *GateMsgHandler) ClientClose(client *client.Client) {
	for i := 0; i < len(client.GetServerIds()); i++ {
		connector := client.GetExistConnector(uint32(i))
		if connector != nil {
			//xsf_log.Infof("ClientKernel OnRpcNetError send Gt_GtA_ClientClose, i=%v, Client id=%v", i, client.ID.Get())
			msg := &xsf_pb.Gt_GtA_ClientClose{}
			msg.ClientId = client.ID.Load()
			connector.SendMessage(msg, network.WithRaID(client.ID.Load()))
		}
	}
}
