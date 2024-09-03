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

type MsgHandler interface {
	INetHandler

	Start()
	OnNetMessage(pack *network.Packet)
	MsgRegister(processor *network.Processor)
}

// INetHandler kernel 需要实现
type INetHandler interface {
	//Start()
	//DoRegist(nk *NetKernel)
	//OnNetData(np *network.NetPoint, msgID uint16, rawID uint32, data []byte)
	//OnNetMessage(np *network.NetPoint, msgID uint16, rawID uint32, message xsf_net.IMessage)
	//OnStartCheck() int
	//OnCloseCheck() int
	//DoClose()
	//OnStartClose()
	OnOk()
	OnUpdate()
	OnServerOk()
	OnNPAdd(np *network.NetPoint)

	OnNetError(np *network.NetPoint, acceptor *network.Acceptor)
	OnNetConnected(np *network.NetPoint)
	OnRpcNetAccept(np *network.NetPoint, acceptor *network.Acceptor)
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
