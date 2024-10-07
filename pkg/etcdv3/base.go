package etcdv3

import (
	"go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"strings"
)

const (
	BasePath      = "/etcd-v3"
	NamespacePath = "namespace"
	Registry      = "registry"
	ConfigPath    = "config"
	KvDataPath    = "kv-data"
)

type AuthorizedPath struct {
	Path    string
	PathEnd string
	Type    clientv3.PermissionType
}

type BaseOptionFun func(*baseOptions)
type baseOptions struct {
	logger *zap.Logger
}

// WithLogger zap.Logger
func WithLogger(logger *zap.Logger) BaseOptionFun {
	return func(opt *baseOptions) {
		opt.logger = logger
	}
}

func GetAuthPaths(namespace string) []AuthorizedPath {
	paths := make([]AuthorizedPath, 0)
	paths = append(paths, AuthorizedPath{
		Path:    strings.Join([]string{BasePath, ConfigPath, namespace}, "/") + "/",
		PathEnd: strings.Join([]string{BasePath, ConfigPath, namespace}, "/") + "0",
		Type:    clientv3.PermissionType(clientv3.PermRead),
	})
	paths = append(paths, AuthorizedPath{
		Path:    strings.Join([]string{BasePath, Registry, namespace}, "/") + "/",
		PathEnd: strings.Join([]string{BasePath, Registry, namespace}, "/") + "0",
		Type:    clientv3.PermissionType(clientv3.PermReadWrite),
	})
	paths = append(paths, AuthorizedPath{
		Path:    strings.Join([]string{BasePath, KvDataPath, namespace}, "/") + "/",
		PathEnd: strings.Join([]string{BasePath, KvDataPath, namespace}, "/") + "0",
		Type:    clientv3.PermissionType(clientv3.PermReadWrite),
	})
	return paths
}
