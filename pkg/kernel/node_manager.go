package kernel

import (
	"container/list"
	"fmt"
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module/module_def"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"
	"go.uber.org/zap"
)

var (
	NodeManager      module_def.INodeManager
	ConnKernel       *ConnectorKernel
	ClientManager    module_def.IGtClientManager
	GtAClientManager module_def.IGtAClientManager
)

//type INodeManager interface {
//	Start()
//	GetIndex(sid *server.ServerID)
//	AddNode(id uint32, ip string, ports []uint32) *network.ServerInfo
//	Send(np *network.NetPoint, si *network.ServerInfo)
//	OnNodeLost(id uint32, ep uint8)
//	OnNodeOK(id uint32)
//	Broadcast(si *network.ServerInfo)
//}

type nodeManager struct {
	innerPort  uint32
	outPort    uint32
	nodesIndex []int
	lostList   []list.List
	nodes      map[uint32]*network.ServerInfo
}

func NewNodeManager() module_def.INodeManager {
	return &nodeManager{
		lostList:   make([]list.List, server.EP_Max),
		nodes:      make(map[uint32]*network.ServerInfo),
		nodesIndex: make([]int, server.EP_Max),
	}
}

func (nm *nodeManager) Start() {
	for i := 0; i < server.EP_Max; i++ {
		nm.nodesIndex[i] = 1
	}
	serverCnf := configs.All()
	nm.innerPort = uint32(serverCnf.InnerPort)
	nm.outPort = uint32(serverCnf.OutPort)
	nm.nodes = make(map[uint32]*network.ServerInfo)
}

func (nm *nodeManager) GetIndex(sid *server.ServerID) {
	sid.Index = uint8(nm.nodesIndex[sid.Type])
	nm.nodesIndex[sid.Type] = int(sid.Index) + 1
}

func (nm *nodeManager) GetOldSI(ep uint8, id uint32, CheckEqual bool) (si *network.ServerInfo) {
	list := nm.lostList[ep]
	for s := list.Front(); s != nil; s = s.Next() {
		serverInfo := s.Value.(*network.ServerInfo)
		if CheckEqual {
			if serverInfo.ID == id {
				si = serverInfo
				list.Remove(s)
				return
			}
		} else {
			si = serverInfo
			list.Remove(s)
			return
		}
	}

	return nil
}

func (nm *nodeManager) getNextPort(IsInner bool) uint32 {
	var port uint32
	if IsInner {
		port = nm.innerPort
		nm.innerPort++
	} else {
		port = nm.outPort
		nm.outPort++
	}

	return port
}

func (nm *nodeManager) GetPorts(ep uint8) []uint32 {
	ports := make([]uint32, server.EP_Max)
	switch ep {
	case server.EP_Gate:
		ports[server.EP_Client] = nm.getNextPort(false)
		ports[server.EP_Robot] = nm.getNextPort(true)
	case server.EP_Game:
		ports[server.EP_Gate] = nm.getNextPort(true)
		ports[server.EP_Robot] = nm.getNextPort(true)
	case server.EP_Login:
		ports[server.EP_Gate] = nm.getNextPort(true)
		ports[server.EP_Robot] = nm.getNextPort(true)
	case server.EP_Manager:
	case server.EP_World:
		ports[server.EP_Gate] = nm.getNextPort(true)
		ports[server.EP_Robot] = nm.getNextPort(true)
	case server.EP_Mail:
		ports[server.EP_Gate] = nm.getNextPort(true)
	case server.EP_Center:
		ports[server.EP_Robot] = nm.getNextPort(true)
	}
	fmt.Println("ports:", ports)

	return ports
}

func (nm *nodeManager) AddNode(id uint32, ip string, ports []uint32) *network.ServerInfo {
	var SID server.ServerID
	server.ID2Sid(id, &SID)

	var si *network.ServerInfo
	if SID.Index == 0 { // 如果是新来的服务器
		oldSi := nm.GetOldSI(SID.Type, id, false)
		if oldSi == nil {
			nm.GetIndex(&SID)
			newPorts := nm.GetPorts(SID.Type)

			si = new(network.ServerInfo)
			si.ID = server.Sid2ID(&SID)
			si.Ports = [server.EP_Max]uint32(newPorts)
			si.IP = ip
			si.Status = network.ServerInfo_New

			//xsf_log.Info("nodeManager AddNode new", xsf_log.Uint32("id", si.ID), xsf_log.String("ip", ip))
		} else {
			si = oldSi
			si.Status = network.ServerInfo_New
			si.IP = ip

			//xsf_log.Info("nodeManager AddNode new, find old server info", xsf_log.Uint32("id", si.ID), xsf_log.String("ip", ip))
		}

		nm.nodes[si.ID] = si

	} else { // 如果是旧服务器
		osi, ok := nm.nodes[id]
		if ok {
			si = osi // 已经存在了，直接返回之前的
			//xsf_log.Info("nodeManager AddNode exsit", xsf_log.Uint32("id", si.ID), xsf_log.String("ip", ip))
		} else {
			oldSi := nm.GetOldSI(SID.Type, id, true)
			if oldSi == nil {
				//xsf_log.Error("nodeManager AddNode not found old server info", xsf_log.Uint32("id", id), xsf_log.String("ip", ip))
				return nil
			} else {
				if oldSi.IP != ip {
					nm.lostList[SID.Type].PushBack(oldSi)
					//xsf_log.Error("nodeManager AddNode, ip error", xsf_log.String("old ip", oldSi.IP), xsf_log.String("new ip", ip))
					return nil
				}

				for i := 0; i < server.EP_Max; i++ {
					if oldSi.Ports[i] != ports[i] {
						//xsf_log.Error("nodeManager AddNode, port error", xsf_log.Int("ep", i), xsf_log.Uint32("old port", oldSi.Ports[i]), xsf_log.Uint32("new port", ports[i]))
						return nil
					}
				}
			}

			si = oldSi
			nm.nodes[si.ID] = si
			//xsf_log.Info("nodeManager AddNode find old", xsf_log.Uint32("id", si.ID), xsf_log.String("ip", ip))
		}
	}

	return si
}

