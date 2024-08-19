package handler

import (
	"fmt"
	"github.com/evanyxw/monster-go/internal/servers"
	"github.com/evanyxw/monster-go/internal/servers/gate/client"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
)

type GateMsgHandler struct {
	ClientManager module.ClientManager
}

func New() *GateMsgHandler {
	servers.ClientManager = client.NewClientManager()
	return &GateMsgHandler{
		ClientManager: servers.ClientManager,
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

func (m *GateMsgHandler) OnRpcNetAccept(np *network.NetPoint) {
	//c := client.NewClient(np)
	//c.Init()
	//fmt.Println("client:", c)
	//np.Connect()
	newClient, isNew := m.ClientManager.NewClient(np)
	if newClient != nil {

		if isNew {

		} else {

		}

		async.Go(func() {
			np.Connect()
		})
	} else {
		//Error("Client Login is full")
		fmt.Println("OnRpcNetAccept NewClient newClient is nil")
		np.Close()
	}
}

func (m *GateMsgHandler) OnNetError(np *network.NetPoint) {

}

func (m *GateMsgHandler) OnServerOk() {

}

func (m *GateMsgHandler) OnOk() {

}

func (m *GateMsgHandler) OnNPAdd(np *network.NetPoint) {

}

func (m *GateMsgHandler) MsgRegister(processor *network.Processor) {
}
