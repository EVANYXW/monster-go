package client

import (
	"bilibili/monster-go/internal/rpc/client"
	"context"
	"fmt"
	"github.com/evanyxw/game_proto/msg"
	"hl.hexinchain.com/welfare-center/basic/etcdv3"
	"hl.hexinchain.com/welfare-center/basic/logs"
	"hl.hexinchain.com/welfare-center/basic/service"
)

func Broadcast() {
	logs.NewLogger(
		logs.WithFilePath(fmt.Sprintf("log/%s.log", service.Merchant)),
		logs.WithCompress(false),
		logs.WithPrettyPrint(false),
		logs.WithFormat("json"),
		logs.WithLevel(5),
		logs.WithMaxSize(100),
		logs.WithServerName(service.Merchant),
	)
	
	etcd, err := etcdv3.NewEtcd([]string{"127.0.0.1:2379"}, "default", "",
		"default", 3, nil)
	if err != nil {
		panic(err)
	}

	client.Init(etcd)
	rpcClient, err := client.NewWorldRpcClient()
	if err != nil {
		panic(err)
	}

	rpcClient.Broadcast(context.Background(), &msg.Req{})
}
