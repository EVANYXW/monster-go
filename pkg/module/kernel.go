package module

import (
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
)

type Kernel struct {
	NoWaitStart bool
	msgHandler  MsgHandler
	processor   *network.Processor
	rpcAcceptor *rpc.Acceptor
}

func NewKernel(msgHandler MsgHandler, rpcAcceptor *rpc.Acceptor, processor *network.Processor) *Kernel {
	kernel := &Kernel{
		NoWaitStart: false,
		msgHandler:  msgHandler,
		processor:   processor,
		rpcAcceptor: rpcAcceptor,
	}
	//kernel.Init()
	return kernel
}

func (n *Kernel) Init() bool {
	return true
}

func (n *Kernel) DoRegister() {
	//n.msgHandler.MsgRegister(n.processor)
	// fixMe login 服务器在register的时候会重复注册，导致报错
	//if n.rpcAcceptor != nil {
	//	n.rpcAcceptor.Regist(rpc.RPC_NET_ACCEPT, n.OnRpcNetAccept)
	//	n.rpcAcceptor.Regist(rpc.RPC_NET_CONNECTED, n.OnRpcNetConnected)
	//	n.rpcAcceptor.Regist(rpc.RPC_NET_ERROR, n.OnRpcNetError)
	//}

	if n.msgHandler != nil {
		n.msgHandler.MsgRegister(n.processor)
	}
}

func (n *Kernel) GetNPManager() network.INPManager {
	return nil
}

func (n *Kernel) GetStatus() int {
	return 0
}

func (n *Kernel) DoRun() {
	n.msgHandler.Start()
}

func (n *Kernel) DoWaitStart() {

}

func (n *Kernel) DoRelease() {

}

func (n *Kernel) Update() {

}

func (n *Kernel) OnOk() {
	n.msgHandler.OnOk()
}

func (n *Kernel) OnStartClose() {

}

func (n *Kernel) DoClose() {

}

func (n *Kernel) OnStartCheck() int {
	return 0
}

func (n *Kernel) OnCloseCheck() int {
	return 0
}

func (n *Kernel) GetNoWaitStart() bool {
	return n.NoWaitStart
}

func (n *Kernel) MessageHandler(packet *network.Packet) {

}

func (n *Kernel) OnRpcNetAccept(args []interface{}) {

}

func (n *Kernel) OnRpcNetConnected(args []interface{}) {

}

func (n *Kernel) OnRpcNetError(args []interface{}) {

}

func (n *Kernel) OnRpcNetData(args []interface{}) {

}

func (n *Kernel) OnRpcNetMessage(args []interface{}) {

}
