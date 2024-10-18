// Package server @Author evan_yxw
// @Date 2024/10/16 19:13:00
// @Desc
package server

type Info struct {
	ServerName string
	Ip         string
	Port       uint32
	Env        string
	Address    string
	RpcAddr    string
}

var (
	info *Info
)

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
