package manager

import (
	"github.com/evanyxw/monster-go/internal/servers"
	"github.com/evanyxw/monster-go/internal/servers/gate/manager/handler"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/server"
	"math/rand"
)

type ep uint8
type id uint8

type ConnectorManager struct {
	*module.BaseModule
	kernel      *module.Kernel
	collections []map[uint32]*module.ConnectorKernel
	ID          int32
}

func NewConnectorManager(id int32) *ConnectorManager {
	c := &ConnectorManager{
		ID: id,
		//collections: ,
		//kernel:     module.NewKernel(handler.NewManager(false)),
	}
	msgHandler := handler.NewManager()
	c.kernel = module.NewKernel(msgHandler, servers.NetPointManager.GetRpcAcceptor(), servers.NetPointManager.GetProcessor())
	module.NewBaseModule(c)

	return c
}

func (c ConnectorManager) GetID() int32 {
	return c.ID
}

func (c *ConnectorManager) Init() {
	//c.collections = make([]connectCollection, xsf_util.EP_Max)
	c.collections = []map[uint32]*module.ConnectorKernel{}
	for i := 0; i < server.EP_Max; i++ {
		c.collections = append(c.collections, make(map[uint32]*module.ConnectorKernel))
	}
	//c.kernel.Init()
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
			ck.Release()
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
	ck.SetID(id)

	c.collections[ck.SID.Type][id] = ck

	ck.DoRegist()
	async.Go(func() {
		ck.Start()
	})
	// fixMe 这里会不会没有运行好，在发送Handshake
	msgHandler.SendHandshake(ck)
	return ck
}

//
//Start()
//OnNetMessage(pack *network.Packet)
//MsgRegister(processor *network.Processor)
