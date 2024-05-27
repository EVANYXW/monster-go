package module

import (
	"github.com/evanyxw/monster-go/pkg/network"
)

// 服务器信息变化处理
type IServerInfoHandler interface {
	OnServerNew(Info *network.ServerInfo) // 有一个服务器连入到集群
	OnServerLost(id uint32)               // 有一个服务器断开
	OnServerOk(Info *network.ServerInfo)  // 服务器已准备好
	OnServerOpenComplete()
}

// IModule 定义
type IModule interface {
	GetID() int32

	Init()
	DoRun()
	DoStart()
	DoRelease()
	OnStartCheck() int
	OnCloseCheck() int
	Update()
}

// 模块内核定义、继承该接口将成为 Module
type IModuleKernel interface {
	Init() bool
	AddModules()
	DoRegist()
	Start()
	DoStart()
	Release()
	OnUpdate(timeDelta uint32)
	OnOK()
	OnStartClose()
	DoClose()
	OnStartCheck() int
	OnCloseCheck() int
}

// INetHandler kernel 需要实现
type INetHandler interface {
	//Start()
	//DoRegist(nk *NetKernel)
	//OnNetData(np *network.NetPoint, msgID uint16, rawID uint32, data []byte)
	//OnNetMessage(np *network.NetPoint, msgID uint16, rawID uint32, message xsf_net.IMessage)
	//OnStartCheck() int
	//OnCloseCheck() int
	//DoClose()
	//OnStartClose()
	OnNetError(np *network.NetPoint)
	OnServerOk()
	OnNPAdd(np *network.NetPoint)
}

// 网络事件处理器
type INetEventHandler interface {
	OnRpcNetAccept(args []interface{})
	OnRpcNetConnected(args []interface{})
	OnRpcNetError(args []interface{})
	OnRpcNetData(args []interface{})
	OnRpcNetMessage(args []interface{})
}

const (
	ModuleRunCode_Ok = iota
	ModuleRunCode_Wait
	ModuleRunCode_Error
)

const (
	ModuleRunStatus_None = iota
	ModuleRunStatus_Running
	ModuleRunStatus_WaitStart
	ModuleRunStatus_Start

	ModuleRunStatus_Stop
	ModuleRunStatus_WaitStop
	ModuleRunStatus_WaitOK
)

const (
	ModuleID_Schema = iota
	ModuleID_SM     // 通知类型
	ModuleID_CenterConnector
	ModuleID_Client
	ModuleID_ConnectorManager

	ModuleID_Notice
	ModuleID_Pprof

	ModuleID_CC
	ModuleID_Nats
	ModuleID_Mongo
	ModuleID_Actor // 玩家模块放在最后
	ModuleID_Purchase

	ModuleID_Max
)

var moduleMap = map[int]string{
	ModuleID_SM:               "Center 网络模块",
	ModuleID_CenterConnector:  "Center 的连接器",
	ModuleID_Client:           "Client 客户端",
	ModuleID_ConnectorManager: "Gate 的多链接管理器",
}

func ModuleId2Name(id int) string {
	if val, ok := moduleMap[id]; ok {
		return val
	}
	return ""
}
