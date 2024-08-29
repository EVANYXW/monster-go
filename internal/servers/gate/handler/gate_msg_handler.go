package handler

import (
	"fmt"
	"github.com/evanyxw/monster-go/internal/servers/gate/client"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"net"
	"time"
)

type GateMsgHandler struct {
	ClientManager module.IClientManager
}

func NewMsg() *GateMsgHandler {
	module.ClientManager = client.NewClientManager()
	return &GateMsgHandler{
		ClientManager: module.ClientManager,
	}
}

func (m *GateMsgHandler) Start() {

}

//func (m *GateMsgHandler) GetIsHandle() bool {
//	return m.isHandle
//}

func (m *GateMsgHandler) OnNetMessage(pack *network.Packet) {

}

func (m *GateMsgHandler) OnNetConnected(np *network.NetPoint) {

}

func (m *GateMsgHandler) OnRpcNetAccept(np *network.NetPoint, acceptor *network.Acceptor) {
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

func (m *GateMsgHandler) OnNPAdd(np *network.NetPoint) {

}

func (m *GateMsgHandler) MsgRegister(processor *network.Processor) {
}
