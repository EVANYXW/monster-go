package core

import (
	"github.com/evanyxw/monster-go/pkg/network"
)

type Client struct {
	messageHandlers network.HandlerMap
	*network.Client
}

func NewClient() *Client {
	client := &Client{
		messageHandlers: make(network.HandlerMap, network.Pool_id_Max),
		Client:          network.NewClient("", nil),
	}
	return client
}

func (c *Client) OnMessageCb(packet *network.Packet) {
	//fmt.Println(c.messageHandlers)
	//fmt.Println(packet.Msg.ID)
	//if handler, ok := c.messageHandlers[messageId.MessageId(packet.Msg.ID)]; ok {
	//	handler(packet)
	//}
	handler := c.messageHandlers[packet.Msg.ID]
	handler(packet)
}
