package module

import (
	"context"
	"fmt"
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/pkg/ipPort"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"
	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

type EtcdKernel struct {
	RpcAcceptor *rpc.Acceptor
	ID          uint32
	SID         server.ServerID
	wg          sync.WaitGroup
	runStatus   int
	NoWaitStart bool
	etcdClient  *clientv3.Client
	servername  string
	logger      *zap.Logger
	netType     NetType
	isWatch     bool
}

const (
	serverKey        = "monster-go"
	serviceKeyPrefix = "monster-go"
	leaseTTL         = 5 // 租约的生存时间（秒）
)

func NewEtcdKernel(servername string, isWatch bool, netType NetType, etcdClient *clientv3.Client, logger *zap.Logger, options ...ckernelOption) *EtcdKernel {
	opt := &ckOptions{}

	rpcAcceptor := rpc.NewAcceptor(10000)
	connector := &EtcdKernel{
		RpcAcceptor: rpcAcceptor,
		NoWaitStart: false,
		etcdClient:  etcdClient,
		logger:      logger,
		servername:  servername,
		netType:     netType, // 默认内网
		isWatch:     isWatch,
	}

	for _, fn := range options {
		fn(opt)
	}
	connector.NoWaitStart = opt.NoWaitStart
	return connector
}

func (c *EtcdKernel) SetID(id uint32) {
	c.ID = id
	server.ID2Sid(id, &c.SID)
}

func (c *EtcdKernel) Init(baseModule *BaseModule) bool {
	c.runStatus = ModuleRunStatus_Start
	return true
}

func (c *EtcdKernel) DoRegister() {
	c.RpcAcceptor.Regist(rpc.RPC_NET_CONNECTED, c.OnRpcNetConnected)
	c.RpcAcceptor.Regist(rpc.RPC_NET_ERROR, c.OnRpcNetError)

	//if c.msgHandler != nil {
	//	c.msgHandler.MsgRegister(c.processor)
	//}
}

func (c *EtcdKernel) DoRun() {
	c.RpcAcceptor.Run()
	//c.Client.Run(c.RpcAcceptor)
	c.runStatus = ModuleRunStatus_Running

	outPort := configs.All().OutPort

	ip, port, err := ipPort.GetDynamicIpAndRangePort(int(outPort), int(outPort)+100)
	if err != nil {
		c.logger.Error("ipPort GetDynamicIpAndRangePort is failed", zap.Error(err))
		return
	}
	addr := fmt.Sprintf("%s:%d", ip, port)
	c.RegisterService(c.servername, addr)
	if c.isWatch {
		go c.watchService()
	}

	if c.netType == Outer {
		server.Ports[server.EP_Client] = uint32(port)
	} else {
		server.Ports[server.EP_Gate] = uint32(port)
	}

	DoWaitStart()
	//c.msgHandler.Start()
}

func (c *EtcdKernel) DoWaitStart() {

}

func (c *EtcdKernel) DoRelease() {
	//c.Client.OnClose()
}

func (c *EtcdKernel) Update() {

}

func (c *EtcdKernel) OnOk() {
	//c.msgHandler.OnOk()
}

func (c *EtcdKernel) OnStartClose() {

}

func (c *EtcdKernel) DoClose() {

}

func (c *EtcdKernel) OnStartCheck() int {
	if c.runStatus == ModuleRunStatus_Running {
		return ModuleRunCode_Ok
	}
	return ModuleRunCode_Wait
}

func (c *EtcdKernel) GetNoWaitStart() bool {
	return c.NoWaitStart
}

func (c *EtcdKernel) OnCloseCheck() int {
	return 0
}

func (c *EtcdKernel) GetNPManager() network.INPManager {
	return nil
}

func (c *EtcdKernel) GetStatus() int {
	return 0
}

func (c *EtcdKernel) RegisterMsg(msgId uint16, handlerFunc network.HandlerFunc) {
	//c.processor.RegisterMsg(msgId, handlerFunc)
}

func (c *EtcdKernel) MessageHandler(packet *network.Packet) {
	//c.processor.MessageHandler(packet)
}

func (c *EtcdKernel) OnRpcNetAccept(args []interface{}) {

}

func (c *EtcdKernel) OnRpcNetConnected(args []interface{}) {
	//np := args[0].(*network.NetPoint)
	//c.msgHandler.OnNetConnected(np)
}

func (c *EtcdKernel) OnRpcNetError(args []interface{}) {
	//np := args[0].(*network.NetPoint)
	//
	//if c.msgHandler != nil {
	//	c.msgHandler.OnNetError(np, nil)
	//}
	//fmt.Println("EtcdKernel OnRpcNetError np close")
	//np.Close()

	//fixMe OnRpcNetError 还没做其他处理!!!
	fmt.Println("OnRpcNetError 还没做其他处理!!!")
}

func (c *EtcdKernel) OnRpcNetClose(args []interface{}) {

}

func (c *EtcdKernel) OnRpcNetData(args []interface{}) {

}

func (c *EtcdKernel) OnRpcNetMessage(args []interface{}) {

}

func (c *EtcdKernel) RegisterService(servername, serviceAddr string) {
	leaseId, err := registerService(c.etcdClient, servername, serviceAddr)
	if err != nil {
		c.logger.Error("Failed to register service", zap.Error(err))
		return
	}
	c.logger.Info("Registered service leaseId", zap.Int("leaseId", int(leaseId)))
}

func (c *EtcdKernel) DiscoverServices(servername string) string {
	services, err := discoverServices(c.etcdClient, servername)
	if err != nil {
		c.logger.Error("Failed to discover services", zap.Error(err))
		return ""
	}
	serverAddr := balance(services)
	return serverAddr
}

func registerService(etcdClient *clientv3.Client, serviceName, serviceAddr string) (clientv3.LeaseID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 创建租约
	resp, err := etcdClient.Grant(ctx, leaseTTL)
	if err != nil {
		return 0, fmt.Errorf("failed to create lease: %v", err)
	}

	// 生成唯一的服务名称
	uniqueServiceName := fmt.Sprintf("%s:%s-%s", serverKey, serviceName, uuid.New().String())

	// 将服务地址注册到 etcd，并与租约绑定
	_, err = etcdClient.Put(ctx, uniqueServiceName, serviceAddr, clientv3.WithLease(resp.ID))
	if err != nil {
		return 0, fmt.Errorf("failed to register service: %v", err)
	}
	fmt.Printf("Service %s registered with address %s\n", uniqueServiceName, serviceAddr)

	// 启动一个 goroutine 定期续租
	go keepAlive(etcdClient, resp.ID)

	return resp.ID, nil
}

func discoverServices(etcdClient *clientv3.Client, serviceName string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 从 etcd 获取所有服务实例，使用 WithPrefix 进行前缀匹配
	resp, err := etcdClient.Get(ctx, serviceName, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %v", err)
	}

	var services []string
	for _, kv := range resp.Kvs {
		services = append(services, string(kv.Value))
	}

	return services, nil
}

func balance(services []string) string {
	// 随机选择一个服务实例
	var selectedService string
	if len(services) > 0 {
		selectedService = services[rand.Intn(len(services))]
		fmt.Printf("Selected service address: %s\n", selectedService)
	} else {
		fmt.Println("No services found")
	}

	return selectedService
}

func keepAlive(etcdClient *clientv3.Client, leaseID clientv3.LeaseID) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		_, err := etcdClient.KeepAliveOnce(ctx, leaseID)
		if err != nil {
			log.Printf("Failed to keep alive for lease %v: %v", leaseID, err)
			return
		}
		time.Sleep(leaseTTL / 2 * time.Second) // 每半个TTL时间续租一次
	}
}

// watchService 监听服务的变化
func (c *EtcdKernel) watchService() {
	//watchChan := c.etcdClient.Watch(context.Background(), serviceKey)
	watchChan := c.etcdClient.Watch(context.Background(), serviceKeyPrefix, clientv3.WithPrefix())
	fmt.Println("Watching for changes...")
	owner := GetModule(ModuleID_ConnectorManager).GetOwner()
	connectorManager := owner.(*Manager)
	//connectorManager := owner.(tcp_manager.TcpConnectorManager)
	for watchResp := range watchChan {
		for _, ev := range watchResp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				fmt.Printf("Service updated: %s : %s\n", ev.Kv.Key, ev.Kv.Value)
				ip, port, err := net.SplitHostPort(string(ev.Kv.Value))
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				portInt, _ := strconv.ParseInt(port, 10, 64)
				time.Sleep(time.Second * 10)
				conn := connectorManager.CreateConnector(0, ip, uint32(portInt))
				if conn == nil {

				}
			case clientv3.EventTypeDelete:
				fmt.Printf("Service deleted: %s\n", ev.Kv.Key)
			}
		}
	}
}
