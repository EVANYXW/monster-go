package module

import (
	"github.com/evanyxw/monster-go/internal/servers/login/config"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module/module_def"
	"github.com/evanyxw/monster-go/pkg/network"
)

type LoginConfig struct {
	kernel module_def.IKernel
	id     int32
}

func NewLoginConfig() *LoginConfig {
	//h := handler.NewCommonMsgHandler()
	id := module_def.GetModuleId(module_def.ModuleLoginConfig)
	l := &LoginConfig{
		id:     id,
		kernel: config.New(),
	}

	//module.NewBaseModule(id, l)
	//h.Init(baseModule) //fixMe 这个看能否改为kernel 里去调用
	return l
}

func (l *LoginConfig) Init(baseModule module_def.IBaseModule) bool {
	l.kernel.Init(baseModule)
	return true
}

func (l *LoginConfig) DoRun() {
	l.kernel.DoRun()
}

func (l *LoginConfig) DoWaitStart() {
	l.kernel.DoWaitStart()
}

func (l *LoginConfig) OnOk() {
	l.kernel.OnOk()
}

func (l *LoginConfig) OnStartCheck() int {
	return module_def.ModuleOk()
}

func (l *LoginConfig) OnCloseCheck() int {
	return l.kernel.OnCloseCheck()
}

func (l *LoginConfig) GetID() int32 {
	return l.id
}

func (l *LoginConfig) GetKernel() module_def.IKernel {
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
