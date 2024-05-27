package server

type Info struct {
	ServerName string
	Env        string
	Address    string
	RpcAddr    string
}

const (
	CN_RunStep_None = iota
	CN_RunStep_StartServer
	CN_RunStep_WaitHandshake
	CN_RunStep_HandshakeDone
	CN_RunStep_Done
)
