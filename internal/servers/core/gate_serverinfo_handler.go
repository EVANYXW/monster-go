package core

import (
	"fmt"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"go.uber.org/zap"
)

type gateServerInfoHandler struct {
}

func NewGateServerInfoHandler() *gateServerInfoHandler {
	return &gateServerInfoHandler{}
}

func (h *gateServerInfoHandler) OnServerNew(Info *network.ServerInfo) {

}

// OnServerLost 有一个服务器断开
func (h *gateServerInfoHandler) OnServerLost(id uint32) {

}

// OnServerOk 服务器已准备好
func (h *gateServerInfoHandler) OnServerOk(info *network.ServerInfo) {
	logger.Info("siHandler OnServerOk", zap.Uint("id", uint(network.ID)))

	var SID network.ServerID
	network.ID2Sid(info.ID, &SID)

	if SID.Type == network.EP_Login || SID.Type == network.EP_Game || SID.Type == network.EP_Mail || SID.Type == network.EP_World {
		moduler := module.GetModule(module.ModuleID_ConnectorManager).GetOwner()
		connectorManager := moduler.(*ConnectorManager)

		fmt.Println("ports:", network.Ports)
		conn := connectorManager.CreateConnector(info.ID, info.IP, info.Ports[network.EP_Gate])
		if conn == nil {
			logger.Error("siHandler OnServerOk create connector error", zap.Uint("server", uint(SID.ID)),
				zap.String("type", network.EP2Name(SID.Type)),
				zap.Uint("index", uint(SID.Index)))
		}
	}

	// todo
	//if SID.Type == network.EP_Mail {
	//	si.mailID.Set(Info.ID)
	//} else if SID.Type == xsf_util.EP_Manager {
	//	si.managerID.Set(Info.ID)
	//}
}

// OnServerOpenComplete
func (h *gateServerInfoHandler) OnServerOpenComplete() {

}
