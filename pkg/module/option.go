// Package module @Author evan_yxw
// @Date 2024/10/16 14:42:00
// @Desc
package module

type options struct {
	handler        Handler
	msgHandlerImpl MsgHandler
}
type Option func(options *options)

func WithHandler(handler Handler) Option {
	return func(options *options) {
		options.handler = handler
	}
}

func WithMsgHandler(handler MsgHandler) Option {
	return func(options *options) {
		options.msgHandlerImpl = handler
	}
}
