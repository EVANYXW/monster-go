package module

import (
	"fmt"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/output"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"
	"net"
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
	netType     NetType
	NPManager   network.INPManager
	msgHandler  MsgHandler
	processor   *network.Processor
	NetAcceptor *network.Acceptor
	RpcAcceptor *rpc.Acceptor
	Status      int
	//handlers    network.HandlerMap
	closeChan   chan struct{}
	port        uint32
	netMaxCount uint32
	NoWaitStart bool
	packer      network.Packer
}

func NewNetKernel(maxConnNum uint32, info server.Info, msgHandler MsgHandler, packer network.Packer, options ...kernelOption) *NetKernel {
	rpcAcceptor := rpc.NewAcceptor(10000)
	processor := network.NewProcessor()
	nodePointManager := network.NewNormal(maxConnNum, rpcAcceptor, processor, packer)

	//if network.GlobalProcess == nil {
	//	network.GlobalProcess = network.NewProcessor()
	//}

	kernel := &NetKernel{
		NPManager:   nodePointManager,
		NetAcceptor: network.NewAcceptor(maxConnNum, info, nodePointManager),
		RpcAcceptor: rpcAcceptor,
		processor:   processor,
		closeChan:   make(chan struct{}),
		NoWaitStart: false,
		netType:     Inner, // 默认内网
		msgHandler:  msgHandler,
	}

	for _, fn := range options {
		fn(kernel)
	}

	kernel.NetAcceptor.MessageHandler = kernel.MessageHandler
	kernel.Init()
	return kernel
}

func (n *NetKernel) Init() bool {
	AddManager(ModuleID_SM, n.NPManager)
	return true
}

func (n *NetKernel) DoRegist() {
	n.RpcAcceptor.Regist(rpc.RPC_NET_ACCEPT, n.OnRpcNetAccept)
	n.RpcAcceptor.Regist(rpc.RPC_NET_CONNECTED, n.OnRpcNetConnected)
	n.RpcAcceptor.Regist(rpc.RPC_NET_ERROR, n.OnRpcNetError)

	if n.msgHandler != nil {
		n.msgHandler.MsgRegister(n.processor)
	}

}
func (n *NetKernel) start(options ...network.Options) {
	async.Go(func() {
		n.NetAcceptor.Connect(options...)
		n.Status = server.Net_RunStep_Done
		n.NetAcceptor.Run() // 会阻塞
	})
	n.msgHandler.Start()
}

func (n *NetKernel) Start() {
	n.Status = server.Net_RunStep_Start
	if n.NoWaitStart {
		n.start()
	}
}

func (n *NetKernel) DoStart() {
	port := server.Ports[server.EP_Client]
	if n.netType == Inner {
		port = server.Ports[server.EP_Gate]
	}
	addr := fmt.Sprintf(":%d", port)
	output.Oput.SetServerAddr(addr)
	n.start(network.WithAddr(addr))
	//async.Go(func() {
	//	n.NetAcceptor.Connect(network.WithAddr(addr))
	//	n.NetAcceptor.Run()
	//	n.status = server.CN_RunStep_Done
	//})
	//n.msgHandler.Start()

	//n.NetAcceptor.DoStart()
}

func (n *NetKernel) Release() {
	n.NetAcceptor.OnClose()
}

func (n *NetKernel) OnUpdate(timeDelta uint32) {

}

func (n *NetKernel) OnOK() {
	n.msgHandler.OnOk()
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

//func (n *NetKernel) RegisterMsg(msgId uint16, handlerFunc network.HandlerFunc) {
//	//n.handlers[msgId] = handlerFunc
//	n.processor.RegisterMsg(msgId, handlerFunc)
//}

func (n *NetKernel) MessageHandler(packet *network.Packet) {
	//if n.msgHandler != nil && n.msgHandler.GetIsHandle() {
	//	n.msgHandler.OnNetMessage(packet)
	//	return
	//}

	//n.processor.MessageHandler(packet)
	packet.NetPoint.Processor.MessageHandler(packet)
}

func (n *NetKernel) GetNPManager() network.INPManager {
	return n.NetAcceptor.NPManager
}

func (n *NetKernel) OnRpcNetAccept(args []interface{}) {
	np := args[0].(*network.NetPoint)
	acc := args[1].(*network.Acceptor)
	fmt.Println("OnRpcNetAccept ....")
	n.msgHandler.OnRpcNetAccept(np)
	conn := np.Conn.(*net.TCPConn)
	acc.RemoveConn(conn, np)
}

func (n *NetKernel) OnRpcNetConnected(args []interface{}) {
	np := args[0].(*network.NetPoint)
	n.msgHandler.OnNetConnected(np)
}

func (n *NetKernel) OnRpcNetError(args []interface{}) {
	fmt.Println("OnRpcNetError !!!")
	np := args[0].(*network.NetPoint)
	async.Go(func() {
		np.CloseChan <- true
	})
	close(np.CloseChan)
	n.NPManager.Del(np)
	//n.Owner.OnNetError(np)
	n.msgHandler.OnNetError(np)
	fmt.Println("NetKernel OnRpcNetError np close")
	np.Close()
}

func (n *NetKernel) OnRpcNetData(args []interface{}) {

}

func (n *NetKernel) OnRpcNetMessage(args []interface{}) {
	//np := args[0].(*network.NetPoint)
	//message := args[1].(*network.Message)
}
