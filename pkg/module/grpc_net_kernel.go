package module

import (
	"fmt"
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/grpcpool"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/output"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"
)

type GrpcNetKernel struct {
	//msgHandler  MsgHandler
	RpcAcceptor *rpc.Acceptor
	Status      int
	closeChan   chan struct{}
	port        uint32
	netType     NetType
	server      *grpcpool.Server
	NoWaitStart bool
	grpcservers []server.GrpcServer
}

var (
	GrpcServer *grpcpool.Server
)

func NewGrpcNetKernel(servername string, grpcservers []server.GrpcServer) *GrpcNetKernel {
	etcdCnf := configs.All().Etcd
	etcdClient := grpcpool.InitEtcd(etcdCnf.Addr, etcdCnf.User, etcdCnf.Pass)
	rpcAcceptor := rpc.NewAcceptor(10000)
	kernel := &GrpcNetKernel{
		RpcAcceptor: rpcAcceptor,
		closeChan:   make(chan struct{}),
		server: grpcpool.NewServer(servername, etcdClient, grpcpool.WithLogger(logger.GetLogger()),
			grpcpool.WithPorts(configs.All().Server.MinPort, configs.All().Server.MaxPort)),
		NoWaitStart: true,
		grpcservers: grpcservers,
		//msgHandler:  msgHandler,
	}

	GrpcServer = kernel.server
	return kernel
}

func (n *GrpcNetKernel) Init(baseModule *BaseModule) bool {
	//AddManager(ModuleID_SM, n.NPManager)
	return true
}

func (n *GrpcNetKernel) DoRegister() {
	//n.RpcAcceptor.Regist(rpc.RPC_NET_ACCEPT, n.OnRpcNetAccept)
	//n.RpcAcceptor.Regist(rpc.RPC_NET_CONNECTED, n.OnRpcNetConnected)
	//n.RpcAcceptor.Regist(rpc.RPC_NET_ERROR, n.OnRpcNetError)
	//n.RpcAcceptor.Regist(rpc.RPC_NET_CLOSE, n.OnRpcNetClose)

	//if n.msgHandler != nil {
	//	//n.msgHandler.MsgRegister(n.processor)
	//}

}
func (n *GrpcNetKernel) start() {
	async.Go(func() {
		//n.NetAcceptor.Connect(options...)
		//n.Status = server.Net_RunStep_Done
		server.StatusDone(&n.Status)
		n.server.Connect()
		output.Oput.SetServerAddr(n.server.GetAddr())
		for _, s := range n.grpcservers {
			s.TransportRegister()(n.server)
		}
		n.server.Run()
		//n.RpcAcceptor.Run()
		//n.NetAcceptor.Run() // 会阻塞
	})
	//// todo: m.msgHandler 需要判断然后
	//n.msgHandler.Start()
}

func (n *GrpcNetKernel) DoRun() {
	//n.nodeManager.Start()
	//n.Status = server.Net_RunStep_Start
	server.StatusStart(&n.Status)
	if n.NoWaitStart {
		n.start()
	}
}

func (n *GrpcNetKernel) DoWaitStart() {
	//port := server.Ports[server.EP_Client]
	//if n.netType == Inner {
	//	port = server.Ports[server.EP_Gate]
	//}
	//addr := fmt.Sprintf(":%d", 30009)
	//output.Oput.SetServerAddr(n.server.GetAddr())
	n.start()
	//async.Go(func() {
	//	n.NetAcceptor.Connect(network.WithAddr(addr))
	//	n.NetAcceptor.Run()
	//	n.status = server.CN_RunStep_Done
	//})
	//n.msgHandler.Start()

	//n.NetAcceptor.DoStart()
}

func (n *GrpcNetKernel) DoRelease() {
	//n.NetAcceptor.OnClose()
}

func (n *GrpcNetKernel) Update() {

}

func (n *GrpcNetKernel) OnOk() {
	//n.msgHandler.OnOk()
}

func (n *GrpcNetKernel) OnStartClose() {

}

func (n *GrpcNetKernel) DoClose() {

}

func (n *GrpcNetKernel) OnStartCheck() int {
	return 0
}

func (n *GrpcNetKernel) OnCloseCheck() int {
	return 0
}

func (n *GrpcNetKernel) GetNoWaitStart() bool {
	return n.NoWaitStart
}

//func (n *GrpcNetKernel) RegisterMsg(msgId uint16, handlerFunc network.HandlerFunc) {
//	//n.handlers[msgId] = handlerFunc
//	n.processor.RegisterMsg(msgId, handlerFunc)
//}

func (n *GrpcNetKernel) MessageHandler(packet *network.Packet) {
	//if n.msgHandler != nil && n.msgHandler.GetIsHandle() {
	//	n.msgHandler.OnNetMessage(packet)
	//	return
	//}

	//n.processor.MessageHandler(packet)
	packet.NetPoint.Processor.MessageHandler(packet)
}

func (n *GrpcNetKernel) GetNPManager() network.INPManager {
	//return n.NetAcceptor.NPManager
	return nil
}

func (n *GrpcNetKernel) GetStatus() int {
	return n.Status
}

func (n *GrpcNetKernel) OnRpcNetAccept(args []interface{}) {
	//np := args[0].(*network.NetPoint)
	//acc := args[1].(*network.Acceptor)
	//fmt.Println("OnRpcNetAccept ....")
	//n.msgHandler.OnRpcNetAccept(np, acc)
}

func (n *GrpcNetKernel) OnRpcNetConnected(args []interface{}) {
	//np := args[0].(*network.NetPoint)
	//n.msgHandler.OnNetConnected(np)
}

func (n *GrpcNetKernel) OnRpcNetError(args []interface{}) {
	//fixMe OnRpcNetError 还没做其他处理!!!
	fmt.Println("OnRpcNetError 还没做其他处理!!!")
	//np := args[0].(*network.NetPoint)
	//acc := args[1].(*network.Acceptor)

	//n.NPManager.Del(np)
	//n.msgHandler.OnNetError(np, n.NetAcceptor)
	//fmt.Println("GrpcNetKernel OnRpcNetError np close")
	//np.Close()
}

func (n *GrpcNetKernel) OnRpcNetClose(args []interface{}) {
	//fmt.Println("OnRpcNetClose !!!")
	//np := args[0].(*network.NetPoint)
	//
	//n.NPManager.Del(np)
	//n.msgHandler.OnNetError(np, n.NetAcceptor)
	//fmt.Println("GrpcNetKernel OnRpcNetError np close")
	//np.Close()
}

func (n *GrpcNetKernel) OnRpcNetData(args []interface{}) {

}

func (n *GrpcNetKernel) OnRpcNetMessage(args []interface{}) {
	//np := args[0].(*network.NetPoint)
	//message := args[1].(*network.Message)
}
