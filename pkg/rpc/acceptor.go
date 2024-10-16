package rpc

import (
	"fmt"
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/server"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"runtime"
	"strings"
	"time"
)

const LenStackBuf = 4096

// 调用数据
type CallInfo struct {
	f    interface{}
	args []interface{}
}

// RPC接收者
type Acceptor struct {
	functions map[string][]func([]interface{})
	ChanCall  chan *CallInfo
	close     chan bool
	l         int
}

func NewAcceptor(l int) *Acceptor {
	a := new(Acceptor)
	a.l = l
	a.functions = make(map[string][]func([]interface{}))
	a.ChanCall = make(chan *CallInfo, l)
	return a
}

func (a *Acceptor) Regist(id string, f func([]interface{})) {
	if _, ok := a.functions[id]; ok {
		logger.Error(fmt.Sprintf("function id %v: already registered", id))
		return
	}

	//a.functions[id] = f
	a.functions[id] = append(a.functions[id], f)
}

func (a *Acceptor) Close() {
	close(a.ChanCall)
}

func (a *Acceptor) Open() {
	a.ChanCall = make(chan *CallInfo, a.l)
}

func (a *Acceptor) IsFull() bool {
	return len(a.ChanCall) == cap(a.ChanCall)
}

func (a *Acceptor) exec(ci *CallInfo) (err error) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, LenStackBuf)
			l := runtime.Stack(buf, false)
			err = fmt.Errorf("%v: %s", r, buf[:l])
		}
	}()

	ci.f.(func([]interface{}))(ci.args)

	return nil
}

// 接收者调度函数
func (a *Acceptor) Execute(ci *CallInfo) {
	err := a.exec(ci)
	if err != nil {
		logger.Error(fmt.Sprintf("Acceptor.Exec err:%s ", err.Error()))
	}
}

// goroutine safe
func (a *Acceptor) Go(id string, args ...interface{}) {
	f := a.functions[id]
	if f == nil {
		//logger.Error("Acceptor.Go function not found, id=%v, %v", id, args)
		return
	}

	defer func() {
		if x := recover(); x != nil {
			fmt.Println("caught panic in ECBEncrypt()", x)
		}
	}()

	if a.IsFull() {
		return
	}

	for _, h := range f {
		a.ChanCall <- &CallInfo{
			f:    h,
			args: args,
		}
	}
}

func (a *Acceptor) Run() {
	async.Go(func() {
	OUTLABEL:
		for {
			select {
			//case <-npCloseChan:
			//	return
			case <-a.close:
				break OUTLABEL
			case callMsg, ok := <-a.ChanCall:
				if !ok {
					break OUTLABEL
					//return
				}
				if callMsg == nil {
					continue
				}
				a.Execute(callMsg)
			default:
				time.Sleep(100 * time.Millisecond)
			}

			//time.Sleep(100 * time.Millisecond)
		}
	})
}

func GetClientDestEP(id uint64) int {
	msgName := MsgId2Name(int32(id))
	if strings.Contains(msgName, "_G_") {
		return server.EP_Game
	} else if strings.Contains(msgName, "_L_") {
		return server.EP_Login
	} else if strings.Contains(msgName, "_Gt_") {
		return server.EP_Gate
	} else if strings.Contains(msgName, "_Ml_") {
		return server.EP_Mail
	} else if strings.Contains(msgName, "_Mg_") {
		return server.EP_Manager
	} else {
		return GetMsgEp(id)
	}
}

func GetMsgEp(id uint64) int {
	msgName := MsgId2Name(int32(id))
	if strings.Contains(msgName, "Clt_") {
		return server.EP_Client
	} else if strings.Contains(msgName, "G_") {
		return server.EP_Game
	} else if strings.Contains(msgName, "L_") {
		return server.EP_Login
	} else if strings.Contains(msgName, "Gt_") {
		return server.EP_Gate
	} else if strings.Contains(msgName, "C_") {
		return server.EP_Center
	} else if strings.Contains(msgName, "Mg_") {
		return server.EP_Manager
	} else if strings.Contains(msgName, "Ml_") {
		return server.EP_Mail
	} else {
		// fmt.Println("GetEP ep not handler, msgID=", msgID)
		return server.EP_None
	}
}

func MsgId2Name(msgId int32) string {
	if val, ok := xsf_pb.SMSGID_name[msgId]; ok {
		return val
	}
	if val, ok := xsf_pb.MSGID_name[msgId]; ok {
		return val
	}

	return ""
}

func PrintMsgLog(msgId uint64, data []byte, types string) {
	getMessage, _ := GetMessage(msgId)
	m := getMessage.(proto.Message)
	err := proto.Unmarshal(data, m)
	if err == nil {
		{
			msgName := xsf_pb.SMSGID_name[int32(msgId)]
			logger.Info(types+" message:", zap.String("Message", msgName), zap.String("message", m.String()))
		}
	}
}
