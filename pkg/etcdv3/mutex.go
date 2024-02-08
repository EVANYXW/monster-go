package etcdv3

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"strings"
)

type MutexLock struct {
	etcd    *Etcd
	mutex   *concurrency.Mutex
	session *concurrency.Session
	key     string
}

func NewEtcdMutexV2(etcd *Etcd, ttl int64, key string) (*MutexLock, error) {
	session, err := concurrency.NewSession(etcd.client, concurrency.WithTTL(int(ttl)))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	p := make([]string, 0)
	p = append(p, BasePath)
	p = append(p, KvDataPath)
	p = append(p, etcd.namespace)
	mx := concurrency.NewMutex(session, strings.Join(p, "/")+key)
	return &MutexLock{
		session: session,
		mutex:   mx,
	}, nil
}

func (m *MutexLock) Lock() error {
	if err := m.mutex.Lock(context.TODO()); err != nil {
		fmt.Println("获取锁err:", err)
		return err
	}
	return nil
}

func (m *MutexLock) TryLock() error {
	if err := m.mutex.TryLock(context.TODO()); err != nil {
		fmt.Println("获取锁err:", err)
		return err
	}
	return nil
}

func (m *MutexLock) UnLock() error {
	defer m.session.Close()
	err := m.mutex.Unlock(context.TODO())
	if err != nil {
		fmt.Println("释放锁err:", err)
		return err
	}
	return nil
}
