package module

import (
	"fmt"
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/pkg/grpcpool"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/network"
	//"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"
	"sync"
)

type GrpcConnectorKernel struct {
	//*client.Client
	*grpcpool.Connector

	//RpcAcceptor *rpc.Acceptor
	ID          uint32
	SID         server.ServerID
	wg          sync.WaitGroup
	runStatus   int
	NoWaitStart bool
	msgHandler  MsgHandler
	//processor   *network.Processor
}

func NewGrpcConnectorKernel(servername string) *GrpcConnectorKernel {
	//opt := &ckOptions{}
	//rpcAcceptor := rpc.NewAcceptor(10000)
	//processor := network.NewProcessor()
	etcdCnf := configs.All().Etcd
	etcdClient := grpcpool.InitEtcd(etcdCnf.Addr, etcdCnf.User, etcdCnf.Pass)
	connector := &GrpcConnectorKernel{
		//processor:   processor,
		//Client: client.NewClient(fmt.Sprintf("%s:%d", ip, port), processor, packerFactory),
		//RpcAcceptor: rpcAcceptor,
		Connector:   grpcpool.NewConnector(servername, etcdClient, logger.GetLogger()),
		NoWaitStart: false,
		//msgHandler:  msgHandler,
	}
	//connector.Client.OnMessageCb = connector.MessageHandler

	//for _, fn := range options {
	//	fn(opt)
	//}
	//connector.NoWaitStart = opt.NoWaitStart
	return connector
}

func (c *GrpcConnectorKernel) SetID(id uint32) {
	c.ID = id
	server.ID2Sid(id, &c.SID)
}

func (c *GrpcConnectorKernel) Init(baseModule *BaseModule) bool {
	c.runStatus = ModuleRunStatus_Start
	return true
}

func (c *GrpcConnectorKernel) DoRegister() {

}

func (c *GrpcConnectorKernel) DoRun() {
	//c.RpcAcceptor.Run()
	//c.Client.Run(c.RpcAcceptor)
	//c.runStatus = ModuleRunStatus_Running
	//c.msgHandler.Start()
}

func (c *GrpcConnectorKernel) DoWaitStart() {

}

func (c *GrpcConnectorKernel) DoRelease() {
	//c.Client.OnClose()
}

func (c *GrpcConnectorKernel) Update() {

}

func (c *GrpcConnectorKernel) OnOk() {
	//c.msgHandler.OnOk()
}

func (c *GrpcConnectorKernel) OnStartClose() {

}

func (c *GrpcConnectorKernel) DoClose() {

}

func (c *GrpcConnectorKernel) OnStartCheck() int {
	if c.runStatus == ModuleRunStatus_Running {
		return ModuleRunCode_Ok
	}
	return ModuleRunCode_Wait
}

func (c *GrpcConnectorKernel) GetNoWaitStart() bool {
	return c.NoWaitStart
}

func (c *GrpcConnectorKernel) OnCloseCheck() int {
	return 0
}

func (c *GrpcConnectorKernel) GetNPManager() network.INPManager {
	return nil
}

func (c *GrpcConnectorKernel) GetStatus() int {
	return 0
}

func (c *GrpcConnectorKernel) RegisterMsg(msgId uint16, handlerFunc network.HandlerFunc) {
	//c.processor.RegisterMsg(msgId, handlerFunc)
}

func (c *GrpcConnectorKernel) MessageHandler(packet *network.Packet) {
	//c.processor.MessageHandler(packet)
}

func (c *GrpcConnectorKernel) OnRpcNetAccept(args []interface{}) {

}

func (c *GrpcConnectorKernel) OnRpcNetConnected(args []interface{}) {
	np := args[0].(*network.NetPoint)
	c.msgHandler.OnNetConnected(np)
}

func (c *GrpcConnectorKernel) OnRpcNetError(args []interface{}) {
	//fixMe OnRpcNetError 还没做其他处理!!!
	fmt.Println("OnRpcNetError 还没做其他处理!!!")
}

func (c *GrpcConnectorKernel) OnRpcNetClose(args []interface{}) {

}

func (c *GrpcConnectorKernel) OnRpcNetData(args []interface{}) {

}

func (c *GrpcConnectorKernel) OnRpcNetMessage(args []interface{}) {

}
