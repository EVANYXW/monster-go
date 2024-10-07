package middlewaree

import (
	"context"
	"github.com/evanyxw/basic/logs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

type Option func(*option)
type option struct {
	logger *zap.Logger
}

// WithLogger zap.Logger
func WithLogger(logger *zap.Logger) Option {
	return func(opt *option) {
		opt.logger = logger
	}
}

// 服务端拦截器
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	//opt := &option{}
	//for _, f := range opts {
	//	f(opt)
	//}
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	reqIds := md["reqid"]
	reqId := ""
	if len(reqIds) > 0 {
		reqId = reqIds[0]
	}

	logLevel := false
	loglevel := md["middlewareloglevel"]
	if len(loglevel) > 0 {
		if loglevel[0] == "true" {
			logLevel = true
		}
	}

	resp, err := handler(ctx, req)
	if logLevel {
		log := logs.Log.WithField("reqId", reqId)
		log = log.WithField("duration", time.Since(start).String())
		log = log.WithField("method", info.FullMethod)
		log = log.WithField("req", req)
		log = log.WithField("resp", resp)
		if logs.Log.ServerName != "" {
			log = log.WithField("servername", logs.Log.ServerName)
		}
		if err != nil {
			log = log.WithField("err", err.Error())
		}
		log.Info("server after handling.")

		//opt.logger.With(zap.Field{Key: "reqId", String: reqId})
		//opt.logger.With(zap.Field{Key: "duration", String: time.Since(start).String()})
		//opt.logger.With(zap.Field{Key: "method", String: info.FullMethod})
		//opt.logger.With(zap.Field{Key: "req", Interface: req})
		//opt.logger.With(zap.Field{Key: "resp", Interface: resp})
		//if logs.Log.ServerName != "" {
		//	opt.logger.With(zap.Field{Key: "servername", String: logs.Log.ServerName})
		//}
		//if err != nil {
		//	opt.logger.With(zap.Field{Key: "err", String: err.Error()})
		//}
		//opt.logger.Info("server after handling.")
	}
	return resp, err
}

// interceptor 客户端拦截器
func Interceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	//opt := &option{}
	//for _, f := range opts {
	//	f(opt)
	//}
	start := time.Now()
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	reqIds := md["reqid"]
	reqId := ""
	if len(reqIds) > 0 {
		reqId = reqIds[0]
	}

	logLevel := false
	loglevel := md["middlewareloglevel"]
	if len(loglevel) > 0 {
		if loglevel[0] == "true" {
			logLevel = true
		}
	}

	err := invoker(ctx, method, req, reply, cc, opts...)
	if logLevel {
		log := logs.Log.WithField("reqId", reqId)
		log = log.WithField("duration", time.Since(start).String())
		log = log.WithField("method", method)
		log = log.WithField("req", req)
		if logs.Log.ServerName != "" {
			log = log.WithField("servername", logs.Log.ServerName)
		}
		log = log.WithField("resp", reply)
		if err != nil {
			log = log.WithField("err", err.Error())
		}
		log.Info("client after handling.")

		//opt.logger.With(zap.Field{Key: "reqId", String: reqId})
		//opt.logger.With(zap.Field{Key: "duration", String: time.Since(start).String()})
		//opt.logger.With(zap.Field{Key: "method", String: method})
		//opt.logger.With(zap.Field{Key: "req", Interface: req})
		//
		//if logs.Log.ServerName != "" {
		//	opt.logger.With(zap.Field{Key: "servername", String: logs.Log.ServerName})
		//}
		//opt.logger.With(zap.Field{Key: "resp", Interface: reply})
		//if err != nil {
		//	opt.logger.With(zap.Field{Key: "err", String: err.Error()})
		//}
		//opt.logger.Info("client after handling.")
	}
	return err
}
