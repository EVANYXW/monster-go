package etcdv3

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/evanyxw/monster-go/pkg/utils"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"strings"
	"time"
)

type Config struct {
	namespace string `json:"namespace"`
	etcd      *Etcd
}

type ConfigData struct {
	Key        string `json:"key"`
	Val        string `json:"val"`
	UpdateTime string `json:"updateTime"`
	Desc       string `json:"desc"`
}

func NewConfig(etcd *Etcd, namespace string) Config {
	return Config{etcd: etcd, namespace: namespace}
}

func (c *Config) getKeys(param []string) string {
	p := make([]string, 0)
	p = append(p, BasePath)
	p = append(p, ConfigPath)
	p = append(p, c.namespace)

	p = append(p, param...)
	return strings.Join(p, "/")
}

func (c *Config) Set(key, val, desc string) error {
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

func (c *Config) Del(key string) error {
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

func (c *Config) Gets(key string) ([]ConfigData, error) {
	key = c.getKeys([]string{key})
	configVals := make([]ConfigData, 0)
	fmt.Println("key:", key)
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

func (c *Config) Get(key string) (ConfigData, error) {
	key = c.getKeys([]string{key})
	fmt.Println("key:", key)
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

func LoadConfig(etcd *Etcd, key string, fn func([]byte)) error {
	p := make([]string, 0)
	p = append(p, BasePath)
	p = append(p, ConfigPath)
	p = append(p, etcd.namespace)
	p = append(p, key)
	key = strings.Join(p, "/")
	resp, err := etcd.client.Get(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		panic(err.Error() + "    " + key)
	}

	watchChan := etcd.client.Watch(context.Background(), key)
	go func() {
		for {
			select {
			case resp := <-watchChan:
				err := resp.Err()
				if err != nil {
					fmt.Println("err:", err)
					return
				}
				for _, ev := range resp.Events {
					switch ev.Type {
					case mvccpb.PUT:
						conf := ConfigData{}
						err = json.Unmarshal(ev.Kv.Value, &conf)
						if err != nil {
							log.Println("err:", err)
						}
						fn([]byte(conf.Val))
					}
				}
			}
		}
	}()

	if len(resp.Kvs) <= 0 {
		return nil
	}

	conf := ConfigData{}
	err = json.Unmarshal([]byte(resp.Kvs[0].Value), &conf)
	if err != nil {
		log.Println("err:", err)
	}
	fn([]byte(conf.Val))

	return nil
}

func LoadConfigV2(etcd *Etcd, key string, obj interface{}, fn func([]byte, interface{})) error {
	p := make([]string, 0)
	p = append(p, BasePath)
	p = append(p, ConfigPath)
	p = append(p, etcd.namespace)
	p = append(p, key)
	key = strings.Join(p, "/")
	resp, err := etcd.client.Get(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		panic(err.Error() + "    " + key)
	}

	watchChan := etcd.client.Watch(context.Background(), key)
	go func() {
		for {
			select {
			case resp := <-watchChan:
				err := resp.Err()
				if err != nil {
					fmt.Println("err:", err)
					return
				}
				for _, ev := range resp.Events {
					switch ev.Type {
					case mvccpb.PUT:
						conf := ConfigData{}
						err = json.Unmarshal(ev.Kv.Value, &conf)
						if err != nil {
							log.Println("err:", err)
						}
						fn([]byte(conf.Val), obj)
					}
				}
			}
		}
	}()

	if len(resp.Kvs) <= 0 {
		return nil
	}

	conf := ConfigData{}
	err = json.Unmarshal([]byte(resp.Kvs[0].Value), &conf)
	if err != nil {
		log.Println("err:", err)
	}
	fn([]byte(conf.Val), obj)
	return nil
}
