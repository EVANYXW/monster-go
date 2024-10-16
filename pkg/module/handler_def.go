package module

import (
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/golang/protobuf/proto"
	"sync/atomic"
)

var (
	MailID    atomic.Uint32
	ManagerID atomic.Uint32
)

type ClientHandler interface {
	OnHandshakeTicker(netPoint *network.NetPoint)
	SendMessage(message proto.Message, options ...network.PackerOptions)
}

type GateAcceptorHandler interface {
	SendHandshake(ck *ConnectorKernel)
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
	OnInit(baseModule *BaseModule)
	Start()
}

// KernelNetEvent kernel网络事件处理器
type KernelNetEvent interface {
	OnRpcNetAccept(args []interface{})
	OnRpcNetConnected(args []interface{})
	OnRpcNetError(args []interface{})
	OnRpcNetClose(args []interface{})
	OnRpcNetData(args []interface{})
	OnRpcNetMessage(args []interface{})
}
