package module

import (
	"github.com/evanyxw/monster-go/internal/servers/common/handler"
	"github.com/evanyxw/monster-go/internal/servers/login/config"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
)

type LoginConfig struct {
	kernel module.IModuleKernel
}

func NewLoginConfig(id int32) *LoginConfig {
	h := handler.NewCommonMsgHandler()
	l := &LoginConfig{
		kernel: config.New(h),
	}

	baseModule := module.NewBaseModule(id, l)
	h.Init(baseModule) //fixMe 这个看能否改为kernel 里去调用
	return l
}

func (l *LoginConfig) Init() bool {
	l.kernel.Init()
	return true
}

func (l *LoginConfig) DoRun() {
	l.kernel.DoRun()
}

func (l *LoginConfig) DoWaitStart() {
	l.kernel.DoWaitStart()
}

func (l *LoginConfig) OnOk() {

}

func (l *LoginConfig) OnStartCheck() int {
	return module.ModuleRunCode_Ok
}

func (l *LoginConfig) OnCloseCheck() int {
	return l.kernel.OnCloseCheck()
}

func (l *LoginConfig) GetKernel() module.IModuleKernel {
	return l.kernel
}

func (l *LoginConfig) Update() {

}

func (l *LoginConfig) DoRegister() {
	l.kernel.DoRegister()
}

func (l *LoginConfig) DoRelease() {
	l.kernel.DoRelease()
}

func (l *LoginConfig) OnNetError(np *network.NetPoint) {
	logger.Debug("center onNetError")
	//l.nodeManager.OnNodeLost(np.ID, np.SID.Type)
}

func (l *LoginConfig) OnServerOk() {

}

func (l *LoginConfig) OnNPAdd(np *network.NetPoint) {

}
