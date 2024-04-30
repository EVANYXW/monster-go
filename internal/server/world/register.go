package world

import (
	"github.com/evanyxw/game_proto/msg/messageId"
)

func (w *World) HandlerRegister() {
	w.RegisterMsg(messageId.MessageId_CSCreatePlayer, w.CreatePlayer)
	w.RegisterMsg(messageId.MessageId_CSLogin, w.UserLogin)
}
