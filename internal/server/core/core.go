package core

import (
	"fmt"
	"github.com/evanyxw/game_proto/msg/messageId"
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/internal/configure"
	"github.com/evanyxw/monster-go/internal/network"
	"github.com/evanyxw/monster-go/pkg/env"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/timeutil"
	"go.uber.org/zap"
)

var (
	Logger *zap.Logger
)

type handlerFunc func(message *network.Packet)

type Server struct {
	networkServer *network.Server
	handlers      map[messageId.MessageId]handlerFunc
	closeChan     chan struct{}
}

var Oasis *Server

func initLog() {
	log, err := logger.NewJSONLogger(
		logger.WithDisableConsole(),
		logger.WithField("domain", fmt.Sprintf("%s[%s]", configs.ProjectName, env.Active().Value())),
		logger.WithTimeLayout(timeutil.CSTLayout),
		logger.WithFileP(configs.ProjectCronLogFile),
	)
	if err != nil {
		panic(err)
	}
	Logger = log
}

func NewServer(info network.Info, logger *zap.Logger) *Server {
	// 日志初始化
	//initLog()

	config := configs.Get().Server
	w := &Server{
		handlers:  make(map[messageId.MessageId]handlerFunc),
		closeChan: make(chan struct{}),
		networkServer: network.NewServer(fmt.Sprintf("%s", config.Address),
			config.MaxConnNum, config.BuffSize, logger, info),
	}
	w.networkServer.MessageHandler = w.OnSessionPacket

	return w
}

func (s *Server) RegisterMsg(msgId messageId.MessageId, fun handlerFunc) {
	s.handlers[msgId] = fun
}

func (w *Server) Run() {
	// 加载配置
	configure.Global.Load()

	//// pb消息的注册
	//w.HandlerRegister() // todo
	go w.networkServer.Run()

	//worldRpcServer := &rpcServer.WorldServer{}
	//go worldRpcServer.Run()

	// 监听配置文件
	go func() {
	outer:
		for {
			select {
			case <-w.closeChan:
				break outer
			case <-configs.NotifyChan:
				// TODO: 监听 configs的本地配置文件,有修改重新加载
			}
		}
	}()
}

func (w *Server) Destroy() {
	//Logger.Sync()
	go func() {
		w.closeChan <- struct{}{}
		w.networkServer.OnClose()
	}()

}

// OnSessionPacket 根据注册方法调佣
func (w *Server) OnSessionPacket(packet *network.Packet) {
	if handler, ok := w.handlers[messageId.MessageId(packet.Msg.ID)]; ok {
		handler(packet)
		return
	}
}

// OnSystemSignal 监听退出信道
//func (w *Server) OnSystemSignal(signal os.Signal) bool {
//	tag := true
//	switch signal {
//	case syscall.SIGHUP:
//		//todo
//		fmt.Println("SIGHUP")
//	case syscall.SIGPIPE:
//		fmt.Println("SIGPIPE")
//	default:
//		Logger.Debug("[World] 收到信号准备退出 %v \n", zap.String("signal", signal.String()))
//		tag = false
//	}
//	return tag
//}
