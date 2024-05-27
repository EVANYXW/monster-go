package factory

import (
	"github.com/evanyxw/monster-go/pkg/server"
	"os"
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

// CmdServer cmd的server
type CmdServer interface {
	Run()
	Destroy()
	OnSystemSignal(signal os.Signal) bool
}

type ServerNewFunc func(info server.Info) CmdServer

var typeRegistry = make(map[string]ServerNewFunc)

func Register(name string, fn ServerNewFunc) {
	typeRegistry[name] = fn
}

func MakeInstance(name string) ServerNewFunc {
	if fn, ok := typeRegistry[name]; ok {
		return fn
	}
	return nil
}
