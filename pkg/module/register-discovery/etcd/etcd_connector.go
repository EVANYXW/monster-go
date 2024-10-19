package etcd

import (
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/pkg/grpcpool"
	"github.com/evanyxw/monster-go/pkg/kernel"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module/module_def"
)

type EtcdConnector struct {
	kernel     module_def.IKernel
	id         int32
	servername string
}

func NewEtcdConnector(id int32, servername string, isWatch bool, netType kernel.NetType) *EtcdConnector {
	etcdCnf := configs.All().Etcd
	etcdClient := grpcpool.InitEtcdClient(etcdCnf.Addr, etcdCnf.User, etcdCnf.Pass)
	c := &EtcdConnector{
		id: id,
		kernel: kernel.NewEtcdKernel(servername, isWatch, netType, etcdClient, logger.GetLogger(),
			kernel.WithCNoWaitStart(true)),
	}

	//module.ConnKernel = c.kernel.(*module.ConnectorKernel)

	return c
}

func (c *EtcdConnector) Init(baseModule module_def.IBaseModule) bool {
	c.kernel.Init(baseModule)
	return true
}

func (c *EtcdConnector) DoRun() {
	c.kernel.DoRun()
}

func (c *EtcdConnector) DoWaitStart() {

}

func (c *EtcdConnector) DoRelease() {
	c.kernel.DoRelease()
}

func (c *EtcdConnector) OnOk() {
	c.kernel.OnOk()
}

func (c *EtcdConnector) OnStartCheck() int {
	return c.kernel.OnStartCheck()
}

func (c *EtcdConnector) OnCloseCheck() int {
	return c.kernel.OnCloseCheck()
}

func (c *EtcdConnector) GetID() int32 {
	return c.id
}

func (c *EtcdConnector) GetKernel() module_def.IKernel {
	return c.kernel
}

func (c *EtcdConnector) Update() {

}

func (c *EtcdConnector) DoRegister() {
	c.kernel.DoRegister()
}
