package client

import (
	"github.com/evanyxw/monster-go/pkg/client"
	"github.com/evanyxw/monster-go/pkg/network"
	"os"
	"syscall"
	"time"
)

type Client struct {
	cli             *client.Client
	inputHandlers   map[string]InputHandler
	messageHandlers network.HandlerMap
	console         *ClientConsole
	chInput         chan *InputParam
}

func NewClient() *Client {
	//rpcAcceptor := rpc.NewAcceptor(1000)
	processor := network.NewProcessor()
	c := &Client{
		//cli:             client.NewClient(":20001", rpcAcceptor, processor),
		cli:             client.NewClient(":30000", processor, new(network.DefaultPackerFactory)),
		inputHandlers:   map[string]InputHandler{},
		messageHandlers: make(network.HandlerMap, network.Pool_id_Max),
		console:         NewClientConsole(),
	}
	c.cli.OnMessageCb = c.OnMessage
	//c.cli.ChMsg = make(chan *network2.Message, 1)
	//c.chInput = make(chan *InputParam, 1)
	//c.console.chInput = c.chInput
	c.console.chInput = make(chan *InputParam, 1)
	c.chInput = c.console.chInput
	//https://github.com/phuhao00/greatestworks-proto.git
	//github.com/phuhao00/greatestworks-proto
	return c
}

func (c *Client) Run() {
	go func() {
		for {
			select {
			case input := <-c.chInput:
				inputHandler := c.inputHandlers[input.Command]
				if inputHandler != nil {
					inputHandler(input)
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	go c.console.Run()
	go c.cli.Run()
}

func (c *Client) OnMessage(packet *network.Packet) {
	//fmt.Println(c.messageHandlers)
	//fmt.Println(packet.Msg.ID)
	//if handler, ok := c.messageHandlers[messageId.MessageId(packet.Msg.ID)]; ok {
	//	handler(packet)
	//}
}

func (c *Client) OnSystemSignal(signal os.Signal) bool {
	//logger.Logger.InfoF("[Client] 收到信号 %v \n", signal)
	tag := true
	switch signal {
	case syscall.SIGHUP:
		//todo
	case syscall.SIGPIPE:
	default:
		//logger.Logger.InfoF("[Client] 收到信号准备退出...")
		tag = false

	}
	return tag
}
