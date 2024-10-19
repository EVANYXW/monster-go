package kernel

import (
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/grpcpool"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module/module_def"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/output"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"
)

type GrpcNetKernel struct {
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
	}

	GrpcServer = kernel.server
	return kernel
}

func (n *GrpcNetKernel) Init(baseModule module_def.IBaseModule) bool {
	return true
}

func (n *GrpcNetKernel) DoRegister() {

}

func (n *GrpcNetKernel) start() {
	async.Go(func() {
		server.NetStatusDone(&n.Status)
		n.server.Connect()
		output.Oput.SetServerAddr(n.server.GetAddr())
		for _, s := range n.grpcservers {
			s.TransportRegister()(n.server)
		}
		n.server.Run()
	})
}

func (n *GrpcNetKernel) DoRun() {
	server.NetStatusStart(&n.Status)
	if n.NoWaitStart {
		n.start()
	}
}

func (n *GrpcNetKernel) DoWaitStart() {
	n.start()
}

func (n *GrpcNetKernel) DoRelease() {

}

func (n *GrpcNetKernel) Update() {

}

func (n *GrpcNetKernel) OnOk() {
}

func (n *GrpcNetKernel) OnStartClose() {

}

func (n *GrpcNetKernel) DoClose() {

}

func (n *GrpcNetKernel) OnStartCheck() int {
	return module_def.ModuleOk()
}

func (n *GrpcNetKernel) OnCloseCheck() int {
	return module_def.ModuleOk()
}

func (n *GrpcNetKernel) GetNoWaitStart() bool {
	return n.NoWaitStart
}

func (n *GrpcNetKernel) MessageHandler(packet *network.Packet) {

}

func (n *GrpcNetKernel) GetNPManager() network.INPManager {
	//return n.NetAcceptor.NPManager
	return nil
}

func (n *GrpcNetKernel) GetStatus() int {
	return n.Status
}

func (n *GrpcNetKernel) OnRpcNetAccept(args []interface{}) {

}

func (n *GrpcNetKernel) OnRpcNetConnected(args []interface{}) {

}

func (n *GrpcNetKernel) OnRpcNetError(args []interface{}) {

}

func (n *GrpcNetKernel) OnRpcNetClose(args []interface{}) {

}

func (n *GrpcNetKernel) OnRpcNetData(args []interface{}) {

}

func (n *GrpcNetKernel) OnRpcNetMessage(args []interface{}) {

}
