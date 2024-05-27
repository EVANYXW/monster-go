package network

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/utils"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type IConn interface {
	OnConnect()
	OnClose()
	OnMessage(*Message, *NetPoint)
}

type NetPoint struct {
	Conn       net.Conn
	closed     int32
	wg         sync.WaitGroup
	stopped    chan bool
	signal     chan interface{}
	lastSignal chan interface{}
	//msgParser   *BufferPacker
	msgParser   Packer
	timeoutTime int
	verify      int32
	Impl        IConn
	//logger      *zap.Logger
	Time uint32

	SID           ServerID
	ID            uint32
	RemoteIP      string
	isHandshake   bool
	lastHeartbeat uint32
	RpcAcceptor   *rpc.Acceptor
}

func NewNetPoint(conn *net.TCPConn) (*NetPoint, error) {
	return &NetPoint{
		Conn:        conn,
		closed:      -1,
		verify:      0,
		msgParser:   newDefaultPacker(),
		stopped:     make(chan bool, 1),
		signal:      make(chan interface{}, 100),
		lastSignal:  make(chan interface{}, 1),
		timeoutTime: 30,
		RemoteIP:    conn.RemoteAddr().(*net.TCPAddr).IP.String(),
	}, nil

}

func (np *NetPoint) SetID(id uint32) {
	np.ID = id
	ID2Sid(id, &np.SID)
}

func (np *NetPoint) Verified() bool {
	return atomic.LoadInt32(&np.verify) != 0
}

// Verify Client使用
func (np *NetPoint) Verify() {
	atomic.CompareAndSwapInt32(&np.verify, 0, 1)
}

// Reset Client使用
func (np *NetPoint) Reset() {
	if atomic.LoadInt32(&np.closed) == -1 {
		return
	}
	np.closed = -1
	np.verify = 0
	np.stopped = make(chan bool, 1)
	//c.signal = make(chan interface{}, c.msgBuffSize)
	np.lastSignal = make(chan interface{}, 1)
	np.msgParser.Reset()
}

func (np *NetPoint) Connect() {
	if atomic.CompareAndSwapInt32(&np.closed, -1, 0) {
		np.wg.Add(1)
		go np.HandleRead()
		np.wg.Add(1)
		go np.HandleWrite()
	}

	timeout := time.NewTimer(time.Second * time.Duration(np.timeoutTime))

	defer func() {
		timeout.Stop()
		np.wg.Wait()
	}()

L:
	for {
		select {
		// 等待通到返回 返回后检查连接是否验证完成 如果没有验证 则关闭连接
		case <-timeout.C:
			//if !c.Verified() {
			//	//c.logger.ErrorF("[Connect] 验证超时 ip addr %s", c.RemoteAddr())
			//	c.Close()
			//	break L
			//}
		case <-np.stopped:
			break L
		}
	}

}

func (np *NetPoint) Close() {
	marshal, _ := json.Marshal(np)
	logger.Info("NetPoint Close", zap.String("np:", string(marshal)))
	if atomic.CompareAndSwapInt32(&np.closed, 0, 1) {
		logger.Info("NetPoint Close Success", zap.String("np:", string(marshal)))
		np.Conn.Close()
		close(np.stopped)
	}
}

func (np *NetPoint) Read(b []byte) (int, error) {
	return np.Conn.Read(b)
}

func (np *NetPoint) Write(b []byte) (int, error) {
	if np.closed == -1 {
		return 0, nil
	}
	return np.Conn.Write(b)
}
func (np *NetPoint) HandleRead() {
	defer func() {
		if err := recover(); err != nil {
			if err, ok := err.(error); ok {
				logger.Error("Recovered error", zap.Error(err))
			} else {
				logger.Error("Recovered value is not an error")
			}

			//logger.Error("[HandleRead] panic ", err, "\n", string(debug.Stack()))
		}
	}()
	defer np.Close()

	defer np.wg.Done()

	for {
		//message, err := c.msgParser.TestRead(c)
		data, err := np.msgParser.Read(np)

		if err != nil {
			if err != io.EOF {
				np.RpcAcceptor.Go(RPC_NET_ERROR, np)
				logger.Error("read message RPC_NET_ERROR: %v", zap.Error(err))
			}
			break
		}

		message, err := np.msgParser.Unpack(data)
		rpc.PrintMsgLog(message.ID, message.Data, "read")

		//pbMsg := &player.SCLogin{}
		//err = proto.Unmarshal(message.Data, pbMsg)
		//if err == nil {
		//	fmt.Println("proto:", pbMsg)
		//} else {
		//	fmt.Println("data:", string(message.Data))
		//}
		np.Impl.OnMessage(message, np)

	}
}

func (np *NetPoint) HandleWrite() {
	defer func() {
		if err := recover(); err != nil {
			if err, ok := err.(error); ok {
				logger.Error("Recovered error", zap.Error(err))
			} else {
				logger.Error("Recovered value is not an error")
			}
		}
	}()
	defer np.Close()
	defer np.wg.Done()
	defer close(np.signal)
	defer close(np.lastSignal)

	for {
		select {
		case signal := <-np.signal: // 普通消息

			data, ok := signal.([]byte)
			if !ok {
				//c.logger.ErrorF("write message %v error: msg is not bytes", reflect.TypeOf(signal))
				return
			}

			//buffer := bytes.NewBuffer(data)
			//message, err := c.msgParser.Read(buffer)
			//unpack, err := c.msgParser.Unpack(message)
			//rpc.PrintMsgLog(unpack.ID, unpack.Data, "read")

			err := np.msgParser.Write(np, data...)
			if err != nil {
				logger.Error("write message %v error: msg is not bytes", zap.Error(err))
				return
			}

		case signal := <-np.lastSignal: // 最后一个通知消息
			data, ok := signal.([]byte)
			if !ok {
				logger.Error("write message %v error: msg is not bytes")
				return
			}
			err := np.msgParser.Write(np, data...)
			if err != nil {
				logger.Error("write message %v error: msg is not bytes", zap.Error(err))
				return
			}
			time.Sleep(2 * time.Second)
			return
		case <-np.stopped: // 连接关闭通知
			return
		}
	}
}

