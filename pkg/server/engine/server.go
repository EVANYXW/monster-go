package engine

import (
	"fmt"
	"github.com/evanyxw/monster-go/pkg/logger"
	"go.uber.org/zap"
	"os"
	"syscall"
)

// 方案一：
// var typeRegistry = make(map[string]reflect.Type)
//
//	func init() {
//		typeRegistry["world"] = reflect.TypeOf(world.World{})
//
// }
//
//	func MakeInstance(name string) interface{} {
//		v := reflect.New(typeRegistry[name]).Elem()
//
//		return v.Interface()
//	}
//

// IServerKernel Server内核
type IServerKernel interface {
	Run()
	Destroy()
	OnSystemSignal(signal os.Signal) bool
}

type KernelFun func() IServerKernel

var serverEngines = make(map[string]KernelFun)

func Register(name string, fn KernelFun) {
	serverEngines[name] = fn
}

func MakeInstance(name string) KernelFun {
	if fn, ok := serverEngines[name]; ok {
		return fn
	}
	return nil
}

func BaseSystemSignal(signal os.Signal, serverName string) bool {
	tag := true
	switch signal {
	case syscall.SIGHUP:
		//todo
		fmt.Println("SIGHUP")
	case syscall.SIGPIPE:
		fmt.Println("SIGPIPE")
	default:
		logger.Info(fmt.Sprintf("【 %s 】 收到信号准备退出", serverName),
			zap.String("signal", signal.String()))
		tag = false
	}
	return tag
}
