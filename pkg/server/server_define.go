package server

type Info struct {
	ServerName string
	Ip         string
	Port       uint32
	Env        string
	Address    string
	RpcAddr    string
}

// 为服务网络节点的状态，为了表示已经准备好可以其他服务器来链接
const (
	Net_RunStep_None = iota
	Net_RunStep_Start
	Net_RunStep_Done
)

// center 启动其他服务器的状态
const (
	CN_RunStep_None = iota
	CN_RunStep_StartServer
	CN_RunStep_WaitHandshake
	CN_RunStep_HandshakeDone
	CN_RunStep_Done
)
