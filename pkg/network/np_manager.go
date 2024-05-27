package network

import (
	"github.com/golang/protobuf/proto"
	"net"
)

type INPManager interface {
	//Init(module *Module, maxCount uint32)
	New(conn *net.TCPConn) *NetPoint
	Add(np *NetPoint) bool
	Del(np *NetPoint)
	Get(id uint32) *NetPoint
	//Release()
	OnUpdate()
	OnHandshake(np *NetPoint) bool
	//OnCloseCheck() int
	Broadcast(msgId int32, msg proto.Message, skip uint32)
	//CountHandshake() int
}
