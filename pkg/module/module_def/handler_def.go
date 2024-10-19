package module_def

import (
	"sync/atomic"
)

var (
	MailID    atomic.Uint32
	ManagerID atomic.Uint32
)

// KernelNetEvent kernel网络事件处理器
type KernelNetEvent interface {
	OnRpcNetAccept(args []interface{})
	OnRpcNetConnected(args []interface{})
	OnRpcNetError(args []interface{})
	OnRpcNetClose(args []interface{})
	OnRpcNetData(args []interface{})
	OnRpcNetMessage(args []interface{})
}
