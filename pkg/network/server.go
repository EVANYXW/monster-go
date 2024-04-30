package network

import (
	"github.com/evanyxw/monster-go/internal/pkg/alert"
	"github.com/evanyxw/monster-go/internal/pkg/output"
	"github.com/evanyxw/monster-go/pkg/async"
	"go.uber.org/zap"
	"net"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type Info struct {
	ServerName string
	Env        string
	Address    string
	RpcAddr    string
}

type Server struct {
	Info
	alertNotify    alert.Handler
	logger         *zap.Logger
	swg            sync.WaitGroup
	ls             *net.TCPListener
	mux            sync.Mutex
	MessageHandler func(packet *Packet)
	maxConnNum     int
	connBuffSize   int
	pid            int64
	connMap        map[net.Conn]interface{}
	addr           string
	closeChan      chan struct{}
}

func NewServer(addr string, maxConnNum int, buffSize int, log *zap.Logger, info Info) *Server {
	s := &Server{
		addr:         addr,
		maxConnNum:   maxConnNum,
		connBuffSize: buffSize,
		connMap:      make(map[net.Conn]interface{}, 0),
		logger:       log,
		closeChan:    make(chan struct{}),
		alertNotify:  alert.NotifyHandler(log),
		Info:         info,
	}
	s.Init()

	return s
}

func (s *Server) Init() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", s.addr)
	if err != nil {
		s.logger.Fatal("[net] addr resolve error", zap.Error(err))
	}

	ln, err := net.ListenTCP("tcp4", tcpAddr)

	if err != nil {
		s.logger.Fatal("%v", zap.Error(err))
	}

	if s.maxConnNum <= 0 {
		s.maxConnNum = 100
		s.logger.Info("invalid MaxConnNum, reset to %v", zap.Int("maxConnNum", s.maxConnNum))
	}

	s.ls = ln
	s.pid = int64(os.Getpid())
	//s.logger.InfoF("Server Listen %s", s.ln.Addr().String())
}

func (s *Server) addConn(conn *net.TCPConn, tcpConn *TcpConn) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.connMap[conn] = tcpConn
	tcpConn.ConnID = 12323
	if output.Oput != nil {
		output.Oput.Update(output.Data{
			ConnNum: int32(len(s.connMap)),
			GoCount: async.GetGoCount(),
		})
	}
}

func (s *Server) removeConn(conn *net.TCPConn, tcpConn *TcpConn) {
	tcpConn.Close()
	s.mux.Lock()
	defer s.mux.Unlock()
	delete(s.connMap, conn)
	if output.Oput != nil {
		output.Oput.Update(output.Data{
			ConnNum: int32(len(s.connMap)),
			GoCount: async.GetGoCount(),
		})
	}
}

func (s *Server) Run() {
	defer func() {
		if err := recover(); err != nil {
			if err, ok := err.(error); ok {
				s.alertNotify(&alert.AlertMessage{
					ProjectName:  s.Info.ServerName,
					Env:          s.Info.Env,
					HOST:         s.Info.Address,
					ErrorMessage: err.Error, ErrorStack: string(debug.Stack()),
				})
				s.logger.Error("[net run] panic", zap.Error(err))
			} else {
				s.logger.Error("Recovered value is not an error")
			}
		}
	}()
	s.swg.Add(1)
	defer s.swg.Done()

	var tempDelay time.Duration
outer:
	for {
		select {
		case <-s.closeChan:
			break outer
		default:

		}

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
			return
		}

		tempDelay = 0

		if len(s.connMap) >= s.maxConnNum {
			conn.Close()
			s.logger.Info("too many connections %v", zap.Int("connections", len(s.connMap)))
			continue
		}

		tcpConn, err := NewTcpConn(conn, s.logger)
		if err != nil {
			continue
		}

		s.addConn(conn, tcpConn)
		tcpConn.Impl = s
		// s.swg.Add(1)
		go func() {
			tcpConn.Connect()
			s.removeConn(conn, tcpConn)
			// s.swg.Done()
		}()
	}
}

func (s *Server) OnMessage(message *Message, conn *TcpConn) {
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
				s.logger.Error("[OnMessage] panic", zap.Error(err))
			} else {
				s.logger.Error("Recovered value is not an error")
			}
		}
	}()
	if s.MessageHandler != nil {
		s.MessageHandler(&Packet{
			Msg:  message,
			Conn: conn,
		})
	}
}

func (s *Server) OnConnect() {

}

func (s *Server) OnClose() {
	s.closeChan <- struct{}{}
}
