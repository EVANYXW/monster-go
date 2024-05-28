package module

import (
	"fmt"
	"github.com/evanyxw/monster-go/internal/pkg/output"
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"
)

type NetType int

const (
	Inner NetType = iota
	Outer
)

type kernelOption func(kernel *NetKernel)

func WithNoWaitStart(noWaitStart bool) kernelOption {
	return func(kernel *NetKernel) {
		kernel.NoWaitStart = noWaitStart
	}
}

func WithNetType(netType NetType) kernelOption {
	return func(kernel *NetKernel) {
		kernel.netType = netType
	}
}

type NetKernel struct {
	Owner       INetHandler
	netType     NetType
	NetAcceptor *network.Acceptor
	RpcAcceptor *rpc.Acceptor
	handlers    network.HandlerMap
	NPManager   network.INPManager
	closeChan   chan struct{}
	port        uint32
	netMaxCount uint32
	NoWaitStart bool
}

func NewNetKernel(maxConnNum uint32, info server.Info, Owner INetHandler, options ...kernelOption) *NetKernel {
	rpcAcceptor := rpc.NewAcceptor(10000)
	nodePointManager := network.NewNormal(maxConnNum, rpcAcceptor)

	kernel := &NetKernel{
		NPManager:   nodePointManager,
		NetAcceptor: network.NewAcceptor(maxConnNum, info, rpcAcceptor, nodePointManager),
		RpcAcceptor: rpcAcceptor,
		handlers:    make(network.HandlerMap, xsf_pb.SMSGID_Server_Max),
		closeChan:   make(chan struct{}),
		Owner:       Owner,
		NoWaitStart: false,
		netType:     Inner, // 默认内网
	}

	for _, fn := range options {
		fn(kernel)
	}

	kernel.NetAcceptor.MessageHandler = kernel.MessageHandler
	kernel.Init()
	return kernel
}

func (n *NetKernel) Init() bool {
	n.AddModules()
	return true
}

func (n *NetKernel) AddModules() {
	AddManager(ModuleID_SM, n.NPManager)
}

func (n *NetKernel) DoRegist() {
	n.RpcAcceptor.Regist(rpc.RPC_NET_ACCEPT, n.OnRpcNetAccept)
	n.RpcAcceptor.Regist(rpc.RPC_NET_ERROR, n.OnRpcNetError)
}

func (n *NetKernel) Start() {
	if n.NoWaitStart {
		async.Go(func() {
			n.NetAcceptor.Connect()
			n.NetAcceptor.Run()
		})
	}
}

func (n *NetKernel) DoStart() {
	port := server.Ports[server.EP_Client]
	if n.netType == Inner {
		port = server.Ports[server.EP_Gate]
	}
	addr := fmt.Sprintf(":%d", port)
	output.Oput.SetServerAddr(addr)

	async.Go(func() {
		n.NetAcceptor.Connect(network.WithAddr(addr))
		n.NetAcceptor.Run()
	})
	//n.NetAcceptor.DoStart()
}

func (n *NetKernel) Release() {
	n.NetAcceptor.OnClose()
}

func (n *NetKernel) OnUpdate(timeDelta uint32) {

}

func (n *NetKernel) OnOK() {

}

func (n *NetKernel) OnStartClose() {

}

func (n *NetKernel) DoClose() {

}

func (n *NetKernel) OnStartCheck() int {
	return 0
}

func (n *NetKernel) OnCloseCheck() int {
	return 0
}

func (n *NetKernel) GetNoWaitStart() bool {
	return n.NoWaitStart
}

func (n *NetKernel) RegisterMsg(msgId uint16, handlerFunc network.HandlerFunc) {
	n.handlers[msgId] = handlerFunc
}

func (n *NetKernel) MessageHandler(packet *network.Packet) {
	handler := n.handlers[packet.Msg.ID]
	handler(packet)
}

func (n *NetKernel) GetNPManager() network.INPManager {
	return n.NetAcceptor.NPManager
}

func (n *NetKernel) OnRpcNetAccept(args []interface{}) {

}

func (n *NetKernel) OnRpcNetConnected(args []interface{}) {

}

func (n *NetKernel) OnRpcNetError(args []interface{}) {
	np := args[0].(*network.NetPoint)
	n.GetNPManager().Del(np)
	n.Owner.OnNetError(np)
	np.Close()
}

func (n *NetKernel) OnRpcNetData(args []interface{}) {

}

func (n *NetKernel) OnRpcNetMessage(args []interface{}) {

}
