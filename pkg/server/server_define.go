package server

import "google.golang.org/grpc"

type GrpcServer interface {
	TransportRegister() func(grpc.ServiceRegistrar) error
}

type Info struct {
	ServerName string
	Ip         string
	Port       uint32
	Env        string
	Address    string
	RpcAddr    string
}

// center 启动其他服务器的状态
const (
	CN_RunStep_None = iota
	CN_RunStep_StartServer
	CN_RunStep_WaitHandshake
	CN_RunStep_HandshakeDone
	CN_RunStep_Done
)

// 为服务网络节点的状态，为了表示已经准备好可以其他服务器来链接
const (
	Net_RunStep_Start = iota + 1
	Net_RunStep_Done
)

var (
	info *Info
)

func StatusStart(status *int) {
	*status = Net_RunStep_Start
}

func StatusDone(status *int) {
	*status = Net_RunStep_Done
}

func StatusIsDone(status int) bool {
	return status == Net_RunStep_Done
}

func SetServerInfo(i *Info) {
	info = i
}

func SetInfoIP(ip string) {
	info.Ip = ip
}

func SetInfoPort(port uint32) {
	info.Port = port
}

func GetServerInfo() *Info {
	return info
}
