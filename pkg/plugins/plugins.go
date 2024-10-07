// Package plugins @Author evan_yxw
// @Date 2024/10/5 11:51:00
// @Desc
package plugins

import (
	"context"
	"errors"
	"fmt"
	"github.com/evanyxw/monster-go/pkg/plugins/transport/grpc"
	"log"
)

type PluginIns interface {
	Factory() Factory
}

type Factory interface {
	Type() Type
	Name() string
	Setup(context.Context) (PluginIns, error)
	Destroy(PluginIns) error
	Reload(PluginIns, map[string]interface{}) error
	CanUnload(PluginIns) bool
}

type config struct {
	Plugin map[string]map[string]map[string]interface{} `toml:"Plugin"`
}

type Type string

const (
	DB     Type = "db"     // 存储
	Router Type = "router" // 服务发现
	Trace  Type = "trace"  // 链路追踪
	Metric Type = "metric" // 监控
	Broker Type = "broker" // 消息队列
	Log    Type = "log"    // 日志
	Shape  Type = "shape"  // 流量管理
	Gate   Type = "gate"   // CS连接
	Lease  Type = "lease"  // 租约
	Lua    Type = "lua"    // Lua
	Trans  Type = "trans"  // 内部通信
	Grpc   Type = "grpc"   // 内部通信
)

var (
	_pluginserial = []Type{Log, Trans, Broker, Router, Trace, Metric, DB, Lease, Shape, Gate, Lua, Grpc} // Plugin加载顺序
	_factoryMap   = make(map[string]Factory)
)

func init() {
	RegisterFactory(grpc.NewFactory())
}

func RegisterFactory(f Factory) {
	if _, ok := _factoryMap[f.Name()]; ok {
		log.Fatal("RegisterFactory factory name repeatad %s", f.Name())
	}
	_factoryMap[f.Name()] = f
}

func InitPlugins(ctx context.Context) error {
	//var cfg config
	//if err := v.Unmarshal(&cfg); err != nil {
	//	//return errs.PluginsInitFail.Wrap(err, "Unmarshal PluginCfg")
	//	return errors.New(fmt.Sprintf("Unmarshal PluginCfg, error:%s", err))
	//}
	//for _, typ := range _pluginserial {
	//	for name, nm := range cfg.Plugin[string(typ)] {
	//		factory, ok := _factoryMap[name]
	//		if !ok {
	//			return errors.New(fmt.Sprintf("Factory not find %s", name))
	//			//return errs.PluginsFactoryNotRegister.New("Factory not find %s", name)
	//		}
	//		//log.Info("Plugin Setup %s", name)
	//		pluginIns, err := factory.Setup(ctx, nm)
	//		if err != nil {
	//			//return errs.PluginsSetupFail.Wrap(err, "Factory %s", name)
	//			return errors.New(fmt.Sprintf("Factory %s", name))
	//		}
	//		if err := addPluginIns(factory.Type(), name, pluginIns); err != nil {
	//			//return errs.PluginsAddPluginFail.Wrap(err, "Factory %s", name)
	//			return errors.New(fmt.Sprintf("Factory %s", name))
	//		}
	//	}
	//}
	//return nil

	for name, factory := range _factoryMap {
		//log.Info("Plugin Setup %s", name)
		pluginIns, err := factory.Setup(ctx)
		if err != nil {
			//return errs.PluginsSetupFail.Wrap(err, "Factory %s", name)
			return errors.New(fmt.Sprintf("Factory %s", name))
		}
		if err := addPluginIns(factory.Type(), name, pluginIns); err != nil {
			//return errs.PluginsAddPluginFail.Wrap(err, "Factory %s", name)
			return errors.New(fmt.Sprintf("Factory %s", name))
		}
	}

	return nil
}
