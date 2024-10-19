package module

import (
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/output"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"go.uber.org/zap"
	"sync"
	"time"
)

type IBaseModule interface {
	GetRpcAcceptor() *rpc.Acceptor
}

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
	RpcAcceptor *rpc.Acceptor
}

func NewBaseModule(id int32, owner IModule) *BaseModule {
	noWaitStart := false
	if owner.GetKernel() != nil {
		noWaitStart = owner.GetKernel().GetNoWaitStart()
	}

	b := &BaseModule{
		owner:       owner,
		okSig:       make(chan bool),
		closeSig:    make(chan bool),
		ID:          id,
		name:        ModuleIdToName(id),
		NoWaitStart: noWaitStart,
		RpcAcceptor: rpc.NewAcceptor(1000),
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

	m.owner.Init(m)
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
			output.Oput.SetModuleNum(total, int(startOkCount.Load()), id)
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

func (m *BaseModule) GetRpcAcceptor() *rpc.Acceptor {
	return m.RpcAcceptor
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

func (m *BaseModule) RpcEventRun() {
OUTLABEL:
	for {
		select {
		case callMsg, ok := <-m.RpcAcceptor.ChanCall:
			if !ok {
				break OUTLABEL
			}
			if callMsg == nil {
				continue
			}
			m.RpcAcceptor.Execute(callMsg)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (m *BaseModule) Run() {
	if m.NoWaitStart {
		m.owner.DoRun()
	}

	var allModules []int32
	for _, v := range modules {
		if v.module == nil {
			continue
		}
		allModules = append(allModules, v.module.GetID())
	}

	output.Oput.SetAllModules(allModules)
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

			// 1.ModuleRunStatus_Start 和 ModuleRunStatus_WaitStart 需要检测module里的onStartCheck 是否ok
			// 2.如果不ok继续等待,如果ok,状态流转至ModuleRunStatus_WaitOK,进入缓冲等待
			// 3.缓冲module数量与所有数量一致了,所有module状态流转至 ModuleRunStatus_Running
			switch m.runStatus {
			case ModuleRunStatus_WaitStart:
			case ModuleRunStatus_Start:
				res := m.onStartCheck()
				if res == ModuleOk() {
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
				if res == ModuleOk() {
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
