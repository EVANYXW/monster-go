package network

type Message struct {
	ID    uint64
	Data  []byte
	RawID uint64
}

type Packet struct {
	Msg      *Message
	NetPoint *NetPoint
}
