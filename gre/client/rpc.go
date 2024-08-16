package client

import (
	"fmt"
	"github.com/evanyxw/monster-go/pkg/logs"
)

func Broadcast() {
	logs.NewLogger(
		logs.WithFilePath(fmt.Sprintf("log/%s.log", "world")),
		logs.WithCompress(false),
		logs.WithPrettyPrint(false),
		logs.WithFormat("json"),
		logs.WithLevel(5),
		logs.WithMaxSize(100),
		logs.WithServerName("world"),
	)

	//etcd, err := etcdv3.NewEtcd([]string{"127.0.0.1:2379"}, "default", "",
	//	"default", 3, nil)
	//if err != nil {
	//	panic(err)
	//}
	//
	//client.Init(etcd)
	//rpcClient, err := client.NewWorldRpcClient()
	//if err != nil {
	//	panic(err)
	//}
	//
	//rpcClient.Broadcast(context.Background(), &msg.Req{})
}
