package server

func Name2EP(name string) uint8 {
	if val, ok := NameMap[name]; ok {
		return val
	}
	return EP_None
}

func EP2Name(name uint8) string {
	if val, ok := EPMap[name]; ok {
		return val
	}
	return ""
}

// 网络端点定义
const (
	EP_None = iota
	EP_Client
	EP_Center
	EP_Login
	EP_Gate
	EP_Game
	EP_Manager
	EP_World
	EP_Mail
	EP_Robot

	EP_Max
)

var NameMap = map[string]uint8{
	"center":  EP_Center,
	"gate":    EP_Gate,
	"login":   EP_Login,
	"game":    EP_Game,
	"mail":    EP_Mail,
	"manager": EP_Manager,
	"world":   EP_World,
}

var EPMap = map[uint8]string{
	EP_Center:  "center",
	EP_Gate:    "gate",
	EP_Login:   "login",
	EP_Game:    "game",
	EP_Mail:    "mail",
	EP_Manager: "manager",
	EP_World:   "world",
}
