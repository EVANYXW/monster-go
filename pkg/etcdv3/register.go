package etcdv3

import (
	"context"
	"go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

type ServiceRegister struct {
	client        *clientv3.Client
	leaseId       clientv3.LeaseID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string
	val           string
}

func NewServiceRegister(endpoints []string, key, val string, lease int64, dialTimeout int) (*ServiceRegister, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Duration(dialTimeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}
	service := &ServiceRegister{
		client: client,
		key:    key,
		val:    val,
	}
	if err := service.putKeyWithLease(lease); err != nil {
		return nil, err
	}
	return service, nil
}

func (s *ServiceRegister) putKeyWithLease(lease int64) error {

	resp, err := s.client.Grant(context.Background(), lease)
	if err != nil {
		return err
	}

	_, err = s.client.Put(context.Background(), s.key, s.val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}

	leaseRespChan, err := s.client.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return err
	}

	s.leaseId = resp.ID
	s.keepAliveChan = leaseRespChan
	return nil
}

func (s *ServiceRegister) ListenLeaseRespChan() {

	for leaseKeepResp := range s.keepAliveChan {
		log.Println("续约成功:", leaseKeepResp)
	}
	log.Println("关闭续约")
}

func (s *ServiceRegister) Close() error {
	if _, err := s.client.Revoke(context.Background(), s.leaseId); err != nil {
		return err
	}
	return s.client.Close()
}
