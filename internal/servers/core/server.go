package core

import (
	"github.com/evanyxw/monster-go/configs"
	"github.com/evanyxw/monster-go/internal/configure"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/server"
)

type Server struct {
	//modules     map[int]module.IModule
	netAcceptor *network.Acceptor
	handlers    network.HandlerMap
	closeChan   chan struct{}
}

func NewServer(info server.Info) *Server {
	config := configs.Get().Center
	w := &Server{
		handlers:    make(network.HandlerMap, network.Pool_id_Max),
		closeChan:   make(chan struct{}),
		netAcceptor: network.NewAcceptor(config.MaxConnNum, info, nil, nil),
	}
	w.netAcceptor.MessageHandler = w.MessageHandler

	//w.modules[module.ModuleID_Net] = module.NewNetKernel()

	return w
}

func (s *Server) Init() {

}

func (s *Server) Gateg(msgId uint16, fun network.HandlerFunc) {
	s.handlers[msgId] = fun
}

func (w *Server) Run() {
	// 加载配置
	configure.Global.Load()

	go w.netAcceptor.Run()

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
		w.netAcceptor.OnClose()
	}()
}

// MessageHandler 根据注册方法调用
func (w *Server) MessageHandler(packet *network.Packet) {
	//if handler, ok := w.handlers[messageId.MessageId(packet.Msg.ID)]; ok {
	//	handler(packet)
	//	return
	//}
	handler := w.handlers[packet.Msg.ID]
	handler(packet)
}
