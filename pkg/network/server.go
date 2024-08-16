package network

import (
	"github.com/evanyxw/monster-go/pkg/server"
	"net/http"
)

type HandlerFunc func(message *Packet)
type HandlerMap []HandlerFunc

const (
	Pool_id_Max = 100
)

const (
	ServerInfo_None = 0
	ServerInfo_New  = 1
	ServerInfo_Ok   = 2
)

type ServerInfo struct {
	ID     uint32
	IP     string
	Ports  [server.EP_Max]uint32
	Status uint32
}

//const (
//	ServerID_Center = iota
//	ServerID_Gate
//	ServerID_World
//)

//const (
//	ServerStatus_None     = 0
//	ServerStatus_Starting = 1
//	ServerStatus_Running  = 2
//	ServerStatus_Stopping = 3
//)

func OpenPPROF() {
	go func() {
		http.ListenAndServe("localhost:8000", nil)
	}()
}
