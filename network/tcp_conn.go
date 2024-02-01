package network

import (
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
	ConnID      int64
	wg          sync.WaitGroup
	stopped     chan bool
	signal      chan interface{}
	lastSignal  chan interface{}
	msgParser   *BufferPacker
	timeoutTime int
	verify      int32
	Impl        IConn
}

func NewTcpConn(conn *net.TCPConn, msgBuffSize int) (*TcpConn, error) {
	return &TcpConn{
		Conn:        conn,
		verify:      0,
		msgParser:   newInActionPacker(),
		stopped:     make(chan bool, 1),
		signal:      make(chan interface{}, 100),
		lastSignal:  make(chan interface{}, 1),
		timeoutTime: 30,
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
	//if atomic.LoadInt32(&c.closed) == -1 {
	//	return
	//}
	//c.closed = -1
	c.verify = 0
	c.stopped = make(chan bool, 1)
	//c.signal = make(chan interface{}, c.msgBuffSize)
	c.lastSignal = make(chan interface{}, 1)
	c.msgParser.reset()
}

func (c *TcpConn) Connect() {
	c.wg.Add(1)
	go c.HandleRead()
	c.wg.Add(1)
	go c.HandleWrite()

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

func (c *TcpConn) HandleRead() {
	defer func() {
		if err := recover(); err != nil {
			//c.logger.ErrorF("[HandleRead] panic ", err, "\n", string(debug.Stack()))
		}
	}()
	defer c.Close()

	defer c.wg.Done()

	for {
		data, err := c.msgParser.Read(c)
		if err != nil {
			if err != io.EOF {
				//c.logger.ErrorF("read message error: %v", err)
			}
			break
		}
		message, err := c.msgParser.Unpack(data)
		c.Impl.OnMessage(message, c)
	}
}

func (c *TcpConn) Close() {
	//if atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
	c.Conn.Close()
	close(c.stopped)
	//}
}

func (c *TcpConn) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (c *TcpConn) HandleWrite() {
	defer func() {
		if err := recover(); err != nil {
			//c.logger.ErrorF("[HandleWrite] panic", err, "\n", string(debug.Stack()))
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
				//c.logger.ErrorF("write message %v error: %v", reflect.TypeOf(signal), err)
				return
			}
		case signal := <-c.lastSignal: // 最后一个通知消息
			data, ok := signal.([]byte)
			if !ok {
				//c.logger.ErrorF("write message %v error: msg is not bytes", reflect.TypeOf(signal))
				return
			}
			err := c.msgParser.Write(c, data...)
			if err != nil {
				//c.logger.ErrorF("write message %v error: %v", reflect.TypeOf(signal), err)
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
