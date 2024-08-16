package handler

import (
	"fmt"
	"github.com/evanyxw/monster-go/internal/servers"
	"github.com/evanyxw/monster-go/internal/servers/gate/client"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
)

type gateMsgHandler struct {
	ClientManager module.ClientManager
}

func New() *gateMsgHandler {
	servers.ClientManager = client.NewClientManager()
	return &gateMsgHandler{
		ClientManager: servers.ClientManager,
	}
}

func (m *gateMsgHandler) Start() {

}

//func (m *gateMsgHandler) GetIsHandle() bool {
//	return m.isHandle
//}

func (m *gateMsgHandler) OnNetMessage(pack *network.Packet) {

}

func (m *gateMsgHandler) OnNetConnected(np *network.NetPoint) {

}

func (m *gateMsgHandler) OnRpcNetAccept(np *network.NetPoint) {
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

func (m *gateMsgHandler) OnNetError(np *network.NetPoint) {

}

func (m *gateMsgHandler) OnServerOk() {

}

func (m *gateMsgHandler) OnOk() {

}

func (m *gateMsgHandler) OnNPAdd(np *network.NetPoint) {

}

func (m *gateMsgHandler) MsgRegister(processor *network.Processor) {

}
