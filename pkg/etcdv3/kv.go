package etcdv3

import (
	"context"
	"errors"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
)

func Put(etcd *Etcd, key, val string, ttl int64) error {
	resp, err := etcd.client.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}
	_, err = etcd.client.Put(context.Background(), key, val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}
	return nil
}

func Get(etcd *Etcd, key string) (*mvccpb.KeyValue, error) {
	resp, err := etcd.client.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) <= 0 {
		return nil, errors.New("数据不存在")
	}
	return resp.Kvs[0], nil
}

func Del(etcd *Etcd, key string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), etcd.ttl)
	defer cancel()
	// delete the keys
	_, err = etcd.client.Delete(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func KvPut(etcd *Etcd, key, val string) error {
	ctx, cancel := context.WithTimeout(context.Background(), etcd.ttl)
	_, err := etcd.client.Put(ctx, key, val)
	cancel()
	if err != nil {
		return err
	}
	return nil
}

func KvGet(etcd *Etcd, key string) (response *clientv3.GetResponse, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), etcd.ttl)
	resp, err := etcd.client.Get(ctx, key, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return resp, err
	}
	return resp, err
}

func KvDel(etcd *Etcd, key string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), etcd.ttl)
	defer cancel()
	gresp, err := etcd.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	// delete the keys
	dresp, err := etcd.client.Delete(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	if int64(len(gresp.Kvs)) == dresp.Deleted {
		return nil
	}
	return errors.New("删除的数量与查询的数量不一致")
}
