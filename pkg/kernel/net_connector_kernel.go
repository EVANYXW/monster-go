package kernel

import (
	"fmt"
	"github.com/evanyxw/monster-go/pkg/client"
	"github.com/evanyxw/monster-go/pkg/handler"
	"github.com/evanyxw/monster-go/pkg/module/module_def"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"
	"sync"
)

type ConnectorKernel struct {
	*client.Client

	RpcAcceptor *rpc.Acceptor
	ID          uint32
	SID         server.ServerID
	wg          sync.WaitGroup
	runStatus   int
	NoWaitStart bool
	msgHandler  handler.MsgHandler
	processor   *network.Processor
	baseModule  module_def.IBaseModule
}
type ckOptions struct {
	NoWaitStart bool
}

type ckernelOption func(kernel *ckOptions)

func WithCNoWaitStart(noWaitStart bool) ckernelOption {
	return func(kernel *ckOptions) {
		kernel.NoWaitStart = noWaitStart
	}
}

func NewConnectorKernel(ip string, port uint32, msgHandler handler.MsgHandler, packerFactory network.PackerFactory,
	options ...ckernelOption) *ConnectorKernel {
	opt := &ckOptions{}
	rpcAcceptor := rpc.NewAcceptor(10000)
	processor := network.NewProcessor()
	connector := &ConnectorKernel{
		processor:   processor,
		Client:      client.NewClient(fmt.Sprintf("%s:%d", ip, port), processor, packerFactory),
		RpcAcceptor: rpcAcceptor,
		NoWaitStart: false,
		msgHandler:  msgHandler,
	}
	connector.Client.OnMessageCb = connector.MessageHandler

	for _, fn := range options {
		fn(opt)
	}
	connector.NoWaitStart = opt.NoWaitStart
	return connector
}

func (c *ConnectorKernel) SetID(id uint32) {
	c.ID = id
	server.ID2Sid(id, &c.SID)
}

func (c *ConnectorKernel) Init(baseModule module_def.IBaseModule) bool {
	c.baseModule = baseModule
	c.runStatus = module_def.ModuleRunStatus_Start
	return true
}

func (c *ConnectorKernel) DoRegister() {
	c.RpcAcceptor.Regist(rpc.RPC_NET_CONNECTED, c.OnRpcNetConnected)
	c.RpcAcceptor.Regist(rpc.RPC_NET_ERROR, c.OnRpcNetError)

	if c.msgHandler != nil {
		c.msgHandler.MsgRegister(c.processor)
	}
}

func (c *ConnectorKernel) DoRun() {
	c.RpcAcceptor.Run()
	err := c.Client.Run(c.RpcAcceptor)
	if err != nil {
		return
	}

	c.runStatus = module_def.ModuleRunStatus_Running
	c.msgHandler.Start()
}

func (c *ConnectorKernel) DoWaitStart() {

}

func (c *ConnectorKernel) DoRelease() {
	c.Client.OnClose()
}

func (c *ConnectorKernel) Update() {

}

func (c *ConnectorKernel) OnOk() {
	c.msgHandler.OnOk()
}

func (c *ConnectorKernel) OnStartClose() {
	c.Close()
}

func (c *ConnectorKernel) DoClose() {
	c.DoRelease()
}

func (c *ConnectorKernel) OnStartCheck() int {
	if c.runStatus == module_def.ModuleRunStatus_Running {
		return module_def.ModuleOk()
	}
	return module_def.ModuleWait()
}

func (c *ConnectorKernel) GetNoWaitStart() bool {
	return c.NoWaitStart
}

func (c *ConnectorKernel) OnCloseCheck() int {
	return module_def.ModuleOk()
}

func (c *ConnectorKernel) GetNPManager() network.INPManager {
	return nil
}

func (c *ConnectorKernel) GetStatus() int {
	return 0
}

func (c *ConnectorKernel) RegisterMsg(msgId uint16, handlerFunc network.HandlerFunc) {
	c.processor.RegisterMsg(msgId, handlerFunc)
}

func (c *ConnectorKernel) MessageHandler(packet *network.Packet) {
	c.processor.MessageHandler(packet)
}

func (c *ConnectorKernel) OnRpcNetAccept(args []interface{}) {

}

func (c *ConnectorKernel) OnRpcNetConnected(args []interface{}) {
	if c.msgHandler == nil {
		return
	}
	np := args[0].(*network.NetPoint)
	c.msgHandler.OnNetConnected(np)
}

func (c *ConnectorKernel) OnRpcNetError(args []interface{}) {
	if c.msgHandler == nil {
		return
	}
	np := args[0].(*network.NetPoint)
	acc := args[1].(*network.Acceptor)
	c.msgHandler.OnNetError(np, acc)
}

func (c *ConnectorKernel) OnRpcNetClose(args []interface{}) {

}

func (c *ConnectorKernel) OnRpcNetData(args []interface{}) {

}

func (c *ConnectorKernel) OnRpcNetMessage(args []interface{}) {

}
