package handler

import (
	"fmt"
	"github.com/evanyxw/monster-go/internal/servers/gate/manager"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/server"
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
	logger.Info("siHandler OnServerOk", zap.Uint("id", uint(server.ID)))

	var SID server.ServerID
	server.ID2Sid(info.ID, &SID)

	if SID.Type == server.EP_Login || SID.Type == server.EP_Game || SID.Type == server.EP_Mail || SID.Type == server.EP_World {
		owner := module.GetModule(module.ModuleID_ConnectorManager).GetOwner()
		connectorManager := owner.(*manager.ConnectorManager)

		fmt.Println("ports:", server.Ports)
		conn := connectorManager.CreateConnector(info.ID, info.IP, info.Ports[server.EP_Gate])
		if conn == nil {
			logger.Error("siHandler OnServerOk create connector error", zap.Uint("server", uint(SID.ID)),
				zap.String("type", server.EP2Name(SID.Type)),
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
