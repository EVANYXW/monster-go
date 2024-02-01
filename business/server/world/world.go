package world

import (
	"bilibili/monster-go/network"
	"fmt"
	"os"
	"syscall"
)

type MessageId int32

type world struct {
	server   *network.Server
	handlers map[MessageId]func(message *network.Packet)
}

var Oasis *world

func NewWorld() *world {
	return &world{
		server:   network.NewServer(":8023", 100, 200),
		handlers: make(map[MessageId]func(message *network.Packet)),
	}
}

func (w *world) Start() {
	w.HandlerRegister()

	go w.server.Run()
}

func (w *world) Stop() {

}

func (w *world) CreatePlayer(message *network.Packet) {

}

func (w *world) UserLogin(message *network.Packet) {

}

func (w *world) OnSystemSignal(signal os.Signal) bool {
	//logger.Logger.DebugF("[World] 收到信号 %v \n", signal)
	fmt.Printf("[World] 收到信号 %v \n", signal.String())
	tag := true
	switch signal {
	case syscall.SIGHUP:
		//todo
		fmt.Println(11)
	case syscall.SIGPIPE:
		fmt.Println(22)
	default:
		//logger.Logger.DebugF("[World] 收到信号准备退出...")
		tag = false

	}
	return tag
}
