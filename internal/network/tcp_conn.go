package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
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
	OnMessage(*Message, *TcpConn)
}

type TcpConn struct {
	Conn        net.Conn
	closed      int32
	ConnID      int64
	wg          sync.WaitGroup
	stopped     chan bool
	signal      chan interface{}
	lastSignal  chan interface{}
	msgParser   *BufferPacker
	timeoutTime int
	verify      int32
	Impl        IConn
	logger      *zap.Logger
}

func NewTcpConn(conn *net.TCPConn, msgBuffSize int, log *zap.Logger) (*TcpConn, error) {
	return &TcpConn{
		Conn:        conn,
		closed:      -1,
		verify:      0,
		msgParser:   newInActionPacker(),
		stopped:     make(chan bool, 1),
		signal:      make(chan interface{}, 100),
		lastSignal:  make(chan interface{}, 1),
		timeoutTime: 30,
		logger:      log,
	}, nil

}

func (c *TcpConn) Verified() bool {
	return atomic.LoadInt32(&c.verify) != 0
}

// Verify Client使用
func (c *TcpConn) Verify() {
	atomic.CompareAndSwapInt32(&c.verify, 0, 1)
}

// Reset Client使用
func (c *TcpConn) Reset() {
	if atomic.LoadInt32(&c.closed) == -1 {
		return
	}
	c.closed = -1
	c.verify = 0
	c.stopped = make(chan bool, 1)
	//c.signal = make(chan interface{}, c.msgBuffSize)
	c.lastSignal = make(chan interface{}, 1)
	c.msgParser.reset()
}

func (c *TcpConn) Connect() {
	if atomic.CompareAndSwapInt32(&c.closed, -1, 0) {

		c.wg.Add(1)
		go c.HandleRead()
		c.wg.Add(1)
		go c.HandleWrite()
	}

	timeout := time.NewTimer(time.Second * time.Duration(c.timeoutTime))

	defer func() {
		timeout.Stop()
		c.wg.Wait()
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
		case <-c.stopped:
			break L
		}
	}

}

func (c *TcpConn) Close() {
	if atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		c.Conn.Close()
		close(c.stopped)
	}
}

func (c *TcpConn) Read(b []byte) (int, error) {
	return c.Conn.Read(b)
}

func (c *TcpConn) Write(b []byte) (int, error) {
	if c.closed == -1 {
		return 0, nil
	}
	//fmt.Println("write:", b)
	return c.Conn.Write(b)
}
func (c *TcpConn) HandleRead() {
	defer func() {
		if err := recover(); err != nil {
			if err, ok := err.(error); ok {
				c.logger.Error("Recovered error", zap.Error(err))
			} else {
				c.logger.Error("Recovered value is not an error")
			}
			//logger.Error("[HandleRead] panic ", err, "\n", string(debug.Stack()))
		}
	}()
	defer c.Close()

	defer c.wg.Done()

	for {
		data, err := c.msgParser.Read(c)
		//data := make([]byte, 1024)
		//read, err := c.Conn.Read(data)
		//fmt.Println("read:", read)
		fmt.Printf("%s,%s", "read", string(data))
		if err != nil {
			if err != io.EOF {
				//c.logger.ErrorF("read message error: %v", err)
				fmt.Println(err)
			}
			break
		}

		message, err := c.msgParser.Unpack(data)

		//pbMsg := &player.SCLogin{}
		//err = proto.Unmarshal(message.Data, pbMsg)
		//if err == nil {
		//	fmt.Println("proto:", pbMsg)
		//} else {
		//	fmt.Println("data:", string(message.Data))
		//}
		fmt.Println("message:", message)
		c.Impl.OnMessage(message, c)

	}
}

func (c *TcpConn) HandleWrite() {
	defer func() {
		if err := recover(); err != nil {
			if err, ok := err.(error); ok {
				c.logger.Error("Recovered error", zap.Error(err))
			} else {
				c.logger.Error("Recovered value is not an error")
			}
		}
	}()
	defer c.Close()
	defer c.wg.Done()
	for {
		select {
		case signal := <-c.signal: // 普通消息
			data, ok := signal.([]byte)
			if !ok {
				//c.logger.ErrorF("write message %v error: msg is not bytes", reflect.TypeOf(signal))
				return
			}
			err := c.msgParser.Write(c, data...)
			if err != nil {
				c.logger.Error("write message %v error: msg is not bytes", zap.Error(err))
				return
			}
		case signal := <-c.lastSignal: // 最后一个通知消息
			data, ok := signal.([]byte)
			if !ok {
				c.logger.Error("write message %v error: msg is not bytes")
				return
			}
			err := c.msgParser.Write(c, data...)
			if err != nil {
				c.logger.Error("write message %v error: msg is not bytes", zap.Error(err))
				return
			}
			time.Sleep(2 * time.Second)
			return
		case <-c.stopped: // 连接关闭通知
			return
		}
	}
}

func (c *TcpConn) SetSignal(byte []byte) {
	c.signal <- byte
}

// OnConnect ...
func (c *TcpConn) OnConnect() {
	//c.logger.DebugF("[OnConnect] 建立连接 local:%s remote:%s", c.LocalAddr(), c.RemoteAddr())
}

func (c *TcpConn) OnClose() {
	//c.logger.InfoF("[OnConnect] 断开连接 local:%s remote:%s", c.LocalAddr(), c.RemoteAddr())
}

func (c *TcpConn) IsClosed() bool {
	return atomic.LoadInt32(&c.closed) != 0
}

func (c *TcpConn) IsShutdown() bool {
	return atomic.LoadInt32(&c.closed) == 1
}

func (c *TcpConn) Signal(signal []byte) error {
	select {
	case c.signal <- signal:
		return nil
	default:
		{
			cmd := binary.LittleEndian.Uint16(signal[2:4])
			return fmt.Errorf("[Signal] buffer full blocking connID:%v cmd:%v", c.ConnID, cmd)
		}
	}
}

func (c *TcpConn) AsyncSend(msgID uint64, msg interface{}) bool {

	if c.IsShutdown() {
		return false
	}

	data, err := c.Pack(msgID, msg)
	if err != nil {
		c.logger.Error("[AsyncSend] Pack msgID:%v and msg to bytes error:%v", zap.Error(err))

		return false
	}

	if uint32(len(data)) > c.msgParser.maxMsgLen {
		c.logger.Error("[AsyncSend] 发送的消息包体过长 msgID:%v", zap.Uint64("msgID", msgID))
		return false
	}

	err = c.Signal(data)
	if err != nil {
		c.Close()
		c.logger.Error("%v", zap.Error(err))
		return false
	}

	return true
}

func (c *TcpConn) Pack(msgID uint64, msg interface{}) ([]byte, error) {
	data, err := c.msgParser.Pack(msgID, msg)
	if err != nil {
		c.logger.Error("[AsyncSend] Pack msgID:%v and msg to bytes error:%v", zap.Uint64("msgID", msgID))
		return data, err
	}
	return data, nil
}

func (c *TcpConn) PackWrite(msg interface{}) error {
	data, ok := msg.([]byte)
	if !ok {
		return errors.New("write message error:msg is not bytes")
	}
	err := c.msgParser.Write(c, data...)
	if err != nil {
		c.logger.Error("[PackWrite]  error:%v", zap.Error(err))
		return err
	}
	return nil
}

func (c *TcpConn) GetMessageIdByCmd(cmd string) messageId.MessageId {
	mid, ok := messageId.MessageId_value[cmd]
	if ok {
		fmt.Println("mid-2", messageId.MessageId(mid))
		return messageId.MessageId(mid)
	}
	return messageId.MessageId_None
}
