package module

import (
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/golang/protobuf/proto"
)

// 服务器信息变化处理
type IServerInfoHandler interface {
	OnServerNew(Info *network.ServerInfo) // 有一个服务器连入到集群
	OnServerLost(id uint32)               // 有一个服务器断开
	OnServerOk(Info *network.ServerInfo)  // 服务器已准备好
	OnServerOpenComplete()
}

// IModuleFlow Module流程
type IModuleFlow interface {
	Init(baseModule *BaseModule) bool
	DoRegister()
	DoRun()
	DoWaitStart()
	DoRelease()
	OnOk()
	OnStartCheck() int
	OnCloseCheck() int
	Update()
}

// IModule Module
type IModule interface {
	IModuleFlow
	GetKernel() IModuleKernel
	GetID() int32
}

// IModuleKernel 模块内核定义、继承该接口将成为 Module
type IModuleKernel interface {
	IModuleFlow
	GetNoWaitStart() bool
	OnStartClose()
	DoClose()
	GetNPManager() network.INPManager
	GetStatus() int
}

type Client interface {
	SendMessage(message proto.Message)
	SetSignal(data []byte)
	GetID() uint32
	GetLastHeartbeat() *uint64
	Close()
	GoDisconnect(id uint32)
	SetServerID(args []interface{})
	GetServerIds() []uint32
}

// IGtClientManager  gate的client管理器
type IGtClientManager interface {
	NewClient(np *network.NetPoint) (Client, bool)
	GetClient(id uint32) Client
}

// IGtAClientManager gate消息接受管理器
type IGtAClientManager interface {
	CloseClient(id uint32)
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
	ModuleID_GateAcceptor
	ModuleID_ConnectorManager
	ModuleID_LoginManager
	ModuleID_LoginConfig
	ModuleID_Redis

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
