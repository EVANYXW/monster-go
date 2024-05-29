package network

import "github.com/evanyxw/monster-go/message/pb/xsf_pb"

type Processor struct {
	handlers HandlerMap
}

func NewProcessor() *Processor {
	return &Processor{
		make(HandlerMap, xsf_pb.SMSGID_Server_Max),
	}
}

func (c *Processor) RegisterMsg(msgId uint16, handlerFunc HandlerFunc) {
	c.handlers[msgId] = handlerFunc
}

func (c *Processor) MessageHandler(packet *Packet) {
	handler := c.handlers[packet.Msg.ID]
	handler(packet)
}
