package network

import (
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/golang/protobuf/proto"
	"net"
)

var NetPointManager INPManager

type INPManager interface {
	//Init(module *Module, maxCount uint32)
	New(conn *net.TCPConn) *NetPoint
	Add(np *NetPoint) bool
	Del(np *NetPoint)
	Get(id uint32) *NetPoint
	GetProcessor() *Processor
	GetRpcAcceptor() *rpc.Acceptor
	GetMaxConnNum() uint32
	//Release()
	OnUpdate()
	OnHandshake(np *NetPoint) bool
	//OnCloseCheck() int
	Broadcast(msg proto.Message, skip uint32)
	//CountHandshake() int
}
