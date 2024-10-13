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

// INetHandler kernel 需要实现
type INetHandler interface {
	OnOk()
	OnUpdate()
	OnServerOk()
	OnNPAdd(np *network.NetPoint)

	OnNetError(np *network.NetPoint, acceptor *network.Acceptor)
	OnNetConnected(np *network.NetPoint)
	OnRpcNetAccept(np *network.NetPoint, acceptor *network.Acceptor)
}

type MsgHandler interface {
	INetHandler

	OnInit(baseModule *BaseModule)
	Start()
	OnNetMessage(pack *network.Packet)
	MsgRegister(processor *network.Processor)
}

// INetEventHandler 网络事件处理器
type INetEventHandler interface {
	OnRpcNetAccept(args []interface{})
	OnRpcNetConnected(args []interface{})
	OnRpcNetError(args []interface{})
	OnRpcNetClose(args []interface{})
	OnRpcNetData(args []interface{})
	OnRpcNetMessage(args []interface{})
}
