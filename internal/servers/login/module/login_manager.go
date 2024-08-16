package module

import (
	"github.com/evanyxw/monster-go/internal/servers"
	"github.com/evanyxw/monster-go/internal/servers/login/handler"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
)

type LoginManager struct {
	*module.BaseModule
	kernel *module.Kernel
	ID     int32
}

func NewLoginManager(id int32) *LoginManager {
	l := &LoginManager{
		ID: id,
		kernel: module.NewKernel(handler.NewLoginMsgHandler(false),
			servers.NetPointManager.GetRpcAcceptor(),
			servers.NetPointManager.GetProcessor(),
		),
	}

	baseModule := module.NewBaseModule(l)
	l.BaseModule = baseModule
	return l
}

// 外部通知开启Module
//func (l *LoginManager) Run() {
//	l.BaseModule.Run()
//}

func (l *LoginManager) Init() {
	l.kernel.Init()
}

func (l *LoginManager) DoRun() {
	l.kernel.Start()
}

func (l *LoginManager) DoWaitStart() {
	l.kernel.DoStart()
}

func (l *LoginManager) DoRelease() {
	l.kernel.Release()
}

func (l *LoginManager) OnOk() {

}

func (l *LoginManager) OnStartCheck() int {
	return module.ModuleRunCode_Ok
}

func (l *LoginManager) OnCloseCheck() int {
	return l.kernel.OnCloseCheck()
}

func (l *LoginManager) GetKernel() module.IModuleKernel {
	return l.kernel
}

func (l *LoginManager) Update() {

}

func (l *LoginManager) GetID() int32 {
	return l.ID
}

func (l *LoginManager) DoRegister() {
	l.kernel.DoRegist()
}

func (l *LoginManager) Release() {
	l.kernel.Release()
}

func (l *LoginManager) OnNetError(np *network.NetPoint) {
	logger.Debug("center onNetError")
	//l.nodeManager.OnNodeLost(np.ID, np.SID.Type)
}

func (l *LoginManager) OnServerOk() {

}

func (l *LoginManager) OnNPAdd(np *network.NetPoint) {

}
