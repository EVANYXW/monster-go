package engine

import (
	"fmt"
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/internal/servers/gate/handler"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/env"
	"github.com/evanyxw/monster-go/pkg/grpcpool"
	handler2 "github.com/evanyxw/monster-go/pkg/handler"
	"github.com/evanyxw/monster-go/pkg/kernel"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/logs"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/module/module_def"
	register_discovery "github.com/evanyxw/monster-go/pkg/module/register-discovery"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/output"
	"github.com/evanyxw/monster-go/pkg/server"
	"github.com/evanyxw/monster-go/pkg/timeutil"
	"go.uber.org/zap"
	"log"
	"os"
)

type clientType int

const (
	grpcClientType clientType = iota
	httpClientType
)

type BaseEngine struct {
	servername string
	isPprof    bool
	isOutput   bool
	output     *output.Output
}

type Options func(opt *option)

type option struct {
	isPprof    bool
	isOutput   bool
	output     *output.Config
	modules    map[int32]module_def.IModule
	clientType clientType
}

func WithPprof() Options {
	return func(opt *option) {
		opt.isPprof = true
	}
}

func WithOutput(config *output.Config) Options {
	return func(opt *option) {
		opt.isOutput = true
		opt.output = config
	}
}

func WithModule(m module_def.IModule) Options {
	return func(opt *option) {
		module_def.NewBaseModule(m.GetID(), m)
	}
}

func WithModules(modules map[int32]module_def.IModule) Options {
	return func(opt *option) {
		for id, m := range opt.modules {
			opt.modules[id] = m
		}
	}
}

func WithGrpcClient() Options {
	return func(opt *option) {
		opt.clientType = grpcClientType
	}
}

func (b *BaseEngine) WithModule(m module_def.IModule) *BaseEngine {
	module_def.NewBaseModule(m.GetID(), m)
	return b
}

func (b *BaseEngine) WithOutput(config *output.Config) *BaseEngine {
	//b.output = output.NewOutput(config)
	b.isOutput = true
	b.output = output.NewOutput(config, module_def.GetModuleMap())
	async.Go(func() {
		b.output.Run()
	})
	return b
}

// initLog init log  etcd 在用
func initLog(servername string) {
	logs.NewLogger(
		logs.WithFilePath(fmt.Sprintf("log/%s.log", servername)),
		logs.WithCompress(false),
		logs.WithPrettyPrint(false),
		logs.WithFormat("json"),
		logs.WithLevel(5),
		logs.WithMaxSize(100),
		logs.WithServerName(servername),
	)
}

func Init() {
	{
		_, _ = logger.NewJSONLogger(
			logger.WithField("domain", fmt.Sprintf("%s[%s]", configs.ProjectName, env.Active().Value())),
			logger.WithTimeLayout(timeutil.CSTLayout),
			logger.WithFileP(configs.LogFile, server.GetServerInfo().ServerName),
		)

		server.SID.Type = server.Name2EP(server.GetServerInfo().ServerName)
		server.SID.ID = 1
		if server.Name2EP(server.GetServerInfo().ServerName) == server.EP_Center {
			server.SID.Index = 1
		}
		logger.Info("SID:", zap.Uint16("ID", server.SID.ID),
			zap.Uint8("type", server.SID.Type), zap.Uint8("index", server.SID.Index))

		server.UpdateID()

		module_def.Init()

		initLog(server.GetServerInfo().ServerName) // 主要etcd 包里在用

		etcdCnf := configs.All().Etcd
		etcdClient := grpcpool.InitEtcd(etcdCnf.Addr, etcdCnf.User, etcdCnf.Pass)
		grpcpool.NewConnector(server.GetServerInfo().ServerName, etcdClient, logger.GetLogger())
		//if err := plugins.InitPlugins(context.Background()); err != nil {
		//	panic(err)
		//}

		// mysql、redis
		//newModel()
		// redis 发布订阅
		//redisSub(recvPublish)

		// 初始化etcd
		//etcd := initEtcd()
		//
		//// tcp 服务注册etcd
		//tcpEtcd := registerEtcd(etcd, serverInfo.ServerName, serverInfo.Address)
		//etcdServerArr.add(tcpEtcd)
		//
		//// rpc 服务注册etcd
		//rpcEtcd := registerEtcd(etcd, fmt.Sprintf("%s%s", servername, servers.Rpc), serverInfo.RpcAddr)
		//etcdServerArr.add(rpcEtcd)

		// 开启pprof
		//pprofUrl := ""
		//if allConfig.Center.Pprof {
		//	//pprofUrl = allConfig.Center.PprofAddress
		//	pprofPort := allConfig.PprofPort + int64(server.Name2EP(serverInfo.ServerName))
		//	pprofUrl = fmt.Sprintf(":%d", pprofPort)
		//
		//	async.Go(func() {
		//		pprofs(pprofUrl)
		//	})
		//}
		//async.Go(func() {
		//	output.NewOutput(&output.Config{
		//		Name: server.GetServerInfo().ServerName,
		//		Addr: server.GetServerInfo().Address,
		//		Url:  "",
		//	}).Run()
		//})

		//async.Go(func() {
		//	output.NewOutput(serverInfo.ServerName, serverInfo.Address,
		//		pprofUrl).Run()
		//})
	}
}

