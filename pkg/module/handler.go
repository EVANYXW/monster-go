package module

import "github.com/evanyxw/monster-go/pkg/network"

type MsgHandler interface {
	INetHandler
	Start()
	GetIsHandle() bool
	HandleMsg(pack *network.Packet)
	//MsgRegister(kernel *NetKernel)
	MsgRegister(kernel *network.Processor)
}
