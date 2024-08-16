package module

import (
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/network"
	"sync/atomic"
)

type moduleNode struct {
	module    *BaseModule
	manager   network.INPManager
	isStartOk atomic.Bool
	isCloseOk atomic.Bool
}

var (
	Status       atomic.Int32
	modules      []moduleNode
	ModuleMax    int
	total        int
	startOkCount atomic.Int32
	closeOkCount atomic.Int32
)

func GetModuleById(id int32) moduleNode {
	return modules[id]
}

func Init(moduleMax int) {
	ModuleMax = moduleMax
	modules = make([]moduleNode, ModuleMax)
}

func AddModule(m *BaseModule) {
	if modules == nil {
		modules = make([]moduleNode, ModuleMax)
	}

	if modules[m.GetID()].module != nil {
		//xsf_log.Panic(fmt.Sprintf("AddModule not nil, id=%d, name=%s", m.ID, m.Name))
		return
	}

	modules[m.GetID()].module = m
	//fmt.Println("modules", modules)
	//xsf_log.Info("AddModule", xsf_log.String("name", m.Name), xsf_log.Int("id", m.ID))
}

func AddManager(id int, m network.INPManager) {
	modules[id].manager = m
}

func GetModule(id int) *BaseModule {
	return modules[id].module
}

func GetManager(id int) network.INPManager {
	return modules[id].manager
}

func GetConnectorManager() IModule {
	return GetModule(ModuleID_ConnectorManager).owner
}

func setStartOK(id int) {
	for i := 0; i < len(modules); i++ {
		module := modules[i].module
		if module != nil && module.GetID() == int32(id) {
			// todo
			//modules[i].isStartOk.Set(true)
			startOkCount.Add(1)
			//xsf_log.Info("===> module start ok", xsf_log.String("name", module.Name), xsf_log.Int("count", int(startOkCount.Get())), xsf_log.Int("total", total))
		}
	}
}

func canStart(id int) bool {
	isOK := true
	for i := 0; i < len(modules); i++ {
		module := modules[i].module
		if module != nil {
			if module.GetID() == int32(id) {
				break
			}
			// todo
			//isOK = modules[i].isStartOk.Get()
		}
	}

	return isOK
}

func setCloseOK(id int) {
	for i := len(modules) - 1; i >= 0; i-- {
		module := modules[i].module
		if module != nil && module.GetID() == int32(id) {
			// todo
			//modules[i].isCloseOk.Set(true)
			closeOkCount.Add(1)
			//xsf_log.Info("===> module close ok", xsf_log.String("name", module.Name), xsf_log.Int("count", int(closeOkCount.Get())), xsf_log.Int("total", total))
		}
	}
}

func canClose(id int) bool {
	isOK := true
	for i := len(modules) - 1; i >= 0; i-- {
		module := modules[i].module
		if module != nil {
			if module.GetID() == int32(id) {
				break
			}
			// todo
			//isOK = modules[i].isCloseOk.Get()
		}
	}

	return isOK
}

func Run() {
	async.Go(func() {
		check()
	})

	for _, moduleNode := range modules {
		if moduleNode.module == nil {
			continue
		}
		total++
		moduleNode.module.Init()
	}

	for _, moduleNode := range modules {
		if moduleNode.module == nil {
			continue
		}
		moduleNode.module.DoRegister()
	}

	for _, moduleNode := range modules {
		if moduleNode.module == nil {
			continue
		}
		moduleNode.module.Run()
	}
}

func Close() {

	for _, moduleNode := range modules {
		if moduleNode.module == nil {
			continue
		}
		moduleNode.module.Close()
	}
}

// DoWaitStart 在C_Cc_Handshake 后会调用
func DoWaitStart() {
	for _, moduleNode := range modules {
		if moduleNode.module == nil {
			continue
		}
		moduleNode.module.DoWaitStart()
	}
}

//
//func isAllClose() bool {
//	return closeOkCount.Get() >= int32(total)
//}
//
//func Init() (err error) {
//	total = 0
//	length := len(modules)
//
//	xsf_server.Status.Set(xsf_server.ServerStatus_Starting)
//
//	for i := 0; i < length; i++ {
//		module := modules[i].module
//		if module == nil {
//			continue
//		}
//
//		if !module.init() {
//			xsf_log.Panic("Module Init return false, name=" + module.Name)
//			return
//		}
//
//		total++
//		xsf_log.Info("Module Init Done", xsf_log.String("name", module.Name), xsf_log.Int("id", module.ID))
//	}
//
//	for i := 0; i < length; i++ {
//		module := modules[i].module
//		if module == nil {
//			continue
//		}
//
//		module.doRegist()
//	}
//
//	for i := 0; i < length; i++ {
//		module := modules[i].module
//		if module == nil {
//			continue
//		}
//
//		if module.NoWaitStart {
//			module.start()
//		}
//	}
//
//	for i := 0; i < length; i++ {
//		module := modules[i].module
//		if module == nil {
//			continue
//		}
//
//		module.wg.Add(1)
//		xsf_pool.Submit(
//			func() {
//				run(module)
//			})
//	}
//
//	xsf_pool.Submit(check)
//
//	return nil
//}
//
//func Release() {
//	xsf_log.Info("modules Release call ...")
//	for i := len(modules) - 1; i >= 0; i-- {
//		module := modules[i].module
//		if module == nil {
//			continue
//		}
//
//		module.closeSig <- true
//	}
//
//	for i := len(modules) - 1; i >= 0; i-- {
//		module := modules[i].module
//		if module == nil {
//			continue
//		}
//
//		module.wg.Wait()
//	}
//
//	xsf_log.Info("modules Release ok")
//}
//
//func DoStart() {
//	xsf_log.Info("modules DoStart call ...")
//	for i := 0; i < len(modules); i++ {
//		module := modules[i].module
//		if module == nil {
//			continue
//		}
//
//		module.Go(rpc_do_start)
//	}
//}
//
//func check() {
//	for {
//		if int(startOkCount.Get()) >= total {
//			xsf_log.Info("================== all server module ok ==================")
//
//			xsf_server.Status.Set(xsf_server.ServerStatus_Running)
//
//			for i := 0; i < len(modules); i++ {
//				module := modules[i].module
//				if module != nil {
//					module.okSig <- true
//				}
//			}
//
//			return
//		}
//	}
//}
//
//func run(module *Module) {
//	module.run()
//	module.wg.Done()
//}
