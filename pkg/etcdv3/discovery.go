/**
 * @api post etcdv3.
 *
 * User: yunshengzhu
 * Date: 2021/12/18
 * Time: 下午4:32
 */
package etcdv3

import (
	"context"
	"go.etcd.io/etcd/client/v3"
	"sync"
)

type ServiceDiscovery struct {
	etcd    *Etcd
	servers map[string]string
	lock    sync.RWMutex
}

func NewServiceDiscovery(etcd *Etcd) (*ServiceDiscovery, error) {
	return &ServiceDiscovery{
		etcd:    etcd,
		servers: make(map[string]string),
	}, nil
}

func (s *ServiceDiscovery) WatchService(prefix string) error {
	resp, err := s.etcd.client.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, kv := range resp.Kvs {
		s.SetServices(string(kv.Key), string(kv.Value))
	}
	go s.watcher(prefix)
	return nil
}

func (s *ServiceDiscovery) watcher(prefix string) {
	rch := s.etcd.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for watch := range rch {
		for _, event := range watch.Events {
			switch int(event.Type) {
			case 0:
				s.SetServices(string(event.Kv.Key), string(event.Kv.Key))
			case 1:
				s.DelServices(string(event.Kv.Key))
			}
		}
	}
}

func (s *ServiceDiscovery) SetServices(key, val string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.servers[key] = val
}

func (s *ServiceDiscovery) DelServices(key string) {
	s.lock.Unlock()
	defer s.lock.Unlock()
	delete(s.servers, key)
}

func (s *ServiceDiscovery) GetServices() []string {
	s.lock.Lock()
	defer s.lock.RUnlock()
	addrs := make([]string, 0, len(s.servers))
	for _, v := range s.servers {
		addrs = append(addrs, v)
	}
	return addrs
}
