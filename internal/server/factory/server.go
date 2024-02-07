package factory

import (
	"bilibili/monster-go/internal/network"
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
// Server cmd的server
type Server interface {
	Start()
	Stop()
}

type ServerNewFunc func(info network.Info) Server

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
