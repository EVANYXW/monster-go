package client

import (
	"fmt"
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/module/connector"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"sync/atomic"
	"time"
)

const (
	rpc_SET_SERVER_ID = "18099"
	rpc_DISCONNECT    = "DISCONNECT"
)

type Client struct {
	processor     *network.Processor
	netPoint      *network.NetPoint
	lastHeartbeat uint64
	rpcAcceptor   *rpc.Acceptor
	ID            atomic.Uint32

	CID         server.ClientID
	isHandshake atomic.Bool
	server_ids  []uint32
}

func NewClient(np *network.NetPoint) *Client {
	return &Client{
		netPoint:   np,
		server_ids: make([]uint32, server.EP_Max),
	}
}

func (c *Client) OnInit(baseModule *module.BaseModule) {

}

func (c *Client) Start() {

}

func (c *Client) Close() {
	c.netPoint.Close()
	c.ID.Store(0)
	c.lastHeartbeat = 0
}

func (c *Client) GoDisconnect(id uint32) {
	fmt.Println("login 让我踢掉你～")
	c.Close()
}

func (c *Client) Init() {
	c.processor = network.NewProcessor()

	c.processor.RegisterMsg(uint16(xsf_pb.MSGID_Clt_L_Login), c.OnNetMessage)

	c.processor.RegisterMsg(uint16(xsf_pb.MSGID_Clt_Gt_Handshake), c.Clt_Gt_Handshake)
	c.processor.RegisterMsg(uint16(xsf_pb.MSGID_Clt_Gt_Heartbeat), c.Clt_Gt_Heartbeat)
	//c.processor.RegisterMsg(uint16(18099), c.setServerID)

	//c.rpcAcceptor = rpc.NewAcceptor(100)
	//c.netPoint.SetNetEventRPC(c.rpcAcceptor)
	c.netPoint.SetProcessor(c.processor)

	c.server_ids[server.EP_Mail] = module.MailID.Load()
	c.server_ids[server.EP_Manager] = module.ManagerID.Load()

	//c.rpcAcceptor.Run()
}

func (c *Client) OnNetMessage(pack *network.Packet) {
	ep := rpc.GetClientDestEP(pack.Msg.ID)

	if pack.Msg.ID == uint64(xsf_pb.MSGID_Clt_G_Relogin) {
		// 短线重连rawId是game服务器的id
		//xsf_log.Debug("client relogin, set server id", xsf_log.Uint32("server", rawID))
		c.server_ids[ep] = pack.Msg.RawID
	}

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
			message, _ := rpc.GetMessage(pack.Msg.ID)
			rpc.Import(pack.Msg.Data, message)
			// fixMe login 没开报错
			connector.SendMessage(message, network.WithRaID(c.ID.Load()))

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

func (c *Client) GetConnector(ep int) *module.ConnectorKernel {
	managerModule := module.GetConnectorManager()
	connectorManager, ok := managerModule.(*connector.Manager)
	if !ok {
		logger.Error("GetConnector of module.GetConnectorManager is error!")
		return nil
	}

	var iConnector module.IModuleKernel
	if ep == server.EP_Login {
		iConnector = connectorManager.GetConnector(uint32(ep), 0)
	} else {
		iConnector = connectorManager.GetConnector(uint32(ep), c.server_ids[ep])
	}

	if iConnector == nil {
		logger.Error("GetConnector of module.IModuleKernel is error!")
		return nil
	}

	connector, ok := iConnector.(*module.ConnectorKernel)
	if !ok {
		logger.Error("GetConnector of interface to module.IModuleKernel is error!")
		return nil
	}

	if ep == server.EP_Login {
		if connector.ID != c.server_ids[ep] {
			c.server_ids[ep] = connector.ID
		}
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

func (c *Client) SetServerID(args []interface{}) {
	ep := args[0].(uint32)
	serverID := args[1].(uint32)

	c.server_ids[ep] = serverID
}

func (c *Client) GetServerIds() []uint32 {
	return c.server_ids
}

func (c *Client) GetExistConnector(ep uint32) *module.ConnectorKernel {
	managerModule := module.GetConnectorManager()
	connectorManager, ok := managerModule.(*connector.Manager)
	if !ok {
		logger.Error("GetConnector of module.GetConnectorManager is error!")
		return nil
	}

	if c.server_ids[ep] == 0 {
		return nil
	} else {
		iGetConnector := connectorManager.GetConnector(ep, c.server_ids[ep])
		connectorKernel := iGetConnector.(*module.ConnectorKernel)
		return connectorKernel
	}
}

func (c *Client) OnNetConnected(np *network.NetPoint) {

}

func (c *Client) OnRpcNetAccept(np *network.NetPoint, acceptor *network.Acceptor) {

}

func (c *Client) OnNetError(np *network.NetPoint, acceptor *network.Acceptor) {

}

func (c *Client) OnServerOk() {

}

func (c *Client) OnOk() {

}

func (c *Client) OnUpdate() {

}

func (c *Client) OnNPAdd(np *network.NetPoint) {

}

func (c *Client) SendMessage(message proto.Message) {
	c.netPoint.SendMessage(message)
}

func (c *Client) SetSignal(data []byte) {
	c.netPoint.SetSignal(data)
}

func (c *Client) GetID() uint32 {
	return c.ID.Load()
}

func (c *Client) GetLastHeartbeat() *uint64 {
	return &c.lastHeartbeat
}

func (c *Client) MsgRegister(processor *network.Processor) {

}

func (c *Client) OnHandshake() {
	c.isHandshake.Store(true)
}

func (c *Client) Clt_Gt_Handshake(message *network.Packet) {
	c.OnHandshake()
	getMessage, _ := rpc.GetMessage(uint64(xsf_pb.MSGID_Gt_Clt_Handshake))
	msg := getMessage.(*xsf_pb.Gt_Clt_Handshake)
	c.SendMessage(msg)
}

func (c *Client) Clt_Gt_Heartbeat(message *network.Packet) {
	c.lastHeartbeat = uint64(time.Now().Unix())
}
