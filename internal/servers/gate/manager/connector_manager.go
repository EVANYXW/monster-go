package manager

import (
	"github.com/evanyxw/monster-go/internal/servers/gate/manager/handler"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/server"
	"math/rand"
	"time"
)

type ep uint8
type id uint8

type ConnectorManager struct {
	kernel      module.IModuleKernel
	collections []map[uint32]*module.ConnectorKernel
	ID          int32
}

func NewConnectorManager(id int32) *ConnectorManager {
	c := &ConnectorManager{
		ID: id,
	}
	msgHandler := handler.NewManager()
	c.kernel = module.NewKernel(msgHandler, network.NetPointManager.GetRpcAcceptor(), network.NetPointManager.GetProcessor())
	module.NewBaseModule(c)

	return c
}

func (c ConnectorManager) GetID() int32 {
	return c.ID
}

func (c *ConnectorManager) Init() bool {
	//c.collections = make([]connectCollection, xsf_util.EP_Max)
	c.collections = []map[uint32]*module.ConnectorKernel{}
	for i := 0; i < server.EP_Max; i++ {
		c.collections = append(c.collections, make(map[uint32]*module.ConnectorKernel))
	}
	//c.kernel.Init()
	return true
}

func (c ConnectorManager) DoRegister() {
	//c.kernel.DoRegist()
}

func (c ConnectorManager) DoRun() {
	//c.kernel.Start()
}

func (c *ConnectorManager) DoWaitStart() {

}

func (c ConnectorManager) DoRelease() {
	for _, ckArr := range c.collections {
		for _, ck := range ckArr {
			ck.DoRelease()
		}
	}
}

func (c *ConnectorManager) OnOk() {

}

func (c ConnectorManager) OnStartCheck() int {
	return module.ModuleRunCode_Ok
}

func (c ConnectorManager) OnCloseCheck() int {
	return module.ModuleRunCode_Ok
}

func (c ConnectorManager) Update() {

}

func (c ConnectorManager) GetKernel() module.IModuleKernel {
	return nil
}

func (c ConnectorManager) GetConnector(ep uint32, id uint32) module.IModuleKernel {
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

func (c *ConnectorManager) CreateConnector(id uint32, ip string, port uint32) *module.ConnectorKernel {
	msgHandler := handler.NewManager()
	ck := module.NewConnectorKernel(ip, port, msgHandler, new(network.ClientPackerFactory))
	//ck := module.NewConnectorKernel(ip, port, msgHandler, new(network.DefaultPackerFactory))
	ck.SetID(id)

	c.collections[ck.SID.Type][id] = ck

	ck.DoRegister()
	//async.Go(func() {
	ck.DoRun()
	//})
	// fixMe 这里会不会没有运行好，在发送Handshake
	time.Sleep(time.Duration(3) * time.Second)
	msgHandler.SendHandshake(ck)

	return ck
}
