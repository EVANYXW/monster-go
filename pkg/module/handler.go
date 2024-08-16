package module

import (
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/golang/protobuf/proto"
)

type ClientHandler interface {
	OnHandshakeTicker(netPoint *network.NetPoint)
	SendMessage(msgId uint64, message proto.Message)
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
	OnNetError(np *network.NetPoint)
	OnNetConnected(np *network.NetPoint)
	OnRpcNetAccept(np *network.NetPoint)
	OnServerOk()
	OnNPAdd(np *network.NetPoint)
}

// 网络事件处理器
type INetEventHandler interface {
	OnRpcNetAccept(args []interface{})
	OnRpcNetConnected(args []interface{})
	OnRpcNetError(args []interface{})
	OnRpcNetData(args []interface{})
	OnRpcNetMessage(args []interface{})
}
