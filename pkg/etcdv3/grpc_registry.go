package etcdv3

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/evanyxw/monster-go/pkg/logs"
	"go.etcd.io/etcd/client/v3"
	"log"
	"strings"
	"time"
)

// 服务信息
type ServiceInfo struct {
	Name       string
	Address    string
	CreateTime time.Time
}

type Service struct {
	ServiceInfo ServiceInfo
	stop        chan error
	leaseId     clientv3.LeaseID
	Etcd        *Etcd
}

func NewService(etcd *Etcd, info ServiceInfo) (service *Service, err error) {
	info.CreateTime = time.Now()
	service = &Service{
		ServiceInfo: info,
		Etcd:        etcd,
		stop:        make(chan error),
	}
	return
}

func (s *Service) Start() (err error) {
	ch, err := s.keepAlive()
	if err != nil {
		return err
	}
	for {
		select {
		case err := <-s.stop:
			return err
		case <-s.Etcd.client.Ctx().Done():
			return errors.New("service closed")
		case _, ok := <-ch:
			// 监听租约
			if !ok {
				log.Printf("Recv reply from service: %s", s.getKey())
				err = s.revoke()
				if err != nil {
					log.Printf("revoke: %s", err.Error())
				}
				return errors.New("退出服务注册")
			}
		}
	}
}

func (s *Service) Stop() {
	err := s.revoke()
	if err != nil {
		logs.Log.Error("解除续约失败:", err.Error())
	} else {
		logs.Log.Info("解除续约成功")
	}
	s.stop <- nil
}

func (s *Service) getKeys(param []string) string {
	p := make([]string, 0)
	p = append(p, BasePath)
	p = append(p, Registry)
	p = append(p, s.Etcd.namespace)
	p = append(p, schema)
	p = append(p, param...)
	return strings.Join(p, "/")
}

func (s *Service) keepAlive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	key := s.getKeys([]string{s.ServiceInfo.Name, s.ServiceInfo.Address})
	val, _ := json.Marshal(s.ServiceInfo)
	// 创建一个租约
	resp, err := s.Etcd.client.Grant(context.Background(), 10)
	if err != nil {
		return nil, err
	}
	_, err = s.Etcd.client.Put(context.Background(), key, string(val), clientv3.WithLease(resp.ID))
	if err != nil {
		return nil, err
	}
	s.leaseId = resp.ID
	return s.Etcd.client.KeepAlive(context.TODO(), s.leaseId)
}

func (s *Service) revoke() error {
	_, err := s.Etcd.client.Revoke(context.TODO(), s.leaseId)
	return err
}

func (s *Service) getKey() string {
	return s.ServiceInfo.Name + "/" + s.ServiceInfo.Address
}
