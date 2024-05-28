package network

import (
	"encoding/binary"
	"errors"
	"github.com/golang/protobuf/proto"
	"io"
)

type ClientBufferPacker struct {
	lenMsgLen int32
	minMsgLen uint32
	maxMsgLen uint32
	recvBuff  *ByteBuffer
	sendBuff  *ByteBuffer
	byteOrder binary.ByteOrder
}

func NewClientPacker() *ClientBufferPacker {
	msgParser := &ClientBufferPacker{
		lenMsgLen: 4,
		minMsgLen: 2,
		maxMsgLen: 2 * 1024 * 1024,
		recvBuff:  NewByteBuffer(),
		sendBuff:  NewByteBuffer(),
		byteOrder: binary.LittleEndian,
	}
	return msgParser
}

func (p *ClientBufferPacker) GetMaxMsgLen() uint32 {
	return p.maxMsgLen
}

func (p *ClientBufferPacker) Reset() {
	p.recvBuff = NewByteBuffer()
	p.sendBuff = NewByteBuffer()
}

func (p *ClientBufferPacker) Write(tcpCoon *NetPoint, buff ...byte) error {
	// get len
	msgLen := uint32(len(buff))

	// check len
	if msgLen > p.maxMsgLen {
		return errors.New("message too long")
	} else if msgLen < p.minMsgLen {
		return errors.New("message too short")
	}

	_, err := tcpCoon.Write(buff)
	return err
}

// conn *NetPoint
func (p *ClientBufferPacker) Read(conn io.Reader) ([]byte, error) {
	var b [MsgLengthSize]byte

	_, err := io.ReadFull(conn, b[:])

	//buffer := bytes.NewBuffer(c)
	//_, err := io.ReadFull(buffer, b[:])

	// 包长度
	packLen := binary.LittleEndian.Uint32(b[:])
	//xsf_log.Debug("defaultPakcer Read", xsf_log.Int("packLen", int(packLen)))
	MaxMsgDataLen := 0xFFFFFFFF
	if packLen > uint32(MaxMsgDataLen) {
		return []byte(""), err
		//return "", nil, 0, 0, fmt.Errorf("message too long, len=%d", packLen)
	} else if packLen < MinMsgDataLen {
		return []byte(""), err
		//return "", nil, 0, 0, fmt.Errorf("message too short, len=%d", packLen)
	}

	totalMsg := make([]byte, packLen)
	//_, err = io.ReadFull(buffer, totalMsg)
	_, err = io.ReadFull(conn, totalMsg)
	// 消息ID
	//msgID := binary.LittleEndian.Uint16(totalMsg[:2])
	//xsf_log.Debug("defaultPakcer Read", xsf_log.Int("bufMsgID", int(msgID)))

	//rawID := binary.LittleEndian.Uint32(totalMsg[2:6])
	//xsf_log.Debug("defaultPakcer Read", xsf_log.Uint32("rawID", rawID))

	msgData := make([]byte, packLen-6)
	copy(msgData, totalMsg[6:])

	return msgData, err
	//return "", msgData, msgID, rawID, err
	//return &Message{ID: uint64(msgID), Data: msgData}, err
}

func (p *ClientBufferPacker) Pack(msgId uint64, msg interface{}) ([]byte, error) {
	pbMsg, ok := msg.(proto.Message)
	if !ok {
		//return []byte{}, fmt.Errorf("msg is not protobuf message")
	}
	// data
	data, err := proto.Marshal(pbMsg)
	if err != nil {
		return data, err
	}
	//data := msg.Data

	msgLen := uint32(len(data))

	extLen := MinMsgDataLen + MsgLengthSize
	// 总长度 = 4字节长度 + 2字节消息id + 4字节rawID + pb消息长度
	// 所以消息最小长度为6
	total := msgLen + uint32(extLen)
	packLen := msgLen + MinMsgDataLen

	//xsf_log.Debug("defaultPakcer Pack", xsf_log.Uint16("id", msg.GetID()), xsf_log.String("message", msg.ToString()))
	//xsf_log.Debug("defaultPakcer Pack", xsf_log.Uint32("msgLen", msgLen), xsf_log.Uint32("packLen", packLen), xsf_log.Uint32("total", total))

	buffer := make([]byte, total)

	// 先写入4个字节的包长度
	binary.LittleEndian.PutUint32(buffer, packLen)

	// 再写入消息id
	binary.LittleEndian.PutUint16(buffer[4:], uint16(msgId))

	// 中间4个字节rawID，默认写入消息ID
	binary.LittleEndian.PutUint32(buffer[6:], uint32(msgId))

	// 写入PB数据
	copy(buffer[10:], data)

	return buffer, err
}

// Unpack id + protobuf data
func (p *ClientBufferPacker) Unpack(data []byte) (*Message, error) {
	if len(data) < 2 {
		return nil, errors.New("protobuf data too short")
	}
	msgID := p.byteOrder.Uint16(data[:2])
	msg := &Message{
		ID:   uint64(msgID),
		Data: data[2:],
	}
	return msg, nil
}
