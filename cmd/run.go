package cmd

import (
	"bilibili/monster-go/configs"
	"bilibili/monster-go/internal/mysql"
	"bilibili/monster-go/internal/redis"
	"bilibili/monster-go/internal/rpc/client"
	"bilibili/monster-go/internal/server/world"
	"fmt"
	"github.com/phuhao00/sugar"
	"github.com/spf13/cobra"
	"hl.hexinchain.com/welfare-center/basic/etcdv3"
	"hl.hexinchain.com/welfare-center/basic/logs"
	"hl.hexinchain.com/welfare-center/basic/service"
	"net/http"
	_ "net/http/pprof"
)

var (
	serverName string
)

var (
	etcdServer *etcdv3.Service
)

type Server interface {
	Start()
	Stop()
}

var WorldCmd = &cobra.Command{
	Use:   "run",
	Short: "games world server",
	Run: func(cmd *cobra.Command, args []string) {
		if serverName == "" {
			fmt.Println("Please specify a server name")
			return
		}
		InitLog() // 主要etcd 包里在用
		// mysql、redis
		newModel()
		// redis 发布订阅
		redisSub(recvPublish)
		// 注册服务到etcd
		registerEtcd()

		pprofs()
		fmt.Println(fmt.Sprintf("Starting【 %s 】server...", serverName))

		var server Server
		switch serverName {
		case "world":
			server = world.NewWorld()
		}

		server.Start()
		fmt.Println(fmt.Sprintf("【 %s 】server is started", serverName))

		sugar.WaitSignal(world.Oasis.OnSystemSignal)

		fmt.Println(fmt.Sprintf("Stopping【 %s 】server...", serverName))
		server.Stop()
		close()

		fmt.Println(fmt.Sprintf("【 %s 】server is stoped", serverName))
	},
}

func init() {

	WorldCmd.Flags().StringVar(&serverName, "server_name", "", "server_name")
}

func InitLog() {
	logs.NewLogger(
		logs.WithFilePath(fmt.Sprintf("log/%s.log", service.Merchant)),
		logs.WithCompress(false),
		logs.WithPrettyPrint(false),
		logs.WithFormat("json"),
		logs.WithLevel(5),
		logs.WithMaxSize(100),
		logs.WithServerName(service.Merchant),
	)
}

func registerEtcd() {
	config := configs.Get()

	etcd, err := etcdv3.NewEtcd(config.Etcd.Address, config.Etcd.Namespace, config.Etcd.Secret,
		config.Etcd.Namespace, 3, nil)
	if err != nil {
		panic(err)
	}
	go func() {
		etcdServer, err = etcdv3.NewService(etcd, etcdv3.ServiceInfo{Name: serverName, Address: config.Server.Address})
		if err != nil {
			panic(err)
		}
		if err = etcdServer.Start(); err != nil {
			fmt.Println(err)
		}

	}()

	client.Init(etcd)
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

func pprofs() {
	serverConfig := configs.Get().Server
	fmt.Println(1111)
	if serverConfig.Pprof {
		go func() {
			http.ListenAndServe(serverConfig.PprofAddress, nil)
		}()
	}
}

func close() {
	mysql.DBRepo.DbRClose()
	mysql.DBRepo.DbWClose()
	redis.RedisManagers.Close()
	etcdServer.Stop()
}
