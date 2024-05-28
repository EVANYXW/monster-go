package module

import (
	"github.com/evanyxw/monster-go/pkg/async"
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
	id          int32
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
		id:          owner.GetID(),
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

func (m *BaseModule) GetID() int32 {
	return m.id
}

func (m *BaseModule) GetOwner() IModule {
	return m.owner
}

func (m *BaseModule) Close() {
	m.closeSig <- true
}

func (m *BaseModule) DoStart() {
	//xsf_log.Info("module onDoStart", xsf_log.String("name", m.Name))
	if m.runStatus == ModuleRunStatus_WaitStart {
		//xsf_log.Info("module call start", xsf_log.String("name", m.Name))
		m.owner.DoStart()
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

func (m *BaseModule) setStartOK() {
	m.SetStatusToRunning()
}

func (m *BaseModule) onCloseCheck() int {
	return m.owner.OnCloseCheck()
}

func (m *BaseModule) setCloseOK() {

}

func (m *BaseModule) release() {
	m.owner.DoRelease()
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
			default:

			}

			switch m.runStatus {
			case ModuleRunStatus_WaitStart:
			case ModuleRunStatus_Start:
				res := m.onStartCheck()
				if res == ModuleRunCode_Ok {
					m.runStatus = ModuleRunStatus_WaitOK
					m.setStartOK()
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

		}
	})
}
