package etcdv3

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"strings"
)

type Namespaces struct {
	etcd *Etcd
}

type NamespaceData struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Secret    string `json:"secret"`
}

func NewNamespace(etcd *Etcd) *Namespaces {
	return &Namespaces{
		etcd: etcd,
	}
}

func (n *Namespaces) getKeys(p []string) string {
	param := make([]string, 0)
	param = append(param, BasePath)
	param = append(param, NamespacePath)
	param = append(param, p...)
	return strings.Join(param, "/")
}

func (n *Namespaces) getRootPath() string {
	param := make([]string, 0)
	param = append(param, BasePath)
	param = append(param, NamespacePath)
	return strings.Join(param, "/")
}

func (n *Namespaces) encode(obj interface{}) string {
	bytes, _ := json.Marshal(obj)
	return string(bytes)
}

func (n *Namespaces) Create(data NamespaceData) error {
	if len(data.Namespace) <= 0 {
		return errors.New("param namespace null")
	}
	_, err := n.etcd.client.Put(context.Background(), n.getKeys([]string{data.Namespace}), n.encode(data))
	return err
}

func (n *Namespaces) decode(data []byte, obj interface{}) error {
	return json.Unmarshal(data, obj)
}

func (n *Namespaces) List() (list []NamespaceData, err error) {
	resp, err := n.etcd.client.Get(context.Background(), n.getRootPath(), clientv3.WithPrefix())
	if err != nil {
		return list, err
	}
	list = make([]NamespaceData, 0)
	for _, v := range resp.Kvs {
		obj := NamespaceData{}
		err = n.decode(v.Value, &obj)
		if err != nil {
			return list, err
		}
		list = append(list, obj)
	}
	return list, nil
}

func (n *Namespaces) Del(namespace string) error {
	if len(namespace) <= 0 {
		return errors.New("param namespace null")
	}
	resp, err := n.etcd.client.Get(context.Background(), n.getKeys([]string{namespace}))
	if err != nil {
		return err
	}
	if len(resp.Kvs) <= 0 {
		return errors.New("namespace does not exist")
	}

	delResp, err := n.etcd.client.Delete(context.Background(), n.getKeys([]string{namespace}))
	if err != nil {
		return err
	}
	fmt.Println("删除数据:", delResp.Deleted)
	return nil
}
