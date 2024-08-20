package network

import (
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/evanyxw/monster-go/pkg/logger"
	"go.uber.org/zap"
)

type Processor struct {
	handlers HandlerMap
}

var (
	GlobalProcess *Processor
)

func GetProcessor() *Processor {
	return GlobalProcess
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

	if handler == nil {
		logger.Error("handler is not found!", zap.Uint64("message_id", packet.Msg.ID))
		return
	}
	handler(packet)
}
