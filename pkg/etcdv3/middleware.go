package etcdv3

import (
	"context"
	"github.com/evanyxw/monster-go/pkg/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

// 服务端拦截器
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
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
			log = log.WithField("serverName", logs.Log.ServerName)
		}
		if err != nil {
			log = log.WithField("err", err.Error())
		}
		log.Info("server after handling.")
	}
	return resp, err
}

func Interceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
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
			log = log.WithField("serverName", logs.Log.ServerName)
		}
		log = log.WithField("resp", reply)
		if err != nil {
			log = log.WithField("err", err.Error())
		}
		log.Info("client after handling.")
	}
	return err
}
