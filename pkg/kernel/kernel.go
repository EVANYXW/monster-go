package kernel

import (
	"fmt"
	handler2 "github.com/evanyxw/monster-go/pkg/handler"
	"github.com/evanyxw/monster-go/pkg/module/module_def"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
)

type Kernel struct {
	NoWaitStart    bool
	handler        handler2.Handler
	msgHandlerImpl handler2.MsgHandler
	processor      *network.Processor
	rpcAcceptor    *rpc.Acceptor
}

func NewKernel(rpcAcceptor *rpc.Acceptor, processor *network.Processor, opts ...Option) *Kernel {
	opt := NewOption()
	for _, f := range opts {
		f(opt)
	}

	kernel := &Kernel{
		NoWaitStart:    false,
		handler:        opt.GetHandler(),
		msgHandlerImpl: opt.GetHandlerImpl(),
		processor:      processor,
		rpcAcceptor:    rpcAcceptor,
	}
	//kernel.Init()
	return kernel
}

func (n *Kernel) Init(baseModule module_def.IBaseModule) bool {
	n.handler.OnInit(baseModule)
	return true
}

func (n *Kernel) DoRegister() {
	//n.handler.MsgRegister(n.processor)
	// fixMe login 服务器在register的时候会重复注册，导致报错
	if n.rpcAcceptor != nil && n.msgHandlerImpl != nil {
		n.rpcAcceptor.Regist(rpc.RPC_NET_ACCEPT, n.OnRpcNetAccept)
		n.rpcAcceptor.Regist(rpc.RPC_NET_CONNECTED, n.OnRpcNetConnected)
		n.rpcAcceptor.Regist(rpc.RPC_NET_ERROR, n.OnRpcNetError)
	}

	if n.handler != nil {
		n.handler.MsgRegister(n.processor)
	}
}

func (n *Kernel) GetNPManager() network.INPManager {
	return nil
}

func (n *Kernel) GetStatus() int {
	return 0
}

func (n *Kernel) DoRun() {
	n.handler.Start()
}

func (n *Kernel) DoWaitStart() {

}

func (n *Kernel) DoRelease() {

}

func (n *Kernel) Update() {
	n.handler.OnUpdate()
}

func (n *Kernel) OnOk() {
	n.handler.OnOk()
}

func (n *Kernel) OnStartClose() {

}

func (n *Kernel) DoClose() {

}

func (n *Kernel) OnStartCheck() int {
	return module_def.ModuleOk()
}

func (n *Kernel) OnCloseCheck() int {
	return module_def.ModuleOk()
}

func (n *Kernel) GetNoWaitStart() bool {
	return n.NoWaitStart
}

func (n *Kernel) MessageHandler(packet *network.Packet) {

}

func (n *Kernel) OnRpcNetAccept(args []interface{}) {
	if n.msgHandlerImpl != nil {
		np := args[0].(*network.NetPoint)
		acc := args[1].(*network.Acceptor)
		fmt.Println("OnRpcNetAccept ....")
		n.msgHandlerImpl.OnNetAccept(np, acc)
	}
}

func (n *Kernel) OnRpcNetConnected(args []interface{}) {
	if n.msgHandlerImpl != nil {
		np := args[0].(*network.NetPoint)
		fmt.Println("OnRpcNetConnected ....")
		n.msgHandlerImpl.OnNetConnected(np)
	}
}

func (n *Kernel) OnRpcNetError(args []interface{}) {
	if n.msgHandlerImpl != nil {
		np := args[0].(*network.NetPoint)
		acc := args[1].(*network.Acceptor)
		n.msgHandlerImpl.OnNetError(np, acc)
		fmt.Println("NetKernel OnRpcNetError np close")
	}
}

func (n *Kernel) OnRpcNetData(args []interface{}) {

}

func (n *Kernel) OnRpcNetMessage(args []interface{}) {

}
