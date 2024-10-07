package etcdv3

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/evanyxw/monster-go/pkg/utils"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

type KvData struct {
	namespace string `json:"namespace"`
	etcd      *Etcd
}

type Kv struct {
	Key        string `json:"key"`
	Val        string `json:"val"`
	UpdateTime string `json:"updateTime"`
	Desc       string `json:"desc"`
}

func NewKvData(etcd *Etcd, namespace string) KvData {
	return KvData{etcd: etcd, namespace: namespace}
}

func (c *KvData) getKeys(param []string) string {
	p := make([]string, 0)
	p = append(p, BasePath)
	p = append(p, KvDataPath)
	p = append(p, c.namespace)
	p = append(p, param...)
	return strings.Join(p, "/")
}

func (c *KvData) Set(key, val, desc string) error {
	configVal := ConfigData{
		Key:        key,
		Val:        val,
		UpdateTime: utils.DateString(time.Now().Unix(), utils.DATE_FORMAT),
		Desc:       desc,
	}
	bytes, err := json.Marshal(configVal)
	if err != nil {
		return err
	}
	key = c.getKeys([]string{key})
	_, err = c.etcd.client.Put(context.TODO(), key, string(bytes))
	if err != nil {
		return err
	}
	return nil
}

func (c *KvData) Del(key string) error {
	configData, err := c.Get(key)
	if err != nil {
		return err
	}
	if configData.UpdateTime == "" {
		return errors.New("数据不存在," + key)
	}
	key = c.getKeys([]string{key})
	resp, err := c.etcd.client.Delete(context.TODO(), key)
	if err != nil {
		return err
	}
	if resp.Deleted <= 0 {
		return errors.New("删除数量为0")
	}
	return nil
}

func (c *KvData) Gets(key string) ([]ConfigData, error) {
	key = c.getKeys([]string{key})
	configVals := make([]ConfigData, 0)
	resp, err := c.etcd.client.Get(context.TODO(), key, clientv3.WithPrefix())
	if err != nil {
		return configVals, err
	}

	if resp.Count <= 0 {
		return configVals, nil
	}

	for _, v := range resp.Kvs {
		temp := ConfigData{}
		err = json.Unmarshal(v.Value, &temp)
		if err != nil {
			return configVals, err
		}
		configVals = append(configVals, temp)
	}
	return configVals, nil
}

func (c *KvData) Get(key string) (ConfigData, error) {
	key = c.getKeys([]string{key})
	configVal := ConfigData{}
	resp, err := c.etcd.client.Get(context.TODO(), key)
	if err != nil {
		return ConfigData{}, err
	}
	if resp.Count <= 0 {
		return ConfigData{}, nil
	}
	if len(resp.Kvs) != 1 {
		return ConfigData{}, errors.New("数据异常")
	}
	err = json.Unmarshal(resp.Kvs[0].Value, &configVal)
	if err != nil {
		return ConfigData{}, err
	}
	return configVal, nil
}
