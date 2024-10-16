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

// IModuleEvent Module流程
type IModuleEvent interface {
	Init(baseModule IBaseModule) bool
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
	IModuleEvent
	GetKernel() IKernel
	GetID() int32
}

// IKernel 模块内核定义、继承该接口将成为 Module
type IKernel interface {
	IModuleEvent
	GetNoWaitStart() bool
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

func ModuleOk() int {
	return ModuleRunCode_Ok
}

func ModuleWait() int {
	return ModuleRunCode_Wait
}

const (
	ModuleRunStatus_Running = iota + 1
	ModuleRunStatus_WaitStart
	ModuleRunStatus_Start

	ModuleRunStatus_Stop
	ModuleRunStatus_WaitStop
	ModuleRunStatus_WaitOK
)

const (
	ModuleID_SM = iota + 1 // 通知类型
	ModuleID_CenterConnector
	ModuleID_Client
	ModuleID_GateAcceptor
	ModuleID_ConnectorManager
	ModuleID_Etcd

	// 外部模块
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

const (
	ModuleCenterServer         = "CenterServer"
	ModuleCenterConnector      = "CenterConnector"
	ModuleClient               = "Client"
	ModuleGateAcceptor         = "GateAcceptor"
	ModuleGateConnectorManager = "GateConnectorManager"
	ModuleLoginManager         = "LoginManager"
	ModuleLoginConfig          = "LoginConfig"
	ModuleRedis                = "Redis"
	ModuleEtcd                 = "Etcd"
)

var moduleMap = map[int32]string{
	ModuleID_SM:               ModuleCenterServer,
	ModuleID_CenterConnector:  ModuleCenterConnector,
	ModuleID_Client:           ModuleClient,
	ModuleID_GateAcceptor:     ModuleGateAcceptor,
	ModuleID_ConnectorManager: ModuleGateConnectorManager,
	ModuleID_LoginManager:     ModuleLoginManager,
	ModuleID_LoginConfig:      ModuleLoginConfig,
	ModuleID_Redis:            ModuleRedis,
	ModuleID_Etcd:             ModuleEtcd,
}

func GetModuleId(name string) int32 {
	for id, moduleName := range moduleMap {
		if name == moduleName {
			return id
		}
	}
	moduleMap[int32(len(moduleMap))] = name
	return int32(len(moduleMap))
}

func ModuleIdToName(id int32) string {
	if val, ok := moduleMap[id]; ok {
		return val
	}
	return "未知"
}

func GetModuleMap() map[int32]string {
	return moduleMap
}
