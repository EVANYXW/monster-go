package module

import (
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/server"
	"math/rand"
)

type ep uint8
type id uint8

type Manager struct {
	kernel      IKernel
	collections []map[uint32]IKernel
	handler     GateAcceptorHandler
	id          int32
	//factory     ManagerFactory
}

func NewManager(id int32, msgHandler MsgHandler) *Manager {
	c := &Manager{
		id: id,
	}
	hdler := msgHandler
	gateAcceptorHandler := msgHandler.(GateAcceptorHandler)
	c.handler = gateAcceptorHandler
	c.kernel = NewKernel(hdler, network.NetPointManager.GetRpcAcceptor(),
		network.NetPointManager.GetProcessor())

	return c
	//return factory.Create(id)
}

func (c *Manager) Init(baseModule *BaseModule) bool {
	//c.collections = make([]connectCollection, xsf_util.EP_Max)
	c.collections = []map[uint32]IKernel{}
	for i := 0; i < server.EP_Max; i++ {
		c.collections = append(c.collections, make(map[uint32]IKernel))
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
	return ModuleRunCode_Ok
}

func (c Manager) OnCloseCheck() int {
	return ModuleRunCode_Ok
}

func (c Manager) Update() {

}

func (c *Manager) GetID() int32 {
	return c.id
}

func (c Manager) GetKernel() IKernel {
	return nil
}

func (c Manager) GetConnector(ep uint32, id uint32) IKernel {
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

func (c Manager) DelConnector(id uint32) {
	var sid server.ServerID
	server.ID2Sid(id, &sid)
	if len(c.collections) > int(sid.Type) {
		if _, ok := c.collections[sid.Type][id]; ok {
			delete(c.collections[sid.Type], id)
		}
	}
}

func (c *Manager) CreateConnector(id uint32, ip string, port uint32) *ConnectorKernel {
	//msgHandler := handler.NewManager()
	msgHandler := c.handler.(MsgHandler)
	ck := NewConnectorKernel(ip, port, msgHandler, new(network.ClientPackerFactory))
	//ck := module.NewConnectorKernel(ip, port, msgHandler, new(network.DefaultPackerFactory))
	ck.SetID(id)
	c.collections[ck.SID.Type][id] = ck
	ck.DoRegister()
	ck.DoRun()
	return ck

	// factory
	//ck := c.factory.CreateConnector(c.handler, id, ip, port)
	//c.collections[ck.SID.Type][id] = ck
	//
	//ck.DoRegister()
	//ck.DoRun()
	//
	//return ck
}
