package network

import (
	"go.uber.org/zap"
	"net"
	"os"
	"sync"
	"time"
)

type Server struct {
	addr           string
	maxConnNum     int
	connBuffSize   int
	swg            sync.WaitGroup
	ls             *net.TCPListener
	connMap        map[net.Conn]interface{}
	mux            sync.Mutex
	pid            int64
	MessageHandler func(packet *Packet)
	logger         *zap.Logger
	closeChan      chan struct{}
}

func NewServer(addr string, maxConnNum int, buffSize int, log *zap.Logger) *Server {
	s := &Server{
		addr:         addr,
		maxConnNum:   maxConnNum,
		connBuffSize: buffSize,
		connMap:      make(map[net.Conn]interface{}, 0),
		logger:       log,
		closeChan:    make(chan struct{}),
	}
	s.Init()

	return s
}

func (s *Server) Init() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", s.addr)

	if err != nil {
		//s.logger.FatalF("[net] addr resolve error", tcpAddr, err)
	}

	//ln, err := net.ListenTCP("tcp6", tcpAddr)
	ln, err := net.ListenTCP("tcp4", tcpAddr)

	if err != nil {
		//s.logger.FatalF("%v", err)
	}

	if s.maxConnNum <= 0 {
		s.maxConnNum = 100
		//s.logger.InfoF("invalid MaxConnNum, reset to %v", s.MaxConnNum)
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
}

func (s *Server) removeConn(conn *net.TCPConn, tcpConn *TcpConn) {
	tcpConn.Close()
	s.mux.Lock()
	defer s.mux.Unlock()
	delete(s.connMap, conn)
}

func (s *Server) Run() {
	defer func() {
		if err := recover(); err != nil {
			//s.logger.ErrorF("[net] panic", err, "\n", string(debug.Stack()))
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
		s.logger.Info("lai ren la")
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
			// s.logger.InfoF("too many connections %v", atomic.LoadInt64(&s.counter))
			continue
		}

		tcpConn, err := NewTcpConn(conn, s.connBuffSize, s.logger)
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
	if s.MessageHandler != nil {
		s.MessageHandler(&Packet{
			Msg:  message,
			Conn: conn,
		})
	}

}

func (s *Server) OnClose() {
	s.closeChan <- struct{}{}
}

func (s *Server) OnConnect() {

}
