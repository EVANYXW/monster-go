package engine

import (
	"fmt"
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/env"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/logs"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/module/register-discovery"
	"github.com/evanyxw/monster-go/pkg/output"
	"github.com/evanyxw/monster-go/pkg/server"
	"github.com/evanyxw/monster-go/pkg/timeutil"
	"go.uber.org/zap"
	"log"
	"os"
)

type BaseEngine struct {
	serverName string
	isPprof    bool
	isOutput   bool
	output     *output.Output
}

type Options func(opt *option)

type option struct {
	isPprof  bool
	isOutput bool
	output   *output.Config
	modules  map[int32]module.IModule
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
		opt.modules[m.GetID()] = m
	}
}

func WithModules(modules map[int32]module.IModule) Options {
	return func(opt *option) {
		for id, m := range opt.modules {
			opt.modules[id] = m
		}
	}
}

func newServer(name string, registerDiscovery module.IModule, options ...Options) *BaseEngine {
	server.SetServerInfo(&server.Info{
		ServerName: name,
		Env:        env.Active().Value(),
	})

	configs.Init()
	Init()
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
	}

	for id, m := range opt.modules {
		module.NewBaseModule(id, m)
	}

	return b
}

func NewServer(name string, factor register_discovery.ConnectorFactory, options ...Options) *BaseEngine {
	if factor == nil {
		log.Fatal("Please provide a factor!")
	}

	rd := factor.CreateConnector()
	return newServer(name, rd, options...)
}

func NewCenterServer(name string, factor register_discovery.NetFactory, options ...Options) *BaseEngine {
	if factor == nil {
		log.Fatal("Please provide a factor!")
	}

	rd := factor.CreateNet()
	return newServer(name, rd, options...)
}

func (b *BaseEngine) WithModule(m module.IModule) *BaseEngine {
	module.NewBaseModule(m.GetID(), m)
	return b
}

func (b *BaseEngine) WithOutput(config *output.Config) *BaseEngine {
	b.output = output.NewOutput(config)
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
		//	output.NewOutput(serverInfo.ServerName, serverInfo.Address,
		//		pprofUrl).Run()
		//})
	}
}

func (b *BaseEngine) Run() {
	if b.output != nil {
		async.Go(func() {
			b.output.Run()
		})
	}
	logger.Info(fmt.Sprintf("【 %s 】Starting server...", server.GetServerInfo().ServerName))

	module.Run()
}

func (b *BaseEngine) Destroy() {
	module.Close()
}

func (b *BaseEngine) OnSystemSignal(signal os.Signal) bool {
	return BaseSystemSignal(signal, b.serverName)
}
