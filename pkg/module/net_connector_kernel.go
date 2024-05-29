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
	//handlers    network.HandlerMap
	//forceClose atomic.Bool
	//Addr string
	//Helper IConnectorHelper
	//conn *net.Conn
	//hbTimer *xsf_timer.Timer
}

type ckernelOption func(kernel *ConnectorKernel)

func WithCNoWaitStart(noWaitStart bool) ckernelOption {
	return func(kernel *ConnectorKernel) {
		kernel.NoWaitStart = noWaitStart
	}
}

func NewConnectorKernel(ip string, port uint32, msgHandler MsgHandler, options ...ckernelOption) *ConnectorKernel {
	rpcAcceptor := rpc.NewAcceptor(10000)
	connector := &ConnectorKernel{
		//handlers:    make(network.HandlerMap, xsf_pb.SMSGID_Server_Max),
		processor:   network.NewProcessor(),
		Client:      client.NewClient(fmt.Sprintf("%s:%d", ip, port), rpcAcceptor),
		RpcAcceptor: rpcAcceptor,
		NoWaitStart: false,
		msgHandler:  msgHandler,
	}
	connector.Client.OnMessageCb = connector.MessageHandler

	for _, fn := range options {
		fn(connector)
	}

	return connector
}

func (c *ConnectorKernel) SetID(id uint32) {
	c.ID = id
	server.ID2Sid(id, &c.SID)
}

func (c *ConnectorKernel) Init() bool {
	c.runStatus = ModuleRunStatus_Start
	return true
}

func (c *ConnectorKernel) DoRegist() {
	c.RpcAcceptor.Regist(rpc.RPC_NET_ACCEPT, c.OnRpcNetAccept)
	c.RpcAcceptor.Regist(rpc.RPC_NET_ERROR, c.OnRpcNetError)

	if c.msgHandler != nil {
		c.msgHandler.MsgRegister(c.processor)
	}
}

func (c *ConnectorKernel) Start() {
	// FixMe 携程会泄漏
	c.RpcAcceptor.Run()
	c.Client.Run()
	c.runStatus = ModuleRunStatus_Running
	c.msgHandler.Start()
}

func (c *ConnectorKernel) DoStart() {

}

func (c *ConnectorKernel) Release() {
	c.Client.OnClose()
}

func (c *ConnectorKernel) OnUpdate(timeDelta uint32) {

}

func (c *ConnectorKernel) OnOK() {

}

func (c *ConnectorKernel) OnStartClose() {

}

func (c *ConnectorKernel) DoClose() {

}

func (c *ConnectorKernel) OnStartCheck() int {
	if c.runStatus == ModuleRunStatus_Start {
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

func (c *ConnectorKernel) RegisterMsg(msgId uint16, handlerFunc network.HandlerFunc) {
	//c.handlers[msgId] = handlerFunc
	c.processor.RegisterMsg(msgId, handlerFunc)
}

func (c *ConnectorKernel) MessageHandler(packet *network.Packet) {
	//handler := c.handlers[packet.Msg.ID]
	//handler(packet)
	c.processor.MessageHandler(packet)
}

func (c *ConnectorKernel) OnRpcNetAccept(args []interface{}) {

}

func (c *ConnectorKernel) OnRpcNetConnected(args []interface{}) {

}

func (c *ConnectorKernel) OnRpcNetError(args []interface{}) {
	np := args[0].(*network.NetPoint)
	// connector manager 就传的nil
	if c.msgHandler != nil {
		c.msgHandler.OnNetError(np)
	}
	np.Close()
}

func (c *ConnectorKernel) OnRpcNetData(args []interface{}) {

}

func (c *ConnectorKernel) OnRpcNetMessage(args []interface{}) {

}
