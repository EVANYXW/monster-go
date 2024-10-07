package etcdv3

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/timeutil"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
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
	logger     *zap.Logger
}

func NewResolver(etcd *Etcd, service string, opts ...BaseOptionFun) resolver.Builder {
	opt := &baseOptions{}
	for _, f := range opts {
		f(opt)
	}

	r := &Resolver{etcd: etcd, service: service, logger: opt.logger}
	r.Init()
	return r
}

func (r *Resolver) Init() {
	// 初始化 cron logger
	if r.logger == nil {
		r.logger, _ = logger.NewJSONLogger(
			logger.WithDisableConsole(),
			logger.WithField("domain", fmt.Sprintf("%s", "grpc_resolver")),
			logger.WithTimeLayout(timeutil.CSTLayout),
			logger.WithFileP("./logs/", "grpc_resolver.log"),
		)
	}
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
		r.logger.Error("err:", zap.Error(err))
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
		//logs.Log.Error("更新状态失败1,err:%s,state:%+v", err, state)
		r.logger.Error("更新状态失败1,err:%s,state:%+v", zap.Error(err), zap.Any("state", state))
	}

	//logs.Log.Infof("更新可用服务信息,prefix:%s,%+v", r.prefix, state.Addresses)
	r.logger.Error("更新可用服务信息,prefix:%s,%+v", zap.String("prefix", r.prefix),
		zap.Any("state.address", state.Addresses))
}

func (r *Resolver) sync() error {
	resp, err := r.etcd.client.Get(context.Background(), r.prefix, clientv3.WithPrefix())
	if err == nil {
		tempSrvAddrMap := make(map[string]resolver.Address)
		for _, kv := range resp.Kvs {
			info := &ServiceInfo{}
			err := json.Unmarshal(kv.Value, info)
			if err != nil {
				//logs.Log.Error("err:", err)
				r.logger.Error("err: ", zap.Error(err))
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
	//logs.Log.Info("watch prefix:", r.prefix)
	r.logger.Error("watch prefix:", zap.Any("prefix", r.prefix))
	ticker := time.NewTicker(time.Minute * 10)
	var errNum int64
	watchCh := r.etcd.client.Watch(context.Background(), r.prefix, clientv3.WithPrefix())
	for {
		select {
		case res := <-watchCh:
			if res.Err() != nil {
				errNum++
				//logs.Log.Errorf("重新监听 watchCh read watch err:%+v", res.Err())
				r.logger.Error("重新监听 watchCh read watch err: ", zap.Error(res.Err()))
				err := r.etcd.Reconnect()
				if err != nil {
					//logs.Log.Errorf("etcd 重连  err:%+v", err)
					r.logger.Error("etcd 重连  err:", zap.Error(err))
				}
				watchCh = r.etcd.client.Watch(context.Background(), r.prefix, clientv3.WithPrefix())
				if errNum > 10 {
					//logs.Log.Error("服务监听异常系统退出:", res.Err())
					r.logger.Error("服务监听异常系统退出:", zap.Error(res.Err()))
					panic(res.Err())
				}
			} else {
				errNum = 0
				//logs.Log.Info("服务发现自动触发更新...")
				r.logger.Info("服务发现自动触发更新...")
				r.update(res.Events)
			}
		case <-ticker.C:
			if err := r.sync(); err != nil {
				//logs.Log.Info("sync failed", err)
				r.logger.Info("sync failed:", zap.Error(err))
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
				//logs.Log.Error("err:", err)
				r.logger.Error("err:", zap.Error(err))
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
