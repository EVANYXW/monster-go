package module

import (
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/output"
	"go.uber.org/zap"
	"sync"
	"time"
)

type BaseModule struct {
	owner       IModule
	name        string
	runStatus   int
	timeStart   int64
	okSig       chan bool
	closeSig    chan bool
	wg          sync.WaitGroup
	ID          int32
	NoWaitStart bool
}

func NewBaseModule(owner IModule) *BaseModule {
	noWaitStart := false
	if owner.GetKernel() != nil {
		noWaitStart = owner.GetKernel().GetNoWaitStart()
	}

	b := &BaseModule{
		owner:       owner,
		okSig:       make(chan bool),
		closeSig:    make(chan bool),
		ID:          owner.GetID(),
		name:        ModuleId2Name(int(owner.GetID())),
		NoWaitStart: noWaitStart,
	}

	AddModule(b)
	return b
}

func (m *BaseModule) Init() {
	if m.NoWaitStart {
		m.runStatus = ModuleRunStatus_Start
	} else {
		m.runStatus = ModuleRunStatus_WaitStart
	}

	m.owner.Init()
}

func (m *BaseModule) DoRegister() {
	m.owner.DoRegister()
}

func (m *BaseModule) GetID() int32 {
	return m.ID
}

func (m *BaseModule) GetOwner() IModule {
	return m.owner
}

func (m *BaseModule) onOK() {
	m.owner.OnOk()
}

func (m *BaseModule) Close() {
	m.closeSig <- true
}

// DoWaitStart 在C_Cc_Handshake 后会调用
func (m *BaseModule) DoWaitStart() {
	//xsf_log.Info("module onDoStart", xsf_log.String("name", m.Name))
	if m.runStatus == ModuleRunStatus_WaitStart {
		//xsf_log.Info("module call start", xsf_log.String("name", m.Name))
		m.owner.DoWaitStart()
		m.runStatus = ModuleRunStatus_Start
		//xsf_log.Info("module call start done", xsf_log.String("name", m.Name))
	}
}

func (m *BaseModule) SetStatusToRunning() {
	m.runStatus = ModuleRunStatus_Running
}

func (m *BaseModule) update() {
	m.owner.Update()
}

func (m *BaseModule) onStartCheck() int {
	return m.owner.OnStartCheck()
}

func (m *BaseModule) setStartOK(id int32) {
	for i := 0; i < len(modules); i++ {
		module := modules[i].module
		if module != nil && module.ID == id {
			modules[i].isStartOk.Swap(true)
			startOkCount.Add(1)
			output.Oput.SetModuleNum(total, int(startOkCount.Load()))
			logger.Info("===> module start ok", zap.String("name", m.name), zap.Int("count", int(startOkCount.Load())), zap.Int("total", total))
		}
	}
	//m.SetStatusToRunning()
}

func (m *BaseModule) onCloseCheck() int {
	return m.owner.OnCloseCheck()
}

func (m *BaseModule) setCloseOK() {

}

func (m *BaseModule) release() {
	m.owner.DoRelease()
}

func check() {
	for {
		if startOkCount.Load() >= int32(total) {
			logger.Info("================== all server module ok ==================")
			Status.Swap(ModuleRunStatus_Running)
			for i := 0; i < len(modules); i++ {
				module := modules[i].module
				if module != nil {
					module.okSig <- true
				}
			}
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (m *BaseModule) Run() {
	if m.NoWaitStart {
		m.owner.DoRun()
	}

	async.Go(func() {
		for {
			select {
			case <-m.closeSig:
				m.runStatus = ModuleRunStatus_WaitStop
			case <-m.okSig:
				m.onOK()
				m.runStatus = ModuleRunStatus_Running
			default:

			}

			switch m.runStatus {
			case ModuleRunStatus_WaitStart:
			case ModuleRunStatus_Start:
				res := m.onStartCheck()
				if res == ModuleRunCode_Ok {
					m.runStatus = ModuleRunStatus_WaitOK
					m.setStartOK(m.ID)
				}
			case ModuleRunStatus_Running:
				curTime := time.Now().Unix()
				timeDelta := curTime - m.timeStart
				m.timeStart = curTime
				if timeDelta > 0 {
					m.update()
				}
			case ModuleRunStatus_WaitStop:
				res := m.onCloseCheck()
				if res == ModuleRunCode_Ok {
					m.runStatus = ModuleRunStatus_Stop
					m.setCloseOK()
				}
			case ModuleRunStatus_Stop:
				m.release()
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	})
}
