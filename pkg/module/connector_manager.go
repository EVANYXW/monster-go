package module

import "github.com/evanyxw/monster-go/pkg/logger"

type ConnectorManager struct {
	*BaseModule
	connectors map[uint32]*ConnectorKernel
	ID         int32
}

func NewConnectorManager(id int32) *ConnectorManager {
	c := &ConnectorManager{
		ID:         id,
		connectors: make(map[uint32]*ConnectorKernel),
	}

	NewBaseModule(c)

	return c
}

func (c ConnectorManager) GetID() int32 {
	return c.ID
}

func (c *ConnectorManager) Init() {
	//c.collections = make([]connectCollection, xsf_util.EP_Max)
	//for i := 0; i < xsf_util.EP_Max; i++ {
	//	c.collections[i].connectors = make(map[uint32]*SingleConnector)
	//}
}

func (c ConnectorManager) DoRun() {

}

func (c ConnectorManager) DoRelease() {
	//for _, conn := range c.connectors {
	//	conn.OnStartClose()
	//}
}

func (c ConnectorManager) OnStartCheck() int {
	return ModuleRunCode_Ok
}

func (c ConnectorManager) OnCloseCheck() int {
	return ModuleRunCode_Ok
}

func (c ConnectorManager) Update() {

}

func (c ConnectorManager) GetKernel() IModuleKernel {
	return nil
}

func (c *ConnectorManager) CreateConnector(id uint32, ip string, port uint32) *ConnectorKernel {
	ck := NewConnectorKernel(ip, port, nil)
	ck.SetID(id)
	c.connectors[id] = ck
	logger.Info("CreateConnector success")
	ck.Start()
	logger.Info("CreateConnector over")
	return ck
}
