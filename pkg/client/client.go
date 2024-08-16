package client

import (
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"go.uber.org/zap"
	"net"
	"sync/atomic"
)

type Client struct {
	*network.NetPoint
	address         string
	running         atomic.Value
	OnMessageCb     func(message *network.Packet)
	OnCloseCallBack func()

	//msgParser   *BufferPacker
	msgParser   network.Packer
	rpcAcceptor *rpc.Acceptor
	processor   *network.Processor

	//closed          int32
	//ChMsg   chan *Message
	//logger      *zap.Logger
	//bufferSize      int
	//logger          *spoor.Spoor
}

func NewClient(address string, rpcAcceptor *rpc.Acceptor, processor *network.Processor, packer network.Packer) *Client {
	client := &Client{
		//bufferSize: connBuffSize,
		address: address,
		//msgParser:   network.NewDefaultPacker(),
		msgParser:   packer,
		rpcAcceptor: rpcAcceptor,
		processor:   processor,
	}

	client.running.Store(false)
	return client
}

func (c *Client) Dial() (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", c.address)
	//logger.Info("ConnectorKernel 创建新链接", zap.String("address:", c.address))
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (c *Client) Run() {
	conn, err := c.Dial()
	if err != nil {
		//c.logger.ErrorF("%v", err)
		logger.Error("Client Run is error:", zap.Error(err))
		return
	}

	tcpConn, err := network.NewNetPoint(conn, c.msgParser)

	if err != nil {
		logger.Error(err.Error())
		return
	}
	c.NetPoint = tcpConn

	c.SetNetEventRPC(c.rpcAcceptor)
	c.SetProcessor(c.processor)
	c.RpcAcceptor.Run(c.NetPoint.CloseChan)

	// fixMe Go到哪里去了
	c.RpcAcceptor.Go(rpc.RPC_NET_CONNECTED, tcpConn)

	c.Impl = c
	c.Reset()
	c.running.Store(true)
	go c.Connect()
}

func (c *Client) OnMessage(data *network.Message, conn *network.NetPoint) {

	//c.Verify()

	defer func() {
		if err := recover(); err != nil {
			//c.logger.ErrorF("[OnMessage] panic ", err, "\n", string(debug.Stack()))
		}
	}()

	if c.OnMessageCb == nil {
		logger.Error("[OnMessage] is nil")
		return
	}

	c.OnMessageCb(&network.Packet{
		Msg:      data,
		NetPoint: conn,
	})
}
func (c *Client) OnClose() {
	if c.OnCloseCallBack != nil {
		c.OnCloseCallBack()
	}
	c.running.Store(false)
	c.NetPoint.Close()
}

func (c *Client) Pack(msgID uint64, msg interface{}) (pack []byte, err error) {
	pack, err = c.msgParser.Pack(msgID, msg)
	return
}

func (c *Client) UnPack(data []byte) (pack *network.Message, err error) {
	pack, err = c.msgParser.Unpack(data)
	return
}
