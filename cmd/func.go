package cmd

import (
	"fmt"
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/internal/mysql"
	"github.com/evanyxw/monster-go/internal/network"
	"github.com/evanyxw/monster-go/internal/pkg/output"
	"github.com/evanyxw/monster-go/internal/redis"
	"github.com/evanyxw/monster-go/internal/rpc/client"
	"github.com/evanyxw/monster-go/internal/server"
	"github.com/evanyxw/monster-go/internal/server/factory"
	"github.com/evanyxw/monster-go/internal/server/world"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/env"
	"github.com/evanyxw/monster-go/pkg/etcdv3"
	"github.com/evanyxw/monster-go/pkg/logs"
	"github.com/phuhao00/sugar"
	"net/http"
)

type EtcdServers []*etcdv3.Service

var etcdServerArr EtcdServers

func (e EtcdServers) add(server *etcdv3.Service) {
	etcdServerArr = append(e, server)
}

func Init(allConfig configs.Config, serverInfo network.Info) {
	{
		initLog() // 主要etcd 包里在用
		// mysql、redis
		newModel()
		// redis 发布订阅
		redisSub(recvPublish)

		// 初始化etcd
		etcd := initEtcd()

		// tcp 服务注册etcd
		tcpEtcd := registerEtcd(etcd, serverInfo.ServerName, serverInfo.Address)
		etcdServerArr.add(tcpEtcd)

		// rpc 服务注册etcd
		rpcEtcd := registerEtcd(etcd, fmt.Sprintf("%s%s", serverName, server.Rpc), serverInfo.RpcAddr)
		etcdServerArr.add(rpcEtcd)

		// 开启pprof
		pprofUrl := ""
		if allConfig.Server.Pprof {
			pprofUrl = allConfig.Server.PprofAddress
			async.Go(func() {
				pprofs(allConfig.Server.PprofAddress)
			})
		}

		async.Go(func() {
			output.NewOutput(serverInfo.ServerName, serverInfo.Address, serverInfo.RpcAddr,
				pprofUrl)
		})
	}
}

func Run(serverName string) {
	allConfig := configs.Get()
	serverInfo := network.Info{
		ServerName: serverName,
		Address:    allConfig.Server.Address,
		RpcAddr:    allConfig.Rpc.Address,
		Env:        env.Active().Value(),
	}

	Init(allConfig, serverInfo)

	fmt.Println(fmt.Sprintf("Starting【 %s 】server...", serverName))

	var server factory.Server
	instance := factory.MakeInstance(serverName)
	if instance == nil {
		panic("找不到对应服务")
	}
	server = instance(serverInfo)
	//内部服务启动
	server.Run()
	defer func() {
		//内部服务关闭
		server.Destroy()
		close()
	}()

	fmt.Println(fmt.Sprintf("【 %s 】server is started", serverName))
	sugar.WaitSignal(world.Oasis.OnSystemSignal)
}

func registerEtcd(etcd *etcdv3.Etcd, serverName, address string) *etcdv3.Service {
	tcpEtcdServe, err := etcdv3.NewService(etcd, etcdv3.ServiceInfo{Name: serverName, Address: address})
	async.Go(func() {
		if err != nil {
			panic(err)
		}
		if err = tcpEtcdServe.Start(); err != nil {
			fmt.Println(err)
		}
	})

	return tcpEtcdServe
}

// initLog init log
func initLog() {
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

func initEtcd() *etcdv3.Etcd {
	config := configs.Get()

	etcd, err := etcdv3.NewEtcd(config.Etcd.Address, config.Etcd.Namespace, config.Etcd.Secret,
		config.Etcd.Namespace, 3, nil)
	if err != nil {
		panic(err)
	}

	client.Init(etcd)
	return etcd
}

func newModel() {
	redisConfig := configs.Get().Redis

	// redis
	redis.RedisManagers = redis.NewRedisManager(&redis.RedisConfig{
		Passwd:    redisConfig.Pass,
		Addr:      []string{redisConfig.Addr},
		PoolSize:  redisConfig.PoolSize,
		IsCluster: redisConfig.IsCluster,
	})

	// mysql
	var err error
	mysql.DBRepo, err = mysql.New()
	if err != nil {
		//Logger.Error("[Mysql] is failed:", zap.Error(err))
		panic(err)
	}

}

func redisSub(subFun redis.SubFun) {
	redis.RedisManagers.Sub(subFun, redis.RedisPublistChannels...)
}

func recvPublish(channel string, data string) {
	// TODO: 通过redis的订阅发布，可以实现一些配置重读等控制
	fmt.Println(channel, data)
}

func pprofs(addr string) {
	http.ListenAndServe(addr, nil)
}

func close() {
	fmt.Println(fmt.Sprintf("Stopping【 %s 】server...", serverName))

	mysql.DBRepo.DbRClose()
	mysql.DBRepo.DbWClose()
	redis.RedisManagers.Close()
	for _, etcdServer := range etcdServerArr {
		etcdServer.Stop()
	}

	fmt.Println(fmt.Sprintf("【 %s 】server is stoped", serverName))

}
