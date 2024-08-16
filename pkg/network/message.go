package network

type Message struct {
	ID    uint64
	Data  []byte
	RawID uint32
}

type Packet struct {
	Msg      *Message
	NetPoint *NetPoint
}
