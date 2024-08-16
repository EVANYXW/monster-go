package engine

import (
	"github.com/evanyxw/monster-go/pkg/module"
	"os"
)

type BaseEngine struct {
	serverName string
}

func NewEngine(name string) *BaseEngine {
	return &BaseEngine{
		serverName: name,
	}
}

func (b *BaseEngine) Run() {
	module.Run()
}

func (b *BaseEngine) Destroy() {
	module.Close()
}

func (b *BaseEngine) OnSystemSignal(signal os.Signal) bool {
	return BaseSystemSignal(signal, b.serverName)
}
