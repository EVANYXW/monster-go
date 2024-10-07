// Package transport @Author evan_yxw
// @Date 2024/10/5 13:12:00
// @Desc
package transport

import (
	"errors"
	"fmt"
	"github.com/evanyxw/monster-go/pkg/plugins"
	"google.golang.org/grpc"
	"log"
)

var (
	_transport Transport
)

type Transport interface {
	plugins.PluginIns
	//Host() string                                                                        // 通讯地址
	//Port() int                                                                           // 通讯地址
	//Prepare(ctx context.Context) error                                                   // 启动前准备
	//Start(ctx context.Context) error                                                     // 启动
	//RegisterSer(f func(grpc.ServiceRegistrar) error) error                               // 注册业务Ser
	//RegisterServerInterceptor(f grpc.UnaryServerInterceptor)                             // 注册Server拦截器
	//RegisterClientInterceptor(f grpc.UnaryClientInterceptor)                             // 注册Client拦截器
	ClientConn(target string, opts ...grpc.DialOption) (grpc.ClientConnInterface, error) // 获取Client对象
}

func New(trans Transport) (Transport, error) {
	if _transport != nil {
		log.Fatal("Trans has New")
		//return nil, errs.TransportInsHasNewed.New("Trans is newed %s", trans.Factory().Name())
		return nil, errors.New(fmt.Sprintf("Trans is newed %s", trans.Factory().Name()))
	}
	_transport = trans
	return _transport, nil
}

func Instance() Transport {
	if _transport == nil {
		log.Fatal("Trans not newd")
		return nil
	}
	return _transport
}
