package module

import (
	"fmt"
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/client"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"
	"sync"
)

type ConnectorKernel struct {
	Owner INetHandler
	*client.Client
	handlers    network.HandlerMap
	RpcAcceptor *rpc.Acceptor
	ID          uint32
	SID         server.ServerID
	wg          sync.WaitGroup
	runStatus   int
	//forceClose atomic.Bool
	//Addr string
	//Helper IConnectorHelper
	//conn *net.Conn
	//hbTimer *xsf_timer.Timer
}

func NewConnectorKernel(owner INetHandler, ip string, port uint32) *ConnectorKernel {
	rpcAcceptor := rpc.NewAcceptor(10000)
	connector := &ConnectorKernel{
		handlers:    make(network.HandlerMap, xsf_pb.SMSGID_Server_Max),
		Client:      client.NewClient(fmt.Sprintf("%s:%d", ip, port), rpcAcceptor),
		RpcAcceptor: rpcAcceptor,
	}
	connector.Client.OnMessageCb = connector.MessageHandler
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

func (c *ConnectorKernel) AddModules() {
	//module, ok := c.Owner.(IModule)
	//if ok {
	//	AddModule(module)
	//}
}

func (c *ConnectorKernel) DoRegist() {
	c.RpcAcceptor.Regist(rpc.RPC_NET_ACCEPT, c.OnRpcNetAccept)
	c.RpcAcceptor.Regist(rpc.RPC_NET_ERROR, c.OnRpcNetError)
}

func (c *ConnectorKernel) Start() {
	// FixMe 携程会泄漏
	c.RpcAcceptor.Run()
	c.Client.Run()
	c.runStatus = ModuleRunStatus_Running
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
	return false
}

func (c *ConnectorKernel) OnCloseCheck() int {
	return 0
}

func (c *ConnectorKernel) RegisterMsg(msgId uint16, handlerFunc network.HandlerFunc) {
	c.handlers[msgId] = handlerFunc
}

func (c *ConnectorKernel) MessageHandler(packet *network.Packet) {
	handler := c.handlers[packet.Msg.ID]
	handler(packet)
}

func (c *ConnectorKernel) OnRpcNetAccept(args []interface{}) {

}

func (c *ConnectorKernel) OnRpcNetConnected(args []interface{}) {

}

func (c *ConnectorKernel) OnRpcNetError(args []interface{}) {
	np := args[0].(*network.NetPoint)
	// connector manager 就传的nil
	if c.Owner != nil {
		c.Owner.OnNetError(np)
	}
	np.Close()
}

func (c *ConnectorKernel) OnRpcNetData(args []interface{}) {

}

func (c *ConnectorKernel) OnRpcNetMessage(args []interface{}) {

}
