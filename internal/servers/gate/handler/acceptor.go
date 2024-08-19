// Package handler @Author evan_yxw
// @Date 2024/8/16 15:30:00
// @Desc
package handler

//type IAcceptor interface {
//	//AddCloseEventRpc(rpc IClientCloseHandler)
//	//RegistClientMessage(msgID uint16, rpc *xsf_rpc.Acceptor)
//
//	// 断开一个客户端连接
//	//DisconnectClient(id uint32, reason uint32)
//
//	// 发送消息到网关
//	//SendMessage2Agent(agentID uint32, msg xsf_net.IMessage)
//
//	// 发送消息到客户端
//	SendMessage2Client(clientID uint32, message proto.Message, nodePoint *network.NetPoint)
//
//	// 广播一个消息给所有客户端
//	//Broadcast(msg xsf_net.IMessage)
//
//	// 设置一个客户端的服务器转播id
//	//SetServerID(clientID uint32, ep uint8, serverID uint32)
//}
//
//type acceptor struct {
//}
//
//var (
//	acc *acceptor
//)
//
//func GetAcceptor() IAcceptor {
//	return &acceptor{}
//}
//
//func (a *acceptor) SendMessage2Client(clientID uint32, message proto.Message, nodePoint *network.NetPoint) {
//	msg, _ := rpc.GetMessage(uint64(xsf_pb.SMSGID_GtA_Gt_ClientMessage))
//	localMsg := msg.(*xsf_pb.GtA_Gt_ClientMessage)
//	localMsg.ClientId = append(localMsg.ClientId, clientID)
//	//localMsg.ClientMessage = network.ClientBufferPacker{}.Pack()
//	nodePoint.SendMessage(localMsg)
//}
