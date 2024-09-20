package cmd

import (
	"fmt"
	"log"

	//"github.com/evanyxw/monster-go/internal/mysql"
	//"github.com/evanyxw/monster-go/internal/redis"
	//"github.com/evanyxw/monster-go/internal/rpc/client"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/etcdv3"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/logs"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/server"
	"github.com/evanyxw/monster-go/pkg/server/engine"
	"github.com/phuhao00/sugar"
	"go.uber.org/zap"
	"net/http"
)

//type EtcdServers []*etcdv3.Service
//
//var etcdServerArr EtcdServers
//
//func (e EtcdServers) add(server *etcdv3.Service) {
//	etcdServerArr = append(e, server)
//}

// func Init(allConfig configs.Config, serverInfo server.Info) {
func Init(serverInfo server.Info) {
	{
		server.SID.Type = server.Name2EP(serverInfo.ServerName)
		server.SID.ID = 1
		if server.Name2EP(serverInfo.ServerName) == server.EP_Center {
			server.SID.Index = 1
		}
		logger.Info("SID:", zap.Uint16("ID", server.SID.ID),
			zap.Uint8("type", server.SID.Type), zap.Uint8("index", server.SID.Index))

		server.UpdateID()

		module.Init(module.ModuleID_Max)

		initLog(serverInfo.ServerName) // 主要etcd 包里在用
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

func Run(serverName string) {
	//serverInfo := server.Info{
	//	ServerName: serverName,
	//	Env: env.Active().Value(),
	//}
	//Init(allConfig, serverInfo)
	//Init(serverInfo)

	var serverKernel engine.IServerKernel
	instance := engine.MakeInstance(serverName)
	if instance == nil {
		log.Fatalf(fmt.Sprintf("找不到[%v]的服务,请通过 engine.Register 进行注册", serverName))
	}

	//server.SetServerInfo(&serverInfo)

	serverKernel = instance()
	//内部服务启动
	serverKernel.Run()
	defer func() {
		//内部服务关闭
		serverKernel.Destroy()
		release(serverName)
	}()

	sugar.WaitSignal(serverKernel.OnSystemSignal)
}

func registerEtcd(etcd *etcdv3.Etcd, serverName, address string) *etcdv3.Service {
	tcpEtcdServe, err := etcdv3.NewService(etcd, etcdv3.ServiceInfo{Name: serverName, Address: address})
	if err != nil {
		panic(err)
	}
	async.Go(func() {
		if err = tcpEtcdServe.Start(); err != nil {
			fmt.Println(err)
		}
	})

	return tcpEtcdServe
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

//func initEtcd() *etcdv3.Etcd {
//	config := configs.Get()
//
//	etcd, err := etcdv3.NewEtcd(config.Etcd.Address, config.Etcd.Namespace, config.Etcd.Secret,
//		config.Etcd.Namespace, 3, nil)
//	if err != nil {
//		panic(err)
//	}
//
//	client.Init(etcd)
//	return etcd
//}

//func newModel() {
//	redisConfig := configs.Get().Redis
//
//	// redis
//	redis.RedisManagers = redis.NewRedisManager(&redis.RedisConfig{
//		Passwd:    redisConfig.Pass,
//		Addr:      []string{redisConfig.Addr},
//		PoolSize:  redisConfig.PoolSize,
//		IsCluster: redisConfig.IsCluster,
//	})
//
//	// mysql
//	var err error
//	mysql.DBRepo, err = mysql.New()
//	if err != nil {
//		//Logger.Error("[Mysql] is failed:", zap.Error(err))
//		panic(err)
//	}
//
//}

//func redisSub(subFun redis.SubFun) {
//	redis.RedisManagers.Sub(subFun, redis.RedisPublistChannels...)
//}

//func recvPublish(channel string, data string) {
//	// TODO: By subscribing to and publishing through Redis,
//	// TODO: some control measures such as configuration rereading can be implemented
//	// TODO: 通过redis的订阅发布，可以实现一些配置重读等控制
//
//	//fmt.Println(channel, data)
//}

func pprofs(addr string) {
	http.ListenAndServe(addr, nil)
}

func release(serverName string) {
	logger.Info(fmt.Sprintf("【 %s 】Stopping server...", serverName))

	//mysql.DBRepo.DbRClose()
	//mysql.DBRepo.DbWClose()
	//redis.RedisManagers.Close()
	//for _, etcdServer := range etcdServerArr {
	//	etcdServer.Stop()
	//}

	logger.Info(fmt.Sprintf("【 %s 】server is stoped", serverName))

}
