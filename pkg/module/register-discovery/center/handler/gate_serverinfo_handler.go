package handler

import (
	"fmt"
	"github.com/evanyxw/monster-go/pkg/ipPort"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/module/module_def"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/server"
	"go.uber.org/zap"
	"sync/atomic"
)

type gateServerInfoHandler struct {
	mailID    atomic.Uint32
	managerID atomic.Uint32
}

func NewGateServerInfoHandler() *gateServerInfoHandler {
	g := &gateServerInfoHandler{}
	return g
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
		owner := module_def.GetModule(module_def.ModuleID_ConnectorManager).GetOwner()
		connectorManager := owner.(*module.Manager)

		fmt.Println("ports:", server.Ports)
		// fixMe 这里暂时用ip，要改注册到center的时候用ip
		ip, _ := ipPort.ExternalIP()
		ipstr := fmt.Sprintf("%s", ip)
		//conn := connectorManager.CreateConnector(info.ID, info.IP, info.Ports[server.EP_Gate])
		conn := connectorManager.CreateConnector(info.ID, ipstr, info.Ports[server.EP_Gate])
		if conn == nil {
			logger.Error("siHandler OnServerOk create connector error", zap.Uint("server", uint(SID.ID)),
				zap.String("type", server.EP2Name(SID.Type)),
				zap.Uint("index", uint(SID.Index)))
		}
	}

	if SID.Type == server.EP_Mail {
		//h.mailID.Store(info.ID)
		module_def.MailID.Store(info.ID)
	} else if SID.Type == server.EP_Manager {
		module_def.ManagerID.Store(info.ID)
		//h.managerID.Store(info.ID)
	}
}

// OnServerOpenComplete
func (h *gateServerInfoHandler) OnServerOpenComplete() {

}
