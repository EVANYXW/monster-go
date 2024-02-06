package cmd

import (
	"bilibili/monster-go/internal/network"
	"bilibili/monster-go/internal/pkg/output"
	"bilibili/monster-go/pkg/async"
	"bilibili/monster-go/pkg/env"
	"fmt"
	"net/http"
	_ "net/http/pprof" // for side effects only

	"bilibili/monster-go/configs"
	"bilibili/monster-go/internal/mysql"
	"bilibili/monster-go/internal/redis"
	"bilibili/monster-go/internal/rpc/client"
	"bilibili/monster-go/internal/server"
	"bilibili/monster-go/internal/server/world"

	"github.com/phuhao00/sugar"
	"github.com/spf13/cobra"
	"hl.hexinchain.com/welfare-center/basic/etcdv3"
	"hl.hexinchain.com/welfare-center/basic/logs"
	"hl.hexinchain.com/welfare-center/basic/service"
)

var (
	serverName string
)

// Server cmd的server
type Server interface {
	Start()
	Stop()
}

func init() {
	ServerCmd.Flags().StringVar(&serverName, "server_name", "", "server_name")
}

// ServerCmd server 服务的cmd方法、
var ServerCmd = &cobra.Command{
	Use:   "run",
	Short: "games world server",
	Run: func(cmd *cobra.Command, args []string) {

		if serverName == "" {
			fmt.Println("Please specify a server name")
			return
		}

		allConfig := configs.Get()
		serverInfo := network.Info{
			ServerName: serverName,
			Address:    allConfig.Server.Address,
			RpcAddr:    allConfig.Rpc.Address,
			Env:        env.Active().Value(),
		}

		{
			InitLog() // 主要etcd 包里在用
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
			if allConfig.Server.Pprof {
				async.Go(func() {
					pprofs(allConfig.Server.PprofAddress)
				})
			}

			async.Go(func() {
				output.NewOutput(serverInfo.ServerName, serverInfo.Address, serverInfo.RpcAddr)
			})
		}

		fmt.Println(fmt.Sprintf("Starting【 %s 】server...", serverName))

		var server Server
		switch serverName {
		case "world":
			server = world.NewWorld(serverInfo)
		}

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
	},
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

// InitLog init log
func InitLog() {
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
