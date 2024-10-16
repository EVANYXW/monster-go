package module

import (
	"fmt"
	"github.com/evanyxw/monster-go/pkg/client"
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
	msgHandler  MsgHandler
	processor   *network.Processor
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

func NewConnectorKernel(ip string, port uint32, msgHandler MsgHandler, packerFactory network.PackerFactory,
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

func (c *ConnectorKernel) Init(baseModule IBaseModule) bool {
	c.runStatus = ModuleRunStatus_Start
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
	c.Client.Run(c.RpcAcceptor)
	c.runStatus = ModuleRunStatus_Running
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

}

func (c *ConnectorKernel) DoClose() {

}

func (c *ConnectorKernel) OnStartCheck() int {
	if c.runStatus == ModuleRunStatus_Running {
		return ModuleRunCode_Ok
	}
	return ModuleRunCode_Wait
}

func (c *ConnectorKernel) GetNoWaitStart() bool {
	return c.NoWaitStart
}

func (c *ConnectorKernel) OnCloseCheck() int {
	return 0
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
	np := args[0].(*network.NetPoint)
	c.msgHandler.OnNetConnected(np)
}

func (c *ConnectorKernel) OnRpcNetError(args []interface{}) {
	//np := args[0].(*network.NetPoint)
	//
	//if c.msgHandler != nil {
	//	c.msgHandler.OnNetError(np, nil)
	//}
	//fmt.Println("ConnectorKernel OnRpcNetError np close")
	//np.Close()

	//fixMe OnRpcNetError 还没做其他处理!!!
	fmt.Println("OnRpcNetError 还没做其他处理!!!")
}

func (c *ConnectorKernel) OnRpcNetClose(args []interface{}) {

}

func (c *ConnectorKernel) OnRpcNetData(args []interface{}) {

}

func (c *ConnectorKernel) OnRpcNetMessage(args []interface{}) {

}
