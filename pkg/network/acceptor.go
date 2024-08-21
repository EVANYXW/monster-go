package network

import (
	"fmt"
	"github.com/evanyxw/monster-go/pkg/alert"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/output"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"
	"go.uber.org/zap"
	"net"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type Options func(acceptor *Acceptor)

func WithAddr(addr string) Options {
	return func(acceptor *Acceptor) {
		acceptor.Addr = addr
	}
}

type Acceptor struct {
	server.Info
	alertNotify    alert.Handler
	swg            sync.WaitGroup
	ls             *net.TCPListener
	mux            sync.Mutex
	MessageHandler func(packet *Packet)
	maxConnNum     uint32
	//connBuffSize   int
	pid       int64
	connMap   map[net.Conn]interface{}
	NPManager INPManager
	Addr      string
	closeChan chan struct{}
	//Owner     module.NetKernel
}

func NewAcceptor(maxConnNum uint32, info server.Info, nodePointManager INPManager) *Acceptor {
	s := &Acceptor{
		Addr:       fmt.Sprintf("%s:%d", info.Ip, info.Port),
		maxConnNum: maxConnNum,
		//connBuffSize: buffSize,
		connMap:     make(map[net.Conn]interface{}, 0),
		closeChan:   make(chan struct{}),
		alertNotify: alert.NotifyHandler(),
		Info:        info,
		NPManager:   nodePointManager,
	}

	return s
}

func (s *Acceptor) Connect(options ...Options) {
	for _, option := range options {
		option(s)
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp4", s.Addr)
	if err != nil {
		logger.Fatal("[net] addr resolve error", zap.Error(err))
	}

	ln, err := net.ListenTCP("tcp4", tcpAddr)

	if err != nil {
		logger.Fatal("%v", zap.Error(err))
	}

	logger.Info(fmt.Sprintf("新建立网络模块：%s", s.Addr))
	if s.maxConnNum <= 0 {
		s.maxConnNum = 100
		logger.Info("invalid MaxConnNum, reset to %v", zap.Int("maxConnNum", int(s.maxConnNum)))
	}

	s.ls = ln
	s.pid = int64(os.Getpid())
	//s.logger.InfoF("Server Listen %s", s.ln.Addr().String())
}

//func (s *Acceptor) DoStart() {
//	port := Ports[EP_Client]
//	addr := fmt.Sprintf(":%d", port)
//	output.Oput.SetServerAddr(addr)
//
//	async.Go(func() {
//		s.Connect(WithAddr(addr))
//		s.Run()
//	})
//}

func (s *Acceptor) addConn(conn *net.TCPConn) *NetPoint {
	s.mux.Lock()
	defer s.mux.Unlock()
	point := s.NPManager.New(conn)
	s.connMap[conn] = point

	//point.RpcAcceptor.Run(point.Stopped) //CloseChan
	point.RpcAcceptor.Run() // 这个应该只Run一个

	if output.Oput != nil {
		output.Oput.SetData(output.Data{
			ConnNum: int32(len(s.connMap)),
			GoCount: async.GetGoCount(),
		})
	}
	return point
}

func (s *Acceptor) RemoveConn(conn *net.TCPConn, tcpConn *NetPoint) {
	fmt.Println("RemoveConn RemoveConn !!!")
	tcpConn.Close()
	s.mux.Lock()
	defer s.mux.Unlock()
	delete(s.connMap, conn)
	if output.Oput != nil {
		output.Oput.SetData(output.Data{
			ConnNum: int32(len(s.connMap)),
			GoCount: async.GetGoCount(),
		})
	}
}

func (s *Acceptor) Run() {
	defer func() {
		if err := recover(); err != nil {
			if err, ok := err.(error); ok {
				s.alertNotify(&alert.AlertMessage{
					ProjectName:  s.Info.ServerName,
					Env:          s.Info.Env,
					HOST:         fmt.Sprintf("%s:%d", s.Info.Ip, s.Info.Port),
					ErrorMessage: err.Error, ErrorStack: string(debug.Stack()),
				})
				logger.Error("[net run] panic", zap.Error(err))
			} else {
				logger.Error("Recovered value is not an error")
			}
		}
	}()
	s.swg.Add(1)
	defer s.swg.Done()
	// fixMe 这个携程回泄漏
	//s.RpcAcceptor.Run()

	var tempDelay time.Duration
outer:
	for {
		select {
		case <-s.closeChan:
			break outer
		default:
			conn, err := s.ls.AcceptTCP()
			if err != nil {
				if _, ok := err.(net.Error); ok {
					if tempDelay == 0 {
						tempDelay = 5 * time.Millisecond
					} else {
						tempDelay *= 2
					}
					if max := 1 * time.Second; tempDelay > max {
						tempDelay = max
					}
					// s.logger.InfoF("accept error: %v; retrying in %v", err, tempDelay)
					time.Sleep(tempDelay)
					continue
				}
				break outer
			}

			tempDelay = 0

			if uint32(len(s.connMap)) >= s.maxConnNum {
				conn.Close()
				logger.Info("too many connections %v", zap.Int("connections", len(s.connMap)))
				continue
			}

			point := s.addConn(conn)
			point.Impl = s

			logger.Info("net point is create")
			point.RpcAcceptor.Go(rpc.RPC_NET_ACCEPT, point, s)

			//go func() {
			//	point.Connect()
			//	s.removeConn(conn, point)
			//
			//}()
		}

	}
}

func (s *Acceptor) OnMessage(message *Message, conn *NetPoint) {
	defer func() {
		if err := recover(); err != nil {
			s.alertNotify(&alert.AlertMessage{
				ProjectName:  s.Info.ServerName,
				Env:          s.Info.Env,
				HOST:         s.Info.Address,
				ErrorMessage: err,
				ErrorStack:   string(debug.Stack()),
			})
			if err, ok := err.(error); ok {
				logger.Error("[OnMessage] panic", zap.Error(err))
			} else {
				logger.Error("Recovered value is not an error")
			}
		}
	}()
	if s.MessageHandler != nil {
		s.MessageHandler(&Packet{
			Msg:      message,
			NetPoint: conn,
		})
	}
}

func (s *Acceptor) OnConnect() {

}

func (s *Acceptor) OnClose() {
	s.closeChan <- struct{}{}
}
