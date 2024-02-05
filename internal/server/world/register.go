package world

import "github.com/phuhao00/greatestworks-proto/gen/messageId"

func (w *world) HandlerRegister() {
	w.handlers[messageId.MessageId_CSCreatePlayer] = w.CreatePlayer
	w.handlers[messageId.MessageId_CSLogin] = w.UserLogin
}
