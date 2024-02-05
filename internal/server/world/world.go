package world

import (
	"bilibili/monster-go/configs"
	"bilibili/monster-go/internal/network"
	rpcServer "bilibili/monster-go/internal/rpc/server"
	"bilibili/monster-go/pkg/env"
	"bilibili/monster-go/pkg/logger"
	"bilibili/monster-go/pkg/timeutil"
	"fmt"
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"go.uber.org/zap"

	//"github.com/phuhao00/network/example/logger"
	"os"
	"syscall"
)

var (
	Logger *zap.Logger
)

type handlerFunc func(message *network.Packet)

type world struct {
	server    *network.Server
	handlers  map[messageId.MessageId]handlerFunc
	closeChan chan struct{}
}

var Oasis *world

func Init() {
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

func NewWorld() *world {
	Init()

	config := configs.Get().Server
	w := &world{
		handlers:  make(map[messageId.MessageId]handlerFunc),
		closeChan: make(chan struct{}),
	}

	s := network.NewServer(fmt.Sprintf("%s", config.Address),
		config.MaxConnNum, config.BuffSize, Logger)

	s.MessageHandler = w.OnSessionPacket
	w.server = s

	// 监听配置文件

	go func() {
	outer:
		for {
			select {
			case <-w.closeChan:
				break outer
			case s := <-configs.NotifyChan:
				// TODO: 监听 configs的本地配置文件,有修改重新加载
				fmt.Println(s)
			}
		}
	}()

	return w
}

func (w *world) Start() {
	w.HandlerRegister()
	go w.server.Run()

	worldServer := &rpcServer.WorldServer{}
	go worldServer.Run()
}

func (w *world) Stop() {
	// close
	Logger.Sync()
	go func() {
		w.closeChan <- struct{}{}
		w.server.OnClose()
	}()

}

// OnSessionPacket 根据注册方法调佣
func (w *world) OnSessionPacket(packet *network.Packet) {
	if handler, ok := w.handlers[messageId.MessageId(packet.Msg.ID)]; ok {
		handler(packet)
		return
	}

	//
	//if packet.Msg.ID == 3 {
	//	fmt.Println("hello hao")
	//}
	//packet.Conn.AsyncSend(
	//	1,
	//	//&player.SCSendChatMsg{},
	//	&player.SCLogin{
	//		Ok: true,
	//	},
	//)
}

// OnSystemSignal 监听退出信道
func (w *world) OnSystemSignal(signal os.Signal) bool {
	//logger.Logger.DebugF("[World] 收到信号 %v \n", signal)
	fmt.Printf("[World] 收到信号 %v \n", signal.String())
	tag := true
	switch signal {
	case syscall.SIGHUP:
		//todo
		fmt.Println("SIGHUP")
	case syscall.SIGPIPE:
		fmt.Println("SIGPIPE")
	default:
		//logger.Logger.DebugF("[World] 收到信号准备退出...")
		tag = false
	}
	return tag
}
