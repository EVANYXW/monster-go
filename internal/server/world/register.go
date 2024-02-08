package world

import "github.com/evanyxw/game_proto/msg/messageId"

func (w *World) HandlerRegister() {
	w.handlers[messageId.MessageId_CSCreatePlayer] = w.CreatePlayer
	w.handlers[messageId.MessageId_CSLogin] = w.UserLogin
}