func (np *NetPoint) SendMessage(msgId uint64, message proto.Message) {
	pack, _ := np.Pack(msgId, message)
	np.SetSignal(pack)
}

func (np *NetPoint) SetSignal(byte []byte) {
	np.signal <- byte
}

// OnConnect ...
func (np *NetPoint) OnConnect() {
	//c.logger.DebugF("[OnConnect] 建立连接 local:%s remote:%s", c.LocalAddr(), c.RemoteAddr())
}

func (np *NetPoint) OnClose() {
	//c.logger.InfoF("[OnConnect] 断开连接 local:%s remote:%s", c.LocalAddr(), c.RemoteAddr())
}

func (np *NetPoint) IsClosed() bool {
	return atomic.LoadInt32(&np.closed) != 0
}

func (np *NetPoint) IsShutdown() bool {
	return atomic.LoadInt32(&np.closed) == 1
}

func (np *NetPoint) Signal(signal []byte) error {
	select {
	case np.signal <- signal:
		return nil
	default:
		{
			cmd := binary.LittleEndian.Uint16(signal[2:4])
			return fmt.Errorf("[Signal] buffer full blocking connID:%v cmd:%v", np.ID, cmd)
		}
	}
}

func (np *NetPoint) AsyncSend(msgID uint64, msg interface{}) bool {

	if np.IsShutdown() {
		return false
	}

	data, err := np.Pack(msgID, msg)

	if err != nil {
		logger.Error("[AsyncSend] Pack msgID:%v and msg to bytes error:%v", zap.Error(err))

		return false
	}

	if uint32(len(data)) > np.msgParser.GetMaxMsgLen() {
		logger.Error("[AsyncSend] 发送的消息包体过长 msgID:%v", zap.Uint64("msgID", msgID))
		return false
	}

	err = np.Signal(data)
	if err != nil {
		np.Close()
		logger.Error("%v", zap.Error(err))
		return false
	}

	return true
}

func (np *NetPoint) Pack(msgID uint64, msg interface{}) ([]byte, error) {
	data, err := np.msgParser.Pack(msgID, msg)
	if err != nil {
		logger.Error("[AsyncSend] Pack msgID:%v and msg to bytes error:%v", zap.Uint64("msgID", msgID))
		return data, err
	}
	return data, nil
}

func (np *NetPoint) PackWrite(msg interface{}) error {
	data, ok := msg.([]byte)
	if !ok {
		return errors.New("write message error:msg is not bytes")
	}
	err := np.msgParser.Write(np, data...)
	if err != nil {
		logger.Error("[PackWrite]  error:%v", zap.Error(err))
		return err
	}
	return nil
}

//func (c *NetPoint) GetMessageIdByCmd(cmd string) messageId.MessageId {
//	mid, ok := messageId.MessageId_value[cmd]
//	if ok {
//		return messageId.MessageId(mid)
//	}
//	return messageId.MessageId_None
//}

func (np *NetPoint) OnHeartbeat() {
	np.lastHeartbeat = utils.NP_CurrentS()
}

func (np *NetPoint) OnHandshake() {
	np.isHandshake = true
	//np.t = np.module.Dispatcher.AfterFunc(xsf_config.HeartbeatCheck, func() {
	//	//xsf_log.Debug("OnHandshake heartbeat check ..")
	//	curTime := xsf_util.NP_CurrentS()
	//	if curTime > np.lastHeartbeat+xsf_config.HeartbeatTimeout {
	//		xsf_log.Info("NetPoint heartbeat timeout", xsf_log.Uint("last", uint(np.lastHeartbeat)),
	//			xsf_log.Uint("curTime", uint(curTime)),
	//			xsf_log.Uint("server", uint(np.SID.ID)),
	//			xsf_log.Uint("type", uint(np.SID.Type)),
	//			xsf_log.Uint("index", uint(np.SID.Index)))
	//		np.Close()
	//	}
	//}, -1)

	//np.lastHeartbeat = uint32(time.Now().Unix() + 10)
	async.Go(func() {
		ticker := time.NewTicker(7 * time.Second)
		heartbeatTimeout := uint32(60)
		defer ticker.Stop()
		for range ticker.C {
			curTime := time.Now().Unix()
			if uint32(curTime) > np.lastHeartbeat+heartbeatTimeout {
				//xsf_log.Info("NetPoint heartbeat timeout", xsf_log.Uint("last", uint(np.lastHeartbeat)),
				//	xsf_log.Uint("curTime", uint(curTime)),
				//	xsf_log.Uint("server", uint(np.SID.ID)),
				//	xsf_log.Uint("type", uint(np.SID.Type)),
				//	xsf_log.Uint("index", uint(np.SID.Index)))
				fmt.Println("心跳未回，被T掉")
				np.Close()
				return
			}
		}
	})

}

func (np *NetPoint) SetNetEventRPC(rpc *rpc.Acceptor) {
	np.RpcAcceptor = rpc
}
