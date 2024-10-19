package module

import (
	"github.com/evanyxw/monster-go/pkg/kernel"
	"github.com/evanyxw/monster-go/pkg/module/module_def"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/server"
)

type GrpcClient struct {
	kernel module_def.IKernel
	id     int32
}

func NewGrpcClient(id int32) *GrpcClient {
	c := &GrpcClient{
		id:     id,
		kernel: nil,
	}

	return c
}

func (c *GrpcClient) Init(baseModule *module_def.BaseModule) bool {
	c.kernel.Init(baseModule)
	return true
}

// DoRun BaseModule 调用
func (c *GrpcClient) DoRun() {
	c.kernel.DoRun()
}

func (c *GrpcClient) DoWaitStart() {
	c.kernel.DoWaitStart()
}

func (c *GrpcClient) DoRelease() {
	c.kernel.DoRelease()
}

func (c *GrpcClient) OnOk() {

}

func (c *GrpcClient) OnStartCheck() int {
	// TCP链接准备好
	if c.kernel.GetStatus() == server.Net_RunStep_Done {
		return module_def.ModuleOk()
	}
	return module_def.ModuleWait()
}

func (c *GrpcClient) OnCloseCheck() int {
	return c.kernel.OnCloseCheck()
}

func (c *GrpcClient) GetID() int32 {
	return c.id
}

func (c *GrpcClient) GetKernel() module_def.IKernel {
	return c.kernel
}

func (c *GrpcClient) Update() {

}

func (c *GrpcClient) DoRegister() {
	c.kernel.DoRegister()
}

func (c *GrpcClient) OnNetError(np *network.NetPoint) {
	kernel.NodeManager.OnNodeLost(np.ID, np.SID.Type)
}

func (c *GrpcClient) OnServerOk() {

}

func (c *GrpcClient) OnNPAdd(np *network.NetPoint) {

}
