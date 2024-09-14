package network

import (
	"container/list"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"net"
	"time"
)

// 数组存储服务器节点管理
type NPManagerNormal struct {
	maxCount      uint32 // 容纳最大个数
	nps           map[uint32]*NetPoint
	tempConns     *list.List // 临时连接集合
	rpcAcceptor   *rpc.Acceptor
	processor     *Processor
	packerFactory PackerFactory
}

func NewNormal(maxCount uint32, rpcAcceptor *rpc.Acceptor, processor *Processor, packerFactory PackerFactory) *NPManagerNormal {
	return &NPManagerNormal{
		maxCount:      maxCount,
		nps:           make(map[uint32]*NetPoint),
		tempConns:     list.New(),
		rpcAcceptor:   rpcAcceptor,
		processor:     processor,
		packerFactory: packerFactory,
	}
}

func (mn *NPManagerNormal) New(conn *net.TCPConn) *NetPoint {
	point, _ := NewNetPoint(conn, mn.packerFactory)
	logger.Info("New net point:", zap.String("ip:", point.RemoteIP))
	mn.tempConns.PushBack(point)

	point.SetNetEventRPC(mn.rpcAcceptor)
	point.SetProcessor(mn.processor)
	return point
}

func (mn *NPManagerNormal) GetProcessor() *Processor {
	return mn.processor
}

func (mn *NPManagerNormal) GetRpcAcceptor() *rpc.Acceptor {
	return mn.rpcAcceptor
}

func (mn *NPManagerNormal) Get(index uint32) *NetPoint {
	// fixMe 用id做下标就有问题
	//if index >= mn.maxCount {
	//	return nil
	//}

	return mn.nps[index]
}

func (mn *NPManagerNormal) GetMaxConnNum() uint32 {
	return mn.maxCount
}

func (mn *NPManagerNormal) Add(np *NetPoint) bool {
	// 只有服务之间的握手加入了
	if np.SID.ID == 0 {
		//xsf_log.Error("NPManager_Array Add index error", xsf_log.Int("index", int(np.SID.Index)))
		return false
	}

	if len(mn.nps) >= int(mn.maxCount) {
		logger.Error("NPManager_Array add reach max count",
			zap.Int("mn.maxCount", int(mn.maxCount)), zap.Int("index", int(np.SID.Index)))
		return false
	}

	if mn.Get(uint32(np.SID.ID)) != nil {
		logger.Info("npManager_Normal add NetPoint exist", zap.Int("server", int(np.SID.ID)),
			zap.Int("type", int(np.SID.Type)),
			zap.Int("index", int(np.SID.Index)))
		return false
	}

	mn.nps[np.ID] = np
	//mn.nps[uint32(np.SID.ID)] = np

	return true
}

func (mn *NPManagerNormal) Del(np *NetPoint) {
	npCheck := delFromList(mn.tempConns, np)
	if npCheck != nil {
		//xsf_log.Debug("NPManager_Array delete from tempConns")
	} else {
		delete(mn.nps, np.ID)
	}

	//xsf_log.Info(">>>> server logout", xsf_log.Int("remote id", int(np.SID.ID)),
	//	xsf_log.Int("remote type", int(np.SID.Type)),
	//	xsf_log.Int("remote index", int(np.SID.Index)))
}

func (mn *NPManagerNormal) OnUpdate() {
	nCurTime := time.Now().Unix()

	var next *list.Element
	for e := mn.tempConns.Front(); e != nil; e = next {
		next = e.Next()
		np := e.Value.(*NetPoint)
		if uint32(nCurTime) > np.Time+1000 { // 关闭超时的连接
			//xsf_log.Info("NPManager_Normal OnUpdate timeout", xsf_log.Uint("nCurTime", uint(nCurTime)), xsf_log.Uint("np.time", uint(np.time)))
			np.Close()
			mn.tempConns.Remove(e)
		}
	}
}

func (mn *NPManagerNormal) OnHandshake(np *NetPoint) bool {
	//xsf_log.Debug("NPManager_Array OnHandshake")
	var npCheck *NetPoint
	for e := mn.tempConns.Front(); e != nil; e = e.Next() {
		npCheck = e.Value.(*NetPoint)
		if npCheck == np {
			mn.tempConns.Remove(e)
			break
		}
	}

	if npCheck == nil {
		logger.Info("NPManager_Array OnHandshake np not found", zap.Int("remote id", int(np.SID.ID)),
			zap.Int("remote type", int(np.SID.Type)), zap.Int("remote index", int(np.SID.Index)))

		np.Close()
		return false
	}

	ok := mn.Add(np)
	if ok {
		np.OnHandshakeTicker(np)
		logger.Info("<<<< server login", zap.Int("remote id", int(np.SID.ID)),
			zap.Int("remote type", int(np.SID.Type)), zap.Int("remote index", int(np.SID.Index)))
		//todo
		//np.module.Kernel.(*NetKernel).Handler.OnNPAdd(np)

	} else {
		logger.Info("NPManager_Array OnHandshake mn.Add error", zap.Int("remote id", int(np.SID.ID)),
			zap.Int("remote type", int(np.SID.Type)), zap.Int("remote index", int(np.SID.Index)))

		np.Close()
	}

	return ok
}

func (mn *NPManagerNormal) Broadcast(msg proto.Message, skip uint32) {

	for _, v := range mn.nps {
		if v != nil && v.ID != skip {
			packData, _ := v.Pack(msg)
			v.SetSignal(packData)
		}
	}

	//xsf_net.ReturnMessage(msg)
}
