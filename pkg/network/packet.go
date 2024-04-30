package network

type Packet struct {
	Msg  *Message
	Conn *TcpConn
}
