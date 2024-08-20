package network

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"
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
	Stopped    chan bool
	signal     chan interface{}
	lastSignal chan interface{}
	//msgParser   *BufferPacker
	msgParser   Packer
	timeoutTime int
	verify      int32
	Impl        IConn
	Time        uint32

	SID                server.ServerID
	ID                 uint32
	senderID           uint32 // 发送的服务器id,暂时没有用
	RemoteIP           string
	isHandshake        bool
	lastHeartbeat      uint32
	RpcAcceptor        *rpc.Acceptor
	IsRpcAcceptorClose bool
	Processor          *Processor
	CloseChan          chan bool
}

func NewNetPoint(conn *net.TCPConn, packerFactory PackerFactory) (*NetPoint, error) {
	return &NetPoint{
		Conn:      conn,
		closed:    -1,
		verify:    0,
		msgParser: packerFactory.CreatePacker(),
		//msgParser:   packer,
		Stopped:     make(chan bool, 1),
		signal:      make(chan interface{}, 100),
		lastSignal:  make(chan interface{}, 1),
		timeoutTime: 30,
		RemoteIP:    conn.RemoteAddr().(*net.TCPAddr).IP.String(),
		CloseChan:   make(chan bool),
	}, nil

}

func (np *NetPoint) SetProcessor(pro *Processor) {
	np.Processor = pro
}

func (np *NetPoint) SetNetEventRPC(rpc *rpc.Acceptor) {
	np.RpcAcceptor = rpc
}

func (np *NetPoint) SetCloseAccept(set bool) {
	np.IsRpcAcceptorClose = set
}

func (np *NetPoint) SetID(id uint32) {
	np.ID = id
	server.ID2Sid(id, &np.SID)
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
	np.Stopped = make(chan bool, 1)
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
		case <-np.Stopped:
			break L
		}
		time.Sleep(100 * time.Millisecond)
	}

}

func (np *NetPoint) Close() {
	marshal, _ := json.Marshal(np)
	logger.Info("NetPoint Close", zap.String("np:", string(marshal)))
	if atomic.CompareAndSwapInt32(&np.closed, 0, 1) {
		logger.Info("NetPoint Close Success", zap.String("np:", string(marshal)))
		np.Conn.Close()
		async.Go(func() {
			np.Stopped <- true
		})
		close(np.Stopped)
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

OUTLABEL:
	for {
		select {
		case <-np.Stopped:
			// fixMe 如何gate的client 设置了新的rpcAcceptor并且run起来会有2个chan在跑
			np.RpcAcceptor.Go(rpc.RPC_NET_ERROR, np, &Acceptor{})
			break OUTLABEL
		default:
			//message, err := c.msgParser.TestRead(c)
			data, err := np.msgParser.Read(np)
			if err != nil {
				if err != io.EOF && err.Error() != "EOF readLen:0" {
					np.RpcAcceptor.Go(rpc.RPC_NET_ERROR, np, &Acceptor{})
					logger.Debug("read message RPC_NET_ERROR: %v", zap.Error(err), zap.Uint32("server_id", np.ID))
				}
				np.RpcAcceptor.Go(rpc.RPC_NET_ERROR, np, &Acceptor{})
				break OUTLABEL
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
			//np.RpcAcceptor.Go(rpc.RPC_NET_MESSAGE, np, message)
		}
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

OutLabel:
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

			return
		case <-np.Stopped: // 连接关闭通知
			break OutLabel
		}

		time.Sleep(100 * time.Millisecond)
	}
	logger.Info("退出写了")
}

func (np *NetPoint) SendMessage(message proto.Message, options ...PackerOptions) {
	pack, _ := np.Pack(message, options...)
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

func (np *NetPoint) AsyncSend(msg interface{}) bool {

	if np.IsShutdown() {
		return false
	}

	data, err := np.Pack(msg)

	if err != nil {
		logger.Error("[AsyncSend] Pack msgID:%v and msg to bytes error:%v", zap.Error(err))

		return false
	}

	if uint32(len(data)) > np.msgParser.GetMaxMsgLen() {
		localMsg := msg.(proto.Message)
		msgID, _ := rpc.GetMsgID(localMsg)
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

func (np *NetPoint) Pack(msg interface{}, options ...PackerOptions) ([]byte, error) {
	data, err := np.msgParser.Pack(msg, options...)
	if err != nil {
		localMsg := msg.(proto.Message)
		msgID, _ := rpc.GetMsgID(localMsg)
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

func (np *NetPoint) OnHandshakeTicker(netPoint *NetPoint) {
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

	np.lastHeartbeat = uint32(time.Now().Unix())
	async.Go(func() {
		cnf := configs.Get()
		timer := time.NewTimer(time.Duration(cnf.HtCheck))
		heartbeatTimeout := uint64(time.Duration(cnf.HtTimeout))
		defer func() {
			timer.Stop()
		}()
		for {
			select {
			case <-timer.C:
				curTime := time.Now().Unix()
				if uint64(curTime) > uint64(np.lastHeartbeat)+heartbeatTimeout {
					fmt.Println("心跳未回，被T掉")
					np.Close() // evan
					return
				}
			default:
				time.Sleep(500 * time.Millisecond)
			}
			_ = timer.Reset(time.Duration(cnf.HtCheck)) //重制心跳上报时间间隔
		}
	})

}

func (np *NetPoint) OnHandshakeTicker2(netPoint *NetPoint) {
	np.isHandshake = true
	np.lastHeartbeat = uint32(time.Now().Unix())
	// fixMe 性能比较差，吃CPU
	async.Go(func() {
		cnf := configs.Get()
		ticker := time.NewTicker(time.Duration(cnf.HtCheck))
		heartbeatTimeout := uint64(time.Duration(cnf.HtTimeout))
		defer ticker.Stop()

		for range ticker.C {
			curTime := time.Now().Unix()
			if uint64(curTime) > uint64(np.lastHeartbeat)+heartbeatTimeout {
				fmt.Println("心跳未回，被T掉")
				np.Close() // evan
				return
			}
		}
	})

}
