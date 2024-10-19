// Package handler @Author evan_yxw
// @Date 2024/10/19 17:01:00
// @Desc
package handler

import (
	"github.com/evanyxw/monster-go/pkg/module/module_def"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/golang/protobuf/proto"
)

type ClientHandler interface {
	OnHandshakeTicker(netPoint *network.NetPoint)
	SendMessage(message proto.Message, options ...network.PackerOptions)
}

type GateAcceptorHandler interface {
	//SendHandshake(ck *kernel.ConnectorKernel)
	SendHandshake(ck network.IConn)
}

// HandlerNetEvent handler 接收网络事件
type HandlerNetEvent interface {
	OnNetError(np *network.NetPoint, acceptor *network.Acceptor)
	OnNetAccept(np *network.NetPoint, acceptor *network.Acceptor)
}

type HandlerClientNetEvent interface {
	OnNetConnected(np *network.NetPoint)
}

// MsgHandler 消息handler
type MsgHandler interface {
	HandlerEvent
	HandlerNetEvent
	HandlerClientNetEvent

	OnNetMessage(pack *network.Packet)
	MsgRegister(processor *network.Processor)
}

type Handler interface {
	HandlerEvent
	OnNetMessage(pack *network.Packet)
	MsgRegister(processor *network.Processor)
}

// HandlerEvent handler接收事件
type HandlerEvent interface {
	OnOk()
	OnUpdate()
	OnServerOk()
	OnInit(baseModule module_def.IBaseModule)
	Start()
}
