// Package module @Author evan_yxw
// @Date 2024/10/16 14:42:00
// @Desc
package kernel

import handler2 "github.com/evanyxw/monster-go/pkg/handler"

type options struct {
	handler        handler2.Handler
	msgHandlerImpl handler2.MsgHandler
}
type Option func(options *options)

func WithHandler(handler handler2.Handler) Option {
	return func(options *options) {
		options.handler = handler
	}
}

func WithMsgHandler(handler handler2.MsgHandler) Option {
	return func(options *options) {
		options.msgHandlerImpl = handler
	}
}

func NewOption() *options {
	return &options{}
}

func (g *options) GetHandler() handler2.Handler {
	return g.handler
}

func (g *options) GetHandlerImpl() handler2.MsgHandler {
	return g.msgHandlerImpl
}