func (nm *nodeManager) OnNodeLost(id uint32, ep uint8) {
	node, ok := nm.nodes[id]
	if ok {
		nm.lostList[ep].PushBack(node)
		delete(nm.nodes, id)
	} else {
		//xsf_log.Info("nodeManager OnNodeLost node not exist", xsf_log.Uint32("id", id))
		return
	}

	//downMsg := xsf_net.GetMessage(uint16(xsf_pb.SMSGID_C_Cc_ServerLost))
	//localMsg := downMsg.(*xsf_message.MSG_C_Cc_ServerLost)
	//localMsg.PB.ServerId = id

	//downMsg := xsf_pb.C_Cc_ServerLost{}
	//downMsg.ServerId = id

	message, _ := rpc.GetMessage(uint64(xsf_pb.SMSGID_C_Cc_ServerLost))
	downMsg := message.(*xsf_pb.C_Cc_ServerLost)
	downMsg.ServerId = id

	logger.Info("服务器已离线", zap.Uint32("id", id))
	manager := module_def.GetManager(module_def.ModuleID_SM)
	manager.Broadcast(downMsg, 0)
}

// 把当前服务器所有节点数据发送给np
func (nm *nodeManager) Send(np *network.NetPoint, si *network.ServerInfo) {

	localMsg := &xsf_pb.C_Cc_ServerInfo{}

	for _, s := range nm.nodes {
		if s != si {
			msi := &xsf_pb.MSG_ServerInfo{}
			msi.ServerId = s.ID
			msi.Ip = s.IP
			msi.Ports = s.Ports[:]
			msi.Status = s.Status
			localMsg.Infos = append(localMsg.Infos, msi)
		}
	}

	//xsf_log.Info("nodeManager Send", xsf_log.String("message", downMsg.ToString()))

	if len(localMsg.Infos) > 0 {
		np.SendMessage(localMsg)
	}
}

func (nm *nodeManager) Broadcast(si *network.ServerInfo) {
	msg := &xsf_pb.C_Cc_ServerInfo{}
	siMsg := &xsf_pb.MSG_ServerInfo{}
	siMsg.ServerId = si.ID
	siMsg.Ip = si.IP
	siMsg.Ports = si.Ports[:]
	siMsg.Status = si.Status
	msg.Infos = append(msg.Infos, siMsg)

	manager := module_def.GetManager(module_def.ModuleID_SM)
	manager.Broadcast(msg, si.ID)
}

func (nm *nodeManager) OnNodeOK(id uint32) {
	node, ok := nm.nodes[id]
	if ok {
		node.Status = network.ServerInfo_Ok
	} else {
		logger.Error("nodeManager OnNodeOK node not exist", zap.Uint32("id", id))
		return
	}

	downMsg := &xsf_pb.C_Cc_ServerOk{}
	downMsg.ServerId = id

	//xsf_log.Info("【中心服】收到服务器已准备好", xsf_log.Uint32("id", id))
	// todo
	manager := module_def.GetManager(module_def.ModuleID_SM)
	manager.Broadcast(downMsg, id)

	isAllOK := true
	for _, node := range nm.nodes {
		if node.Status != network.ServerInfo_Ok {
			isAllOK = false
		}
	}
	// todo len(NodeList) 服务器配置数据
	if isAllOK && len(nm.nodes) >= 6 {
		version := ""
		// todo
		//path := xsf_server.RunDirPrefix + "/version"
		//vData, err := os.ReadFile(path)
		//if err != nil {
		//	version = "版本数据读取错误：" + err.Error()
		//	logger.Error("read version file error, " + err.Error())
		//} else {
		//	version = string(vData)
		//}

		message := fmt.Sprintf("%v-%v 所有服务器节点都启动完毕\n版本：%v", "center", "01", version)
		logger.Info(message)
		//notice.SendNotice(message)

		server.Ports = [10]uint32(nm.GetPorts(server.EP_Center))
		// todo
		module_def.DoWaitStart()
	}
	logger.Info("OnNodeOK node length:", zap.Int("node length:", len(nm.nodes)))
}
