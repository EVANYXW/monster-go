package etcdv3

import (
	"context"
	"encoding/json"
	"github.com/evanyxw/monster-go/pkg/logs"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"strings"
	"time"
)

const schema = "schema"

type Resolver struct {
	srvAddrMap map[string]resolver.Address
	endpoints  []string
	service    string
	prefix     string
	etcd       *Etcd
	clientConn resolver.ClientConn
}

func NewResolver(etcd *Etcd, service string) resolver.Builder {
	return &Resolver{etcd: etcd, service: service}
}

func (r *Resolver) Scheme() string {
	return schema
}

func (r *Resolver) ResolveNow(rn resolver.ResolveNowOptions) {

}

func (r *Resolver) Close() {

}

func (r *Resolver) getRoot() string {
	p := make([]string, 0)
	p = append(p, BasePath)
	p = append(p, Registry)
	p = append(p, r.etcd.namespace)
	return strings.Join(p, "/")
}

func (r *Resolver) Build(target resolver.Target, clientConn resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.clientConn = clientConn
	r.prefix = r.getRoot() + "/" + target.URL.Scheme + target.URL.Path + "/"
	err := r.sync()
	if err != nil {
		logs.Log.Error("err:", err)
	}
	go r.watch()
	return r, nil
}

func (r *Resolver) updateState() {
	state := resolver.State{}
	for _, v := range r.srvAddrMap {
		state.Addresses = append(state.Addresses, v)
	}
	err := r.clientConn.UpdateState(state)
	if err != nil {
		logs.Log.Error("更新状态失败1,err:%s,state:%+v", err, state)
	}
	logs.Log.Infof("更新可用服务信息,prefix:%s,%+v", r.prefix, state.Addresses)
}

func (r *Resolver) sync() error {
	resp, err := r.etcd.client.Get(context.Background(), r.prefix, clientv3.WithPrefix())
	if err == nil {
		tempSrvAddrMap := make(map[string]resolver.Address)
		for _, kv := range resp.Kvs {
			info := &ServiceInfo{}
			err := json.Unmarshal(kv.Value, info)
			if err != nil {
				logs.Log.Error("err:", err)
				continue
			}
			tempSrvAddrMap[string(kv.Key)] = resolver.Address{Addr: info.Address}
		}
		r.srvAddrMap = make(map[string]resolver.Address)
		r.srvAddrMap = tempSrvAddrMap
		r.updateState()
	}
	return err
}

func (r *Resolver) watch() {
	logs.Log.Info("watch prefix:", r.prefix)
	ticker := time.NewTicker(time.Minute * 10)
	var errNum int64
	watchCh := r.etcd.client.Watch(context.Background(), r.prefix, clientv3.WithPrefix())
	for {
		select {
		case res := <-watchCh:
			if res.Err() != nil {
				errNum++
				logs.Log.Errorf("重新监听 watchCh read watch err:%+v", res.Err())
				err := r.etcd.Reconnect()
				if err != nil {
					logs.Log.Errorf("etcd 重连  err:%+v", err)
				}
				watchCh = r.etcd.client.Watch(context.Background(), r.prefix, clientv3.WithPrefix())
				if errNum > 10 {
					logs.Log.Error("服务监听异常系统退出:", res.Err())
					panic(res.Err())
				}
			} else {
				errNum = 0
				logs.Log.Info("服务发现自动触发更新...")
				r.update(res.Events)
			}
		case <-ticker.C:
			if err := r.sync(); err != nil {
				logs.Log.Info("sync failed", err)
			}
			ticker.Reset(time.Minute * 10)
		}
	}
}

func (r *Resolver) update(events []*clientv3.Event) {
	for _, ev := range events {
		var info ServiceInfo
		switch ev.Type {
		case mvccpb.PUT:
			err := json.Unmarshal(ev.Kv.Value, &info)
			if err != nil {
				logs.Log.Error("err:", err)
				continue
			}
			if _, ok := r.srvAddrMap[string(ev.Kv.Key)]; !ok {
				r.srvAddrMap[string(ev.Kv.Key)] = resolver.Address{Addr: info.Address}
				r.updateState()
			}
		case mvccpb.DELETE:
			delete(r.srvAddrMap, string(ev.Kv.Key))
			r.updateState()
		}
	}
}
