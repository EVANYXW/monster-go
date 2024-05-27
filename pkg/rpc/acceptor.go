package rpc

import (
	"fmt"
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"reflect"
	"runtime"
)

const LenStackBuf = 4096

// 调用数据
type CallInfo struct {
	f    interface{}
	args []interface{}
}

// RPC接收者
type Acceptor struct {
	functions map[string]func([]interface{})
	ChanCall  chan *CallInfo
}

func NewAcceptor(l int) *Acceptor {
	a := new(Acceptor)
	a.functions = make(map[string]func([]interface{}))
	a.ChanCall = make(chan *CallInfo, l)
	return a
}

func (a *Acceptor) Regist(id string, f func([]interface{})) {
	if _, ok := a.functions[id]; ok {
		logger.Error(fmt.Sprintf("function id %v: already registered", id))
		return
	}

	a.functions[id] = f
}

func (a *Acceptor) Close() {
	close(a.ChanCall)
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
		recover()
	}()

	a.ChanCall <- &CallInfo{
		f:    f,
		args: args,
	}
}

func (a *Acceptor) Run() {
	async.Go(func() {
		for {
			select {
			case callMsg := <-a.ChanCall:
				a.Execute(callMsg)
			}
		}
	})
}

func GetMessage(messageID uint64) (interface{}, error) {
	// 利用反射获取消息 ID 对应的协议结构体
	//messageType2 := reflect.TypeOf((*proto.Message)(nil)).Elem()
	for _, messageValue := range xsf_pb.SMSGID_value {
		if messageValue == int32(messageID) {
			messageName := xsf_pb.SMSGID_name[int32(messageValue)]
			messageType := proto.MessageType("NLD_PB." + messageName)
			if messageType == nil {
				return nil, fmt.Errorf("未找到对应的协议结构体")
			}

			instance := reflect.New(messageType.Elem()).Interface()
			msg, ok := instance.(proto.Message)
			if !ok {
				return nil, fmt.Errorf("实例化失败")
			}

			return msg, nil
		}
	}

	return nil, fmt.Errorf("未找到对应的协议结构体")
}

func PrintMsgLog(msgId uint64, data []byte, types string) {
	getMessage, _ := GetMessage(msgId)
	m := getMessage.(proto.Message)
	err := proto.Unmarshal(data, m)
	if err == nil {
		{
			logger.Info(types+" message:", zap.Uint64("Message ID", msgId), zap.String("message", m.String()))
		}
	}
}
