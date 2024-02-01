package world

func (w *world) HandlerRegister() {
	w.handlers[1] = w.CreatePlayer
	w.handlers[2] = w.UserLogin
}
