package engine

import (
	"fmt"
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/internal/servers/gate/handler"
	"github.com/evanyxw/monster-go/pkg/async"
	commonModule "github.com/evanyxw/monster-go/pkg/common/module"
	"github.com/evanyxw/monster-go/pkg/env"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/logs"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/module/connector"
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
	serverName string
	isPprof    bool
	isOutput   bool
	output     *output.Output
}

type Options func(opt *option)

type option struct {
	isPprof    bool
	isOutput   bool
	output     *output.Config
	modules    map[int32]module.IModule
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

func WithModule(m module.IModule) Options {
	return func(opt *option) {
		//opt.modules[m.GetID()] = m
		module.NewBaseModule(m.GetID(), m)
	}
}

func WithModules(modules map[int32]module.IModule) Options {
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

func (b *BaseEngine) WithModule(m module.IModule) *BaseEngine {
	module.NewBaseModule(m.GetID(), m)
	return b
}

func (b *BaseEngine) WithOutput(config *output.Config) *BaseEngine {
	//b.output = output.NewOutput(config)
	b.isOutput = true
	b.output = output.NewOutput(config)
	async.Go(func() {
		b.output.Run()
	})
	return b
}

// initLog init log  etcd 在用
func initLog(serverName string) {
	logs.NewLogger(
		logs.WithFilePath(fmt.Sprintf("log/%s.log", serverName)),
		logs.WithCompress(false),
		logs.WithPrettyPrint(false),
		logs.WithFormat("json"),
		logs.WithLevel(5),
		logs.WithMaxSize(100),
		logs.WithServerName(serverName),
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

		module.Init(module.ModuleID_Max)

		initLog(server.GetServerInfo().ServerName) // 主要etcd 包里在用
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
		//rpcEtcd := registerEtcd(etcd, fmt.Sprintf("%s%s", serverName, servers.Rpc), serverInfo.RpcAddr)
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
		modules: make(map[int32]module.IModule),
	}
	for _, fn := range options {
		fn(opt)
	}

	// 装载注册发现模块
	//opt.modules[module.ModuleID_SM] = registerDiscovery

	b := &BaseEngine{
		serverName: name,
		isOutput:   opt.isOutput,
		isPprof:    opt.isPprof,
	}

	if opt.isOutput {
		b.output = output.NewOutput(opt.output)
		async.Go(func() {
			b.output.Run()
		})
	}

	//for id, m := range opt.modules {
	//	module.NewBaseModule(id, m)
	//}

	return b
}

func NewGateTcpServer(name string, factor register_discovery.ConnectorFactory, options ...Options) *BaseEngine {
	if factor == nil {
		log.Fatal("Please provide a factor!")
	}

	ServerInit(name)
	rd := factor.CreateConnector()
	options = append(options, WithModule(rd), WithModule(commonModule.NewClientNet(
		module.ModuleID_Client,
		5000,
		handler.NewGateMsg(),
		module.Outer,
		new(network.DefaultPackerFactory),
	)))

	if factor.IsConnectorServer() {
		options = append(options, WithModule(connector.NewManager(module.ModuleID_ConnectorManager)))
	}

	return newServer(name, options...)
}

func NewTcpServer(name string, msgHandler module.MsgHandler, factor register_discovery.ConnectorFactory, options ...Options) *BaseEngine {
	if factor == nil {
		log.Fatal("Please provide a factor!")
	}

	ServerInit(name)

	// 注册与发现
	rd := factor.CreateConnector()
	// 网络模块
	clientNet := commonModule.NewClientNet(
		module.ModuleID_GateAcceptor,
		10000,
		msgHandler,
		module.Inner,
		new(network.ClientPackerFactory),
	)
	options = append(options, WithModule(rd), WithModule(clientNet))

	return newServer(name, options...)
}

func NewGrpcServer(name string, msgHandler module.MsgHandler, factor register_discovery.ConnectorFactory, options ...Options) *BaseEngine {
	if factor == nil {
		log.Fatal("Please provide a factor!")
	}

	ServerInit(name)

	// 注册与发现
	rd := factor.CreateConnector()
	// 网络模块
	clientNet := commonModule.NewClientNet(
		module.ModuleID_GateAcceptor,
		10000,
		msgHandler,
		module.Inner,
		new(network.ClientPackerFactory),
	)
	options = append(options, WithModule(rd), WithModule(clientNet))

	return newServer(name, options...)
}

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
	module.Run()
}

func (b *BaseEngine) Destroy() {
	module.Close()
}

func (b *BaseEngine) OnSystemSignal(signal os.Signal) bool {
	return BaseSystemSignal(signal, b.serverName)
}
