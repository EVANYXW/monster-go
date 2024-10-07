// Package ipPort
// @Description: ip port package
// @Author evan_yxw 2024-06-14 19:53:49

package ipPort

import (
	"errors"
	"fmt"
	"net"
)

// GetDynamicIpAndPort get dynamic ip and port
func GetDynamicIpAndPort() (string, int, error) {
	ip := "0.0.0.0"
	port := 80
	IP, err := ExternalIP()
	if err != nil {
		return ip, port, err
	}
	ip = IP.String()
	port, err = GetFreePort()
	return ip, port, err
}

// GetDynamicIpAndRangePort get dynamic ip and range port
func GetDynamicIpAndRangePort(minPort, maxPort int) (string, int, error) {
	ip := "0.0.0.0"
	port := 80
	IP, err := ExternalIP()
	if err != nil {
		return ip, port, err
	}
	ip = IP.String()
	port, err = GetFreeRangePort(ip, minPort, maxPort)
	return ip, port, err
}

// GetFreeRangePort get free range port
func GetFreeRangePort(ip string, minPort, maxPort int) (port int, err error) {
	// 检查端口范围是否有效
	if minPort <= 0 || maxPort > 65535 || minPort > maxPort {
		return 0, errors.New("invalid port range")
	}
	for port = minPort; port <= maxPort; port++ {
		addr := fmt.Sprintf("%s:%d", ip, port)
		listener, err := net.Listen("tcp", addr)
		if err == nil {
			defer listener.Close()
			return listener.Addr().(*net.TCPAddr).Port, nil
		}
	}

	return 0, errors.New("no free ports available in the specified range")
}

// GetFreePort get free port
func GetFreePort() (port int, err error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, nil
	}
	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, nil
	}
	defer listen.Close()
	return listen.Addr().(*net.TCPAddr).Port, nil
}

// ExternalIP get external ip
func ExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}
	return ip
}
