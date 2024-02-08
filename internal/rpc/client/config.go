package client

import "github.com/evanyxw/monster-go/pkg/etcdv3"

var (
	etcdConfig *etcdv3.Etcd
)

func Init(etcd *etcdv3.Etcd) {
	etcdConfig = etcd
}
