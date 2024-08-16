package network

import (
	"container/list"
	"fmt"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"net"
	"time"
)

// 数组存储服务器节点管理
type NPManagerArray struct {
	maxCount    uint32 // 容纳最大个数
	nps         []*NetPoint
	tempConns   *list.List // 临时连接集合
	rpcAcceptor *rpc.Acceptor
	packer      Packer
}

func NewArray(maxCount uint32, rpcAcceptor *rpc.Acceptor) *NPManagerArray {
	return &NPManagerArray{
		maxCount: maxCount,
		//nps:         make([]*NetPoint, 1),
		tempConns:   list.New(),
		rpcAcceptor: rpcAcceptor,
		packer:      NewDefaultPacker(),
	}
}

func delFromList(list *list.List, np *NetPoint) *NetPoint {
	var npCheck *NetPoint
	for e := list.Front(); e != nil; e = e.Next() {
		npCheck = e.Value.(*NetPoint)
		//xsf_log.Info("delFromList", xsf_log.Uint32("id", npCheck.ID), xsf_log.Uint8("type", npCheck.SID.Type), xsf_log.Uint8("index", npCheck.SID.Index))
		if npCheck == np {
			//xsf_log.Info("delFromList remove", xsf_log.Uint32("id", np.ID))
			list.Remove(e)
			return npCheck
		}
	}

	return nil
}

func (mn *NPManagerArray) New(conn *net.TCPConn) *NetPoint {
	point, _ := NewNetPoint(conn, mn.packer)
	mn.tempConns.PushBack(point)
	point.SetNetEventRPC(mn.rpcAcceptor)
	return point
}

func (mn *NPManagerArray) Get(index uint32) *NetPoint {
	if index >= mn.maxCount {
		return nil
	}

	return mn.nps[index]
}

func (mn *NPManagerArray) Add(np *NetPoint) bool {
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

	return true
}

func (mn *NPManagerArray) Del(np *NetPoint) {
	npCheck := delFromList(mn.tempConns, np)
	if npCheck != nil {
		//xsf_log.Debug("NPManager_Array delete from tempConns")
	} else {
		if int(np.SID.Index) >= int(mn.maxCount) || np.SID.Index == 0 {
			//xsf_log.Error("NPManager_Array Del index error", xsf_log.Int("mn.maxCount", int(mn.maxCount)), xsf_log.Int("index", int(np.SID.Index)))
			return
		}

		mn.nps[np.SID.Index] = nil
	}

	//xsf_log.Info(">>>> server logout", xsf_log.Int("remote id", int(np.SID.ID)),
	//	xsf_log.Int("remote type", int(np.SID.Type)),
	//	xsf_log.Int("remote index", int(np.SID.Index)))
}

func (mn *NPManagerArray) OnUpdate() {
	nCurTime := time.Now().Unix()

	var next *list.Element
	for e := mn.tempConns.Front(); e != nil; e = next {
		next = e.Next()
		np := e.Value.(*NetPoint)
		if uint32(nCurTime) > np.Time+1000 { // 关闭超时的连接
			//xsf_log.Info("NPManager_Normal OnUpdate timeout", xsf_log.Uint("nCurTime", uint(nCurTime)), xsf_log.Uint("np.time", uint(np.time)))
			fmt.Println("NPManagerArray OnUpdate np.Close !!!")
			np.Close()
			mn.tempConns.Remove(e)
		}
	}
}

func (mn *NPManagerArray) OnHandshake(np *NetPoint) bool {
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

		//np.module.Kernel.(*NetKernel).Handler.OnNPAdd(np)

	} else {
		logger.Info("NPManager_Array OnHandshake mn.Add error", zap.Int("remote id", int(np.SID.ID)),
			zap.Int("remote type", int(np.SID.Type)), zap.Int("remote index", int(np.SID.Index)))
		fmt.Println("NPManagerArray OnHandshake Add np.Close !!!")
		np.Close()
	}

	return ok
}

func (mn *NPManagerArray) Broadcast(msgId int32, msg proto.Message, skip uint32) {

	for _, v := range mn.nps {
		if v != nil && v.ID != skip {
			packData, _ := v.Pack(uint64(msgId), msg)
			v.SetSignal(packData)
		}
	}

	//xsf_net.ReturnMessage(msg)
}
