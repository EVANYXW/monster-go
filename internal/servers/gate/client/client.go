package client

import (
	"github.com/evanyxw/monster-go/internal/servers/gate/manager"
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"sync/atomic"
)

type client struct {
	processor     *network.Processor
	netPoint      *network.NetPoint
	lastHeartbeat uint32
	rpcAcceptor   *rpc.Acceptor
	ID            atomic.Uint32

	CID         server.ClientID
	isHandshake atomic.Bool
}

func NewClient(np *network.NetPoint) *client {
	return &client{
		netPoint: np,
	}
}

func (c *client) Start() {

}

func (c *client) Init() {
	c.processor = network.NewProcessor()
	c.rpcAcceptor = rpc.NewAcceptor(100)

	c.processor.RegisterMsg(uint16(xsf_pb.MSGID_Clt_L_Login), c.OnNetMessage)
	c.processor.RegisterMsg(uint16(xsf_pb.MSGID_Clt_Gt_Handshake), c.Clt_Gt_Handshake)

	c.netPoint.SetNetEventRPC(c.rpcAcceptor)
	c.netPoint.SetProcessor(c.processor)
}

func (c *client) OnNetMessage(pack *network.Packet) {
	ep := rpc.GetClientDestEP(pack.Msg.ID)
	switch ep {
	case server.EP_Game: // 发往游戏服
		fallthrough
	case server.EP_Login: // 发往登录服
		fallthrough
	case server.EP_Mail: // 发往邮件服
		fallthrough
	case server.EP_Manager:
		connector := c.GetConnector(ep)
		if connector != nil {
			connector.SendMessage(&xsf_pb.Clt_L_Login{}, network.WithRaID(c.ID.Load()))

			// 写入ClientID
			//binary.LittleEndian.PutUint32(data[6:], c.ID.Get())
			//connector.SendData(data) // 直接把数据转发到内部服务器

			//xsf_log.Debug("ClientKernel OnNetData", xsf_log.Uint32("client id", c.ID))
		} else {
			//logger.Error("ClientKernel OnNetData connector not exist", xsf_log.Uint("ep", uint(epDest)), xsf_log.Uint("ep id", uint(c.server_ids[epDest])))
			//
			//if rawID > 0 { // 如果客户端有指定ID，但是连接不存在，则断开连接
			//	xsf_log.Error("ClientKernel OnNetData connector not exist, id error", xsf_log.Uint32("raw id", rawID))
			//	c.Disconnect(int32(xsf_pb.DisconnectReason_MsgInvalid), true)
			//}
		}
	default:
		logger.Error("ClientKernel OnNetData epDest error", zap.Int("epDest", ep))
		//c.Disconnect(int32(xsf_pb.DisconnectReason_MsgInvalid), true)
	}
}

func (c *client) GetConnector(ep int) *module.ConnectorKernel {
	managerModule := module.GetConnectorManager()
	connectorManager, ok := managerModule.(*manager.ConnectorManager)
	if !ok {
		logger.Error("GetConnector of module.GetConnectorManager is error!")
		return nil
	}
	iConnector := connectorManager.GetConnector(uint32(ep), 0)
	if iConnector == nil {
		logger.Error("GetConnector of module.IModuleKernel is error!")
		return nil
	}
	connector, ok := iConnector.(*module.ConnectorKernel)
	if !ok {
		logger.Error("GetConnector of interface to module.IModuleKernel is error!")
		return nil
	}

	//if ep == server.EP_Login {
	//	connector := connectorManager.GetConnector(uint32(ep))
	//	connector = manager.GetConnector(ep, 0)
	//	if connector != nil {
	//		c.server_ids[ep] = connector.ID
	//	}
	//} else {
	//	connector = manager.GetConnector(ep, c.server_ids[ep])
	//}

	return connector
}

func (c *client) OnNetConnected(np *network.NetPoint) {

}

func (c *client) OnRpcNetAccept(np *network.NetPoint) {

}

func (c *client) OnNetError(np *network.NetPoint) {

}

func (c *client) OnServerOk() {

}

func (c *client) OnOk() {

}

func (c *client) OnNPAdd(np *network.NetPoint) {

}

func (c *client) SendMessage(message proto.Message) {
	c.netPoint.SendMessage(message)
}

func (c *client) SetSignal(data []byte) {
	c.netPoint.SetSignal(data)
}

func (c *client) GetID() uint32 {
	return c.ID.Load()
}

func (c *client) MsgRegister(processor *network.Processor) {

}

func (c *client) OnHandshake() {
	c.isHandshake.Store(true)
}

func (c *client) Clt_Gt_Handshake(message *network.Packet) {
	c.OnHandshake()
	getMessage, _ := rpc.GetMessage(uint64(xsf_pb.MSGID_Gt_Clt_Handshake))
	msg := getMessage.(*xsf_pb.Gt_Clt_Handshake)
	c.SendMessage(msg)
}
