package etcdv3

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/timeutil"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
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
	logger      *zap.Logger
}

func NewService(etcd *Etcd, info ServiceInfo, opts ...BaseOptionFun) (service *Service, err error) {

	opt := &baseOptions{}
	for _, f := range opts {
		f(opt)
	}

	info.CreateTime = time.Now()
	service = &Service{
		ServiceInfo: info,
		Etcd:        etcd,
		logger:      opt.logger,
		stop:        make(chan error),
	}
	service.Init()
	return
}

func (s *Service) Init() {
	if s.logger == nil {
		s.logger, _ = logger.NewJSONLogger(
			logger.WithDisableConsole(),
			logger.WithField("domain", fmt.Sprintf("%s", "grpc_registry")),
			logger.WithTimeLayout(timeutil.CSTLayout),
			logger.WithFileP("./logs/", "grpc_registry.log"),
		)
	}
}

func (s *Service) Start() (err error) {
	ch, err := s.keepAlive()
	if err != nil {
		return err
	}
OUT:
	for {
		select {
		case err = <-s.stop:
			err = errors.New("关闭grpc server")
			break OUT
		case <-s.Etcd.client.Ctx().Done():
			err = errors.New("service closed")
			break OUT
		case _, ok := <-ch:
			// 监听租约
			if !ok {
				log.Printf("Recv reply from service: %s", s.getKey())
				err = s.revoke()
				if err != nil {
					s.logger.Info(fmt.Sprintf("revoke: %s", err.Error()))
				}
				err = errors.New("chan被关闭,由断点调试等网络情况会使得keepAlive断开")
				s.logger.Error(fmt.Sprintf("etcd keepalive error:%s", err.Error()))

				go func() {
					printStr := "重新监听 etcd 的 keepAlive..."
					s.logger.Info(printStr)

					if err = s.Start(); err != nil {
						fmt.Println(err)
					}
				}()
				break OUT
			}
		}
	}

	return

}

func (s *Service) Stop() {
	err := s.revoke()
	if err != nil {
		s.logger.Error("解除续约失败:", zap.Error(err))
	} else {
		s.logger.Info("解除续约成功")
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
