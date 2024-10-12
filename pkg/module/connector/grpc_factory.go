// Package connector @Author evan_yxw
// @Date 2024/9/28 14:12:00
// @Desc
package connector

//import (
//	commonHandler "github.com/evanyxw/monster-go/pkg/common/handler"
//	"github.com/evanyxw/monster-go/pkg/module"
//	"github.com/evanyxw/monster-go/pkg/network"
//)
//
//// GrpcManagerFactory grpc 管理器
//type GrpcManagerFactory struct {
//}
//
//func (t *GrpcManagerFactory) Create(id int32) *Manager {
//	c := &Manager{
//		id:      id,
//		factory: t,
//	}
//	hdler := commonHandler.NewCommonMsgHandler()
//	c.handler = hdler
//	c.kernel = module.NewKernel(hdler, network.NetPointManager.GetRpcAcceptor(),
//		network.NetPointManager.GetProcessor())
//	return c
//}
//
//func (t *GrpcManagerFactory) CreateConnector(handler module.GateAcceptorHandler, id uint32, ip string, port uint32) *module.ConnectorKernel {
//
//	return nil
//}
