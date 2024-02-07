package world

import (
	"bilibili/monster-go/configs"
	"bilibili/monster-go/internal/configure"
	"bilibili/monster-go/internal/network"
	rpcServer "bilibili/monster-go/internal/rpc/server"
	"bilibili/monster-go/internal/server/factory"
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

type World struct {
	networkServer *network.Server
	handlers      map[messageId.MessageId]handlerFunc
	closeChan     chan struct{}
}

var Oasis *World

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

func New(info network.Info) factory.Server {
	// 日志初始化
	initLog()

	config := configs.Get().Server
	w := &World{
		handlers:  make(map[messageId.MessageId]handlerFunc),
		closeChan: make(chan struct{}),
	}

	w.networkServer = network.NewServer(fmt.Sprintf("%s", config.Address),
		config.MaxConnNum, config.BuffSize, Logger, info)

	w.networkServer.MessageHandler = w.OnSessionPacket

	return w
}

func (w *World) Start() {

	// 加载配置
	configure.Global.Load()

	// pb消息的注册
	w.HandlerRegister()
	go w.networkServer.Run()

	worldRpcServer := &rpcServer.WorldServer{}
	go worldRpcServer.Run()

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

func (w *World) Stop() {
	Logger.Sync()
	go func() {
		w.closeChan <- struct{}{}
		w.networkServer.OnClose()
	}()

}

// OnSessionPacket 根据注册方法调佣
func (w *World) OnSessionPacket(packet *network.Packet) {
	if handler, ok := w.handlers[messageId.MessageId(packet.Msg.ID)]; ok {
		handler(packet)
		return
	}
}

// OnSystemSignal 监听退出信道
func (w *World) OnSystemSignal(signal os.Signal) bool {
	tag := true
	switch signal {
	case syscall.SIGHUP:
		//todo
		fmt.Println("SIGHUP")
	case syscall.SIGPIPE:
		fmt.Println("SIGPIPE")
	default:
		Logger.Debug("[World] 收到信号准备退出 %v \n", zap.String("signal", signal.String()))
		tag = false
	}
	return tag
}
