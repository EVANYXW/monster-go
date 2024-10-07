package etcdv3

import (
	"context"
	"encoding/json"
	"go.etcd.io/etcd/client/v3"
	"strings"
)

type Application struct {
	namespace string `json:"namespace"`
	etcd      *Etcd
}

type ApplicationData struct {
	Name string   `json:"name"`
	Num  int      `json:"num"`
	List []string `json:"list"`
}

func NewApplication(etcd *Etcd, namespace string) Application {
	return Application{etcd: etcd, namespace: namespace}
}

func (c *Application) getKeys(param []string) string {
	p := make([]string, 0)
	p = append(p, BasePath)
	p = append(p, Registry)
	p = append(p, c.namespace)
	p = append(p, "schema")
	p = append(p, param...)
	return strings.Join(p, "/")
}

func (c *Application) LastKey(str string) string {
	arr := strings.Split(str, "/")
	if len(arr) <= 0 {
		return ""
	}
	return arr[len(arr)-1]
}

func (c *Application) LastTwoKey(str string) string {
	arr := strings.Split(str, "/")
	if len(arr) <= 1 {
		return ""
	}
	return arr[len(arr)-2]
}

func (c *Application) Gets(key string) ([]ApplicationData, error) {
	key = c.getKeys([]string{key})
	apps := make([]ApplicationData, 0)

	resp, err := c.etcd.client.Get(context.TODO(), key, clientv3.WithPrefix())
	if err != nil {
		return apps, err
	}

	if resp.Count <= 0 {
		return apps, nil
	}

	appMap := make(map[string]ApplicationData)
	for _, v1 := range resp.Kvs {
		info := ServiceInfo{}
		err = json.Unmarshal(v1.Value, &info)
		if err != nil {
			return apps, err
		}
		name := c.LastTwoKey(string(v1.Key))
		if v, ok := appMap[name]; ok {
			v.Num += 1
			v.List = append(v.List, info.Address)
			appMap[name] = v
		} else {
			appMap[name] = ApplicationData{
				Name: name,
				Num:  1,
				List: []string{info.Address},
			}
		}
	}
	for _, v := range appMap {
		apps = append(apps, v)
	}
	return apps, nil
}