func ServerInit(name string) {
	server.SetServerInfo(&server.Info{
		ServerName: name,
		Env:        env.Active().Value(),
	})

	configs.Init()
	Init()
}

func newServer(name string, options ...Options) *BaseEngine {

	opt := &option{
		modules: make(map[int32]module_def.IModule),
	}
	for _, fn := range options {
		fn(opt)
	}

	// 装载注册发现模块
	//opt.modules[module.ModuleID_SM] = registerDiscovery

	b := &BaseEngine{
		servername: name,
		isOutput:   opt.isOutput,
		isPprof:    opt.isPprof,
	}

	if opt.isOutput {
		b.output = output.NewOutput(opt.output, module_def.GetModuleMap())
		async.Go(func() {
			b.output.Run()
		})
	}

	return b
}

// NewGateTcpServer
// @Description 创建一个Gate的tcp服务
// @Author evan_yxw 2024-10-07 21:51:23
// @Param name
// @Param factor
// @Param options
// @Return *BaseEngine
func NewGateTcpServer(name string, factor register_discovery.ConnectorFactory, options ...Options) *BaseEngine {
	if factor == nil {
		log.Fatal("Please provide a factor!")
	}

	ServerInit(name)
	factor.SetGateWay()

	// 注册与发现模块,支持内部的center和etcd两种模式
	registerDiscovery := factor.CreateConnector(
		register_discovery.WithServername(name),
		register_discovery.WithWatch(),
		register_discovery.WithNetType(kernel.Outer))

	// 创建tcp网络模块
	tcpNet := module.NewClientNet(
		module_def.ModuleID_Client,
		5000,
		handler.NewGateMsg(),
		kernel.Outer,
		new(network.DefaultPackerFactory),
	)

	//options = append(options, WithModule(registerDiscovery), WithModule(tcpNet))
	options = append(options, WithModule(tcpNet))

	// Center 模式的注册发现,需要加装 connector manager
	if factor.GetType() == register_discovery.TypeCenter {
		//options = append(options, WithModule(factor.CreateConnectorManager(&connector.TcpManagerFactory{})))
		m := factor.CreateConnectorManager()
		iModule := m.(module_def.IModule)
		options = append(options, WithModule(iModule))
		//options = append(options, WithModule(connector.NewManager(module.ModuleID_ConnectorManager)))
	}

	// Etcd 模式的注册发现,需要加载 connector manager
	if factor.GetType() == register_discovery.TypeEtcd {
		// 多个etcd的服务器，gateway需要主动连接
		m := factor.CreateConnectorManager()
		iModule := m.(module_def.IModule)
		options = append(options, WithModule(iModule))
	}

	options = append(options, WithModule(registerDiscovery))
	return newServer(name, options...)
}

// NewTcpServer
// @Description 创建一个tcp服务
// @Author evan_yxw 2024-10-07 21:51:08
// @Param name
// @Param msgHandler
// @Param factor
// @Param options
// @Return *BaseEngine
func NewTcpServer(name string, msgHandler handler2.MsgHandler, factor register_discovery.ConnectorFactory, options ...Options) *BaseEngine {
	if factor == nil {
		log.Fatal("Please provide a factor!")
	}

	ServerInit(name)

	// // 注册与发现模块,支持内部的center和etcd两种模式
	registerDiscovery := factor.CreateConnector(
		register_discovery.WithServername(name),
		register_discovery.WithNetType(kernel.Inner),
	)

	// 创建tcp网络模块
	clientNet := module.NewClientNet(
		module_def.ModuleID_GateAcceptor,
		10000,
		msgHandler,
		kernel.Inner,
		new(network.ClientPackerFactory),
	)
	options = append(options, WithModule(registerDiscovery), WithModule(clientNet))

	return newServer(name, options...)
}

// NewGrpcServer
// @Description 创建一个grpc服务
// @Author evan_yxw 2024-10-07 21:50:51
// @Param name
// @Param grpcServers
// @Param options
// @Return *BaseEngine
func NewGrpcServer(name string, grpcServers []server.GrpcServer, options ...Options) *BaseEngine {
	//factor := center.NewFactor(center.WithServerConnectorManager())
	//if factor == nil {
	//	log.Fatal("Please provide a factor!")
	//}

	ServerInit(name)

	// 注册与发现
	//rd := factor.CreateConnector(name)
	// 网络模块
	clientNet := module.NewGrpcNet(
		module_def.ModuleID_GateAcceptor,
		name,
		grpcServers,
	)

	//options = append(options, WithModule(rd), WithModule(clientNet))
	options = append(options, WithModule(clientNet))

	return newServer(name, options...)
}

// NewCenterServer
// @Description   new一个center
// @Author evan_yxw 2024-10-07 21:50:26
// @Param name
// @Param factor
// @Param options
// @Return *BaseEngine
func NewCenterServer(name string, factor register_discovery.NetFactory, options ...Options) *BaseEngine {
	if factor == nil {
		log.Fatal("Please provide a factor!")
	}

	ServerInit(name)
	rd := factor.CreateNet()
	options = append(options, WithModule(rd))

	return newServer(name, options...)
}

func (b *BaseEngine) Run() {
	module_def.Run()
}

func (b *BaseEngine) Destroy() {
	module_def.Close()
}

func (b *BaseEngine) OnSystemSignal(signal os.Signal) bool {
	return BaseSystemSignal(signal, b.servername)
}
