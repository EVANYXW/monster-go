package module

import (
	"fmt"
	"github.com/evanyxw/monster-go/pkg/async"
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

func NewConnectorKernel(ip string, port uint32, msgHandler MsgHandler, packerFactory network.PackerFactory, options ...ckernelOption) *ConnectorKernel {
	rpcAcceptor := rpc.NewAcceptor(10000)
	processor := network.NewProcessor()
	connector := &ConnectorKernel{
		//handlers:    make(network.HandlerMap, xsf_pb.SMSGID_Server_Max),
		processor:   processor,
		Client:      client.NewClient(fmt.Sprintf("%s:%d", ip, port), rpcAcceptor, processor, packerFactory),
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
	c.RpcAcceptor.Regist(rpc.RPC_NET_CONNECTED, c.OnRpcNetConnected)
	c.RpcAcceptor.Regist(rpc.RPC_NET_ERROR, c.OnRpcNetError)

	if c.msgHandler != nil {
		c.msgHandler.MsgRegister(c.processor)
	}
}

func (c *ConnectorKernel) Start() {
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
	np := args[0].(*network.NetPoint)
	async.Go(func() {
		np.CloseChan <- true
	})
	close(np.CloseChan)
	// connector manager 就传的nil
	if c.msgHandler != nil {
		c.msgHandler.OnNetError(np)
	}
	fmt.Println("ConnectorKernel OnRpcNetError np close")
	np.Close()
}

func (c *ConnectorKernel) OnRpcNetData(args []interface{}) {

}

func (c *ConnectorKernel) OnRpcNetMessage(args []interface{}) {

}
