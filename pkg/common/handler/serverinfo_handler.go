package handler

import (
	"github.com/evanyxw/monster-go/pkg/network"
)

type serverInfoHandler struct {
}

func NewServerInfoHandler() *serverInfoHandler {
	return &serverInfoHandler{}
}

func (h *serverInfoHandler) OnServerNew(Info *network.ServerInfo) {

}

// OnServerLost 有一个服务器断开
func (h *serverInfoHandler) OnServerLost(id uint32) {

}

// OnServerOk 服务器已准备好
func (h *serverInfoHandler) OnServerOk(info *network.ServerInfo) {
}

// OnServerOpenComplete
func (h *serverInfoHandler) OnServerOpenComplete() {

}
