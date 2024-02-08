/**
 * @api post etcdservice.
 *
 * User: yunshengzhu
 * Date: 2019-07-31
 * Time: 14:36
 */
package etcdv3

import (
	"context"
	"errors"
	"fmt"
	"github.com/evanyxw/monster-go/pkg/logs"
	"go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

type EtcdMutex struct {
	etcd    *Etcd
	Ttl     int64              //租约时间
	Key     string             //etcd的key
	cancel  context.CancelFunc //关闭续租的func
	lease   clientv3.Lease
	leaseID clientv3.LeaseID
	txn     clientv3.Txn
}

func NewEtcdMutex(etcd *Etcd, ttl int64, key string) (*EtcdMutex, error) {
	etcdMutex := &EtcdMutex{
		etcd: etcd,
		Ttl:  ttl,
	}
	etcdMutex.Key = etcdMutex.GetKey(key)
	return etcdMutex, nil
}

func (em *EtcdMutex) GetKey(key string) string {
	p := make([]string, 0)
	p = append(p, BasePath)
	p = append(p, KvDataPath)
	p = append(p, em.etcd.namespace)
	return strings.Join(p, "/") + key
}

func (em *EtcdMutex) Lock() error {
	nTime := time.Now().Add(time.Second * 30)
	var count int
	for {
		var err error
		var ctx context.Context
		em.txn = clientv3.NewKV(em.etcd.client).Txn(context.TODO())
		em.lease = clientv3.NewLease(em.etcd.client)
		leaseResp, err := em.lease.Grant(context.TODO(), em.Ttl)
		if err != nil {
			return err
		}
		ctx, em.cancel = context.WithCancel(context.TODO())
		em.leaseID = leaseResp.ID
		_, err = em.lease.KeepAlive(ctx, em.leaseID)
		if err != nil {
			return err
		}

		em.txn.If(clientv3.Compare(clientv3.CreateRevision(em.Key), "=", 0)).
			Then(clientv3.OpPut(em.Key, "", clientv3.WithLease(em.leaseID))).
			Else()
		txnResp, err := em.txn.Commit()
		if err != nil {
			return err
		}
		if txnResp.Succeeded {
			return nil
		} else {
			if time.Now().After(nTime) {
				return errors.New("获取锁失败")
			}
		}
		count++
		logs.Log.Debugf("第%d次获取锁失败，准备sleep", count)
		time.Sleep(time.Millisecond * 50)
	}
}

func (em *EtcdMutex) UnLock() {
	em.cancel()
	em.lease.Revoke(context.TODO(), em.leaseID)
	fmt.Println("释放了锁")
}
