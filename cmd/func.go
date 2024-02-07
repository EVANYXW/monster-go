package cmd

import (
	"bilibili/monster-go/configs"
	"bilibili/monster-go/internal/mysql"
	"bilibili/monster-go/internal/network"
	"bilibili/monster-go/internal/pkg/output"
	"bilibili/monster-go/internal/redis"
	"bilibili/monster-go/internal/rpc/client"
	"bilibili/monster-go/internal/server"
	"bilibili/monster-go/internal/server/factory"
	"bilibili/monster-go/internal/server/world"
	"bilibili/monster-go/pkg/async"
	"bilibili/monster-go/pkg/env"
	"fmt"
	"github.com/phuhao00/sugar"
	"hl.hexinchain.com/welfare-center/basic/etcdv3"
	"hl.hexinchain.com/welfare-center/basic/logs"
	"hl.hexinchain.com/welfare-center/basic/service"
	"net/http"
	"time"
)

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
		defer tcpEtcd.Stop()

		// rpc 服务注册etcd
		rpcEtcd := registerEtcd(etcd, fmt.Sprintf("%s%s", serverName, server.Rpc), serverInfo.RpcAddr)
		defer rpcEtcd.Stop()

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
	{
		//内部服务启动
		server.Start()
		defer func() {
			//内部服务关闭
			server.Stop()
			close()
		}()

		fmt.Println(fmt.Sprintf("【 %s 】server is started", serverName))
		sugar.WaitSignal(world.Oasis.OnSystemSignal)
	}
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

	time.Sleep(time.Duration(1) * time.Second)
	return tcpEtcdServe
}

// initLog init log
func initLog() {
	logs.NewLogger(
		logs.WithFilePath(fmt.Sprintf("log/%s.log", service.Merchant)),
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

	client.Init(etcd) // fixMe 这里初始化有些鸡肋
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
	// TODO: 可以实现好一些配置重读等控制
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

	fmt.Println(fmt.Sprintf("【 %s 】server is stoped", serverName))

}
