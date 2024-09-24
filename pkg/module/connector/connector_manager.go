package connector

import (
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/module/connector/handler"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/server"
	"math/rand"
)

type ep uint8
type id uint8

type Manager struct {
	kernel      module.IModuleKernel
	collections []map[uint32]*module.ConnectorKernel
	handler     module.GateAcceptorHandler
	id          int32
}

func NewManager(id int32) *Manager {
	c := &Manager{
		id: id,
	}
	hdler := handler.NewManagerMsg()
	c.handler = hdler
	c.kernel = module.NewKernel(hdler, network.NetPointManager.GetRpcAcceptor(),
		network.NetPointManager.GetProcessor())

	return c
}

func (c *Manager) Init(baseModule *module.BaseModule) bool {
	//c.collections = make([]connectCollection, xsf_util.EP_Max)
	c.collections = []map[uint32]*module.ConnectorKernel{}
	for i := 0; i < server.EP_Max; i++ {
		c.collections = append(c.collections, make(map[uint32]*module.ConnectorKernel))
	}
	//c.kernel.Init()
	return true
}

func (c Manager) DoRegister() {
	//c.kernel.DoRegist()
}

func (c Manager) DoRun() {
	//c.kernel.Start()
}

func (c *Manager) DoWaitStart() {

}

func (c Manager) DoRelease() {
	for _, ckArr := range c.collections {
		for _, ck := range ckArr {
			ck.DoRelease()
		}
	}
}

func (c *Manager) OnOk() {

}

func (c Manager) OnStartCheck() int {
	return module.ModuleRunCode_Ok
}

func (c Manager) OnCloseCheck() int {
	return module.ModuleRunCode_Ok
}

func (c Manager) Update() {

}

func (c *Manager) GetID() int32 {
	return c.id
}

func (c Manager) GetKernel() module.IModuleKernel {
	return nil
}

func (c Manager) GetConnector(ep uint32, id uint32) module.IModuleKernel {
	eps := c.collections[ep]
	if id == 0 {
		len := len(eps)
		if len <= 0 {
			return nil
		}
		index := rand.Intn(len)
		var num uint32 = 0
		for _, item := range eps {
			if index == int(num) {
				return item
			}
			num++
		}
	} else {
		return eps[id]
	}
	return nil
}

func (c *Manager) CreateConnector(id uint32, ip string, port uint32) *module.ConnectorKernel {
	//msgHandler := handler.NewManager()
	msgHandler := c.handler.(module.MsgHandler)
	ck := module.NewConnectorKernel(ip, port, msgHandler, new(network.ClientPackerFactory))
	//ck := module.NewConnectorKernel(ip, port, msgHandler, new(network.DefaultPackerFactory))
	ck.SetID(id)

	c.collections[ck.SID.Type][id] = ck

	ck.DoRegister()
	ck.DoRun()
	//// fixMe 这里会不会没有运行好，在发送Handshake
	//time.Sleep(time.Duration(3) * time.Second)
	//c.handler.SendHandshake(ck)

	return ck
}
