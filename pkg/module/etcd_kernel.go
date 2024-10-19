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
	"strings"
	"time"
)

type EtcdKernel struct {
	RpcAcceptor *rpc.Acceptor
	ID          uint32
	SID         server.ServerID
	address     string
	runStatus   int
	NoWaitStart bool
	servername  string
	logger      *zap.Logger
	netType     NetType
	isWatch     bool
	etcdClient  *clientv3.Client
	etcdServers map[string]string
}

const (
	serviceKeyPrefix = "monster-go"
	leaseTTL         = 5 // 租约的生存时间（秒）
)

var (
	GlobalEtcdKernel *EtcdKernel
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
	GlobalEtcdKernel = connector
	return connector
}

func (e *EtcdKernel) SetID(id uint32) {
	e.ID = id
	server.ID2Sid(id, &e.SID)
}

func (e *EtcdKernel) Init(baseModule IBaseModule) bool {
	e.runStatus = ModuleRunStatus_Start
	return true
}

func (e *EtcdKernel) DoRegister() {

}

func (e *EtcdKernel) DoRun() {
	e.RpcAcceptor.Run()
	e.runStatus = ModuleRunStatus_Running

	outPort := configs.All().OutPort

	ip, port, err := ipPort.GetDynamicIpAndRangePort(int(outPort), int(outPort)+100)
	if err != nil {
		e.logger.Error("ipPort GetDynamicIpAndRangePort is failed", zap.Error(err))
		return
	}
	e.address = fmt.Sprintf("%s:%d", ip, port)

	if e.netType == Outer {
		server.Ports[server.EP_Client] = uint32(port)
	} else {
		server.Ports[server.EP_Gate] = uint32(port)
	}

	//e.RegisterService(e.servername, addr) // net 准备好了才注册etcd
	DoWaitStart()

	if e.isWatch {
		go e.watchService()
		e.GetServices(e.etcdClient, serviceKeyPrefix)
	}
}

func (e *EtcdKernel) DoWaitStart() {

}

func (e *EtcdKernel) DoRelease() {
}

func (e *EtcdKernel) Update() {

}

func (e *EtcdKernel) OnOk() {
	//e.msgHandler.OnOk()
}

func (e *EtcdKernel) OnStartClose() {

}

func (e *EtcdKernel) DoClose() {

}

func (e *EtcdKernel) OnStartCheck() int {
	if e.runStatus == ModuleRunStatus_Running {
		return ModuleOk()
	}
	return ModuleWait()
}

func (e *EtcdKernel) GetNoWaitStart() bool {
	return e.NoWaitStart
}

func (e *EtcdKernel) OnCloseCheck() int {
	return ModuleOk()
}

func (e *EtcdKernel) GetNPManager() network.INPManager {
	return nil
}

func (e *EtcdKernel) GetStatus() int {
	return 0
}

func (e *EtcdKernel) OnRpcNetAccept(args []interface{}) {

}

func (e *EtcdKernel) OnRpcNetConnected(args []interface{}) {

}

func (e *EtcdKernel) OnRpcNetError(args []interface{}) {

}

func (e *EtcdKernel) OnRpcNetClose(args []interface{}) {

}

func (e *EtcdKernel) OnRpcNetData(args []interface{}) {

}

func (e *EtcdKernel) OnRpcNetMessage(args []interface{}) {

}

// RegisterService
// @Description 注册发现负债均衡
// @Author evan_yxw 2024-10-17 12:11:07
// @Param servername
// @Param serviceAddr
func (e *EtcdKernel) RegisterService() {
	leaseId, err := registerService(e.etcdClient, e.servername, e.address)
	if err != nil {
		e.logger.Error("Failed to register service", zap.Error(err))
		return
	}
	e.logger.Info("Registered service leaseId", zap.Int("leaseId", int(leaseId)))
}

func (e *EtcdKernel) DiscoverServices(servername string) string {
	services, err := discoverServices(e.etcdClient, servername)
	if err != nil {
		e.logger.Error("Failed to discover services", zap.Error(err))
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

	// 生成唯一的服务名称 monster-go:login:37439:6784dcaf-0be5-4433-a315-9271e69b06df
	uniqueServiceName := fmt.Sprintf("%s:%s:%d:%s", serviceKeyPrefix, serviceName, server.ID, uuid.New().String())

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

func (e *EtcdKernel) GetServices(etcdClient *clientv3.Client, serverKey string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 从 etcd 获取所有服务实例，使用 WithPrefix 进行前缀匹配
	resp, err := etcdClient.Get(ctx, serverKey, clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("failed to get services: %v", err)
	}

	for _, service := range resp.Kvs {
		fmt.Println(string(service.Key))
		e.connect(string(service.Key), string(service.Value))
	}

	return nil
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

func (e *EtcdKernel) getServeName(etcdKey string) string {
	splitArr := strings.Split(etcdKey, ":")
	if len(splitArr) != 4 {
		return ""
	}

	return splitArr[1]
}

func (e *EtcdKernel) getServerId(etcdKey string) uint32 {
	splitArr := strings.Split(etcdKey, ":")

	var serverId uint32
	if len(splitArr) > 3 {
		id, _ := strconv.ParseInt(splitArr[2], 10, 64)
		serverId = uint32(id)
	}

	return serverId
}

func (e *EtcdKernel) getIpPort(etcdValue string) (string, uint32) {
	ip, port, err := net.SplitHostPort(etcdValue)
	if err != nil {
		fmt.Println("Error:", err)
		return "", 0
	}
	portInt, _ := strconv.ParseInt(port, 10, 64)
	return ip, uint32(portInt)
}

func (e *EtcdKernel) connect(etcdKey, etcdValue string) bool {
	if e.servername == e.getServeName(etcdKey) {
		return false
	}
	ip, port := e.getIpPort(etcdValue)
	serverId := e.getServerId(etcdKey)
	owner := GetModule(ModuleID_ConnectorManager).GetOwner()
	connectorManager := owner.(*Manager)

	time.Sleep(100 * time.Millisecond) // 延时,对于打印的conn num
	conn := connectorManager.CreateConnector(serverId, ip, port)
	if conn == nil {

	}

	return true
}

// watchService 监听服务的变化
func (e *EtcdKernel) watchService() {
	//watchChan := e.etcdClient.Watch(context.Background(), serviceKey)
	watchChan := e.etcdClient.Watch(context.Background(), serviceKeyPrefix, clientv3.WithPrefix())
	fmt.Println("Watching for changes...")
	owner := GetModule(ModuleID_ConnectorManager).GetOwner()
	connectorManager := owner.(*Manager)
	//connectorManager := owner.(tcp_manager.TcpConnectorManager)
	for watchResp := range watchChan {
		for _, ev := range watchResp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				fmt.Printf("Service updated: %s : %s\n", ev.Kv.Key, ev.Kv.Value)
				e.connect(string(ev.Kv.Key), string(ev.Kv.Value))
			case clientv3.EventTypeDelete:
				fmt.Printf("Service deleted: %s\n", ev.Kv.Key)
				serverId := e.getServerId(string(ev.Kv.Key))
				connectorManager.DelConnector(serverId)
			}
		}
	}
}
