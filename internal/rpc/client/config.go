package client

import "hl.hexinchain.com/welfare-center/basic/etcdv3"

var (
	etcdConfig *etcdv3.Etcd
)

func Init(etcd *etcdv3.Etcd) {
	etcdConfig = etcd
}
