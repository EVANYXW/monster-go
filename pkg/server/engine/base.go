package engine

import (
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/output"
	"os"
)

type BaseEngine struct {
	serverName string
	isPprof    bool
	isOutput   bool
	output     *output.Output
}

type Options func(opt *option)

type option struct {
	isPprof  bool
	isOutput bool
	output   *output.Config
	modules  map[int32]module.IModule
}

func WithPprof() Options {
	return func(opt *option) {
		opt.isPprof = true
	}
}

func WithOutput(config *output.Config) Options {
	return func(opt *option) {
		opt.isOutput = true
		opt.output = config
	}
}

func WithModule(id int32, m module.IModule) Options {
	return func(opt *option) {
		//opt.modules = append(opt.modules, module)
		if opt.modules == nil {
			opt.modules = make(map[int32]module.IModule)
		}
		opt.modules[id] = m
	}
}

func NewEngine(name string, options ...Options) *BaseEngine {
	opt := &option{}
	for _, fn := range options {
		fn(opt)
	}

	b := &BaseEngine{
		serverName: name,
		isOutput:   opt.isOutput,
		isPprof:    opt.isPprof,
	}

	if opt.isOutput {
		b.output = output.NewOutput(opt.output)
	}

	for id, m := range opt.modules {
		module.NewBaseModule(id, m)
	}

	return b
}

func (b *BaseEngine) Run() {
	if b.output != nil {
		async.Go(func() {
			b.output.Run()
		})
	}

	module.Run()
}

func (b *BaseEngine) Destroy() {
	module.Close()
}

func (b *BaseEngine) OnSystemSignal(signal os.Signal) bool {
	return BaseSystemSignal(signal, b.serverName)
}
