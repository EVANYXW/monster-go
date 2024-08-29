package module

import (
	"github.com/evanyxw/monster-go/internal/servers/login/handler"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
)

type LoginManager struct {
	kernel module.IModuleKernel
}

func NewLoginManager(id int32) *LoginManager {
	h := handler.NewLoginMsgHandler()
	l := &LoginManager{
		kernel: module.NewKernel(h,
			network.NetPointManager.GetRpcAcceptor(),
			network.NetPointManager.GetProcessor(),
		),
	}

	baseModule := module.NewBaseModule(id, l)
	h.Init(baseModule) //fixMe 这个看能否改为kernel 里去调用
	return l
}

func (l *LoginManager) Init() bool {
	l.kernel.Init()
	return true
}

func (l *LoginManager) DoRun() {
	l.kernel.DoRun()
}

func (l *LoginManager) DoWaitStart() {
	l.kernel.DoWaitStart()
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

func (l *LoginManager) DoRegister() {
	l.kernel.DoRegister()
}

func (l *LoginManager) DoRelease() {
	l.kernel.DoRelease()
}

func (l *LoginManager) OnNetError(np *network.NetPoint) {
	logger.Debug("center onNetError")
	//l.nodeManager.OnNodeLost(np.ID, np.SID.Type)
}

func (l *LoginManager) OnServerOk() {

}

func (l *LoginManager) OnNPAdd(np *network.NetPoint) {

}
