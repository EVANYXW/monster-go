package network

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
)

type PackerFactory interface {
	CreatePacker() Packer
}

type DefaultPackerFactory struct {
}

func (f *DefaultPackerFactory) CreatePacker() Packer {
	return NewDefaultPacker()
}

type ClientPackerFactory struct {
}

func (f *ClientPackerFactory) CreatePacker() Packer {
	return NewClientPacker()
}

type Packer interface {
	Write(tcpCoon *NetPoint, buff ...byte) error
	Read(conn io.Reader) ([]byte, error)
	Pack(msgID uint64, msg interface{}) ([]byte, error)
	//PackData(msgID uint64, data []byte) ([]byte, error)
	Unpack(data []byte) (*Message, error)
	Reset()
	GetMaxMsgLen() uint32
	GetType() string
}

//var once sync.Once
//var defaultBufferPacker *DefaultBufferPacker

type DefaultBufferPacker struct {
	lenMsgLen int32
	minMsgLen uint32
	maxMsgLen uint32
	recvBuff  *ByteBuffer
	sendBuff  *ByteBuffer
	byteOrder binary.ByteOrder
}

func NewDefaultPacker() *DefaultBufferPacker {
	//once.Do(func() {
	//	defaultBufferPacker = &DefaultBufferPacker{
	//		lenMsgLen: 4,
	//		//lenMsgLen: 2,
	//		minMsgLen: 2,
	//		maxMsgLen: 2 * 1024 * 1024,
	//		recvBuff:  NewByteBuffer(),
	//		sendBuff:  NewByteBuffer(),
	//		byteOrder: binary.LittleEndian,
	//	}
	//})
	return &DefaultBufferPacker{
		lenMsgLen: 4,
		//lenMsgLen: 2,
		minMsgLen: 2,
		maxMsgLen: 2 * 1024 * 1024,
		recvBuff:  NewByteBuffer(),
		sendBuff:  NewByteBuffer(),
		byteOrder: binary.LittleEndian,
	}
	//return defaultBufferPacker

}

func (p *DefaultBufferPacker) GetType() string {
	return "default"
}

func (p *DefaultBufferPacker) GetMaxMsgLen() uint32 {
	return p.maxMsgLen
}

func (p *DefaultBufferPacker) Write(tcpCoon *NetPoint, buff ...byte) error {
	// get len
	msgLen := uint32(len(buff))

	// check len
	if msgLen > p.maxMsgLen {
		return errors.New("message too long")
	} else if msgLen < p.minMsgLen {
		return errors.New("message too short")
	}

	// write len
	switch p.lenMsgLen {
	case 2:
		p.sendBuff.AppendInt16(int16(msgLen))
	case 4:
		p.sendBuff.AppendInt32(int32(msgLen))
	}

	p.sendBuff.Append(buff)
	// write data
	writeBuff := p.sendBuff.ReadBuff()[:p.sendBuff.Length()]

	_, err := tcpCoon.Write(writeBuff)

	p.sendBuff.Reset()

	return err
}

func (p *DefaultBufferPacker) Read(conn io.Reader) ([]byte, error) {
	p.recvBuff.EnsureWritableBytes(p.lenMsgLen)

	readLen, err := io.ReadFull(conn, p.recvBuff.WriteBuff()[:p.lenMsgLen])
	//fmt.Println("read len:", readLen)
	// read len
	if err != nil {
		return nil, fmt.Errorf("%v readLen:%v", err, readLen)
	}
	p.recvBuff.WriteBytes(int32(readLen))

	// parse len
	var msgLen uint32
	switch p.lenMsgLen {
	case 2:
		msgLen = uint32(p.recvBuff.ReadInt16())
	case 4:
		msgLen = uint32(p.recvBuff.ReadInt32())
	}

	// check len
	if msgLen > p.maxMsgLen {
		return nil, errors.New("message too long")
	} else if msgLen < p.minMsgLen {
		return nil, errors.New("message too short")
	}

	p.recvBuff.EnsureWritableBytes(int32(msgLen))

	rLen, err := io.ReadFull(conn, p.recvBuff.WriteBuff()[:msgLen])
	if err != nil {
		return nil, fmt.Errorf("%v msgLen:%v readLen:%v", err, msgLen, rLen)
	}
	p.recvBuff.WriteBytes(int32(rLen))

	/*
		// 保留了2字节flag 暂时未处理
		var flag uint16
		flag = uint16(p.recvBuff.ReadInt16())
	*/
	p.recvBuff.Skip(2) // 跳过2字节保留字段

	// 减去2字节的保留字段长度
	return p.recvBuff.NextBytes(int32(msgLen - 2)), nil
}

func (p *DefaultBufferPacker) Reset() {
	p.recvBuff = NewByteBuffer()
	p.sendBuff = NewByteBuffer()
}

func (p *DefaultBufferPacker) Pack(msgID uint64, msg interface{}) ([]byte, error) {
	//pack := p.TestPack(msgID, msg)
	//return pack, nil

	pbMsg, ok := msg.(proto.Message)
	if !ok {
		return []byte{}, fmt.Errorf("msg is not protobuf message")
	}
	// data
	data, err := proto.Marshal(pbMsg)
	if err != nil {
		return data, err
	}
	// 4byte = len(flag)[2byte] + len(msgID)[2byte]
	buf := make([]byte, 4+len(data))
	mId := uint16(msgID)
	if p.byteOrder == binary.LittleEndian {
		binary.LittleEndian.PutUint16(buf[0:2], 0)
		binary.LittleEndian.PutUint16(buf[2:], mId)
	} else {
		binary.BigEndian.PutUint16(buf[0:2], 0)
		binary.BigEndian.PutUint16(buf[2:], mId)
	}
	copy(buf[4:], data)
	return buf, err
}

//func (p *DefaultBufferPacker) PackData(msgID uint64, data []byte) ([]byte, error) {
//	//pack := p.TestPack(msgID, msg)
//	//return pack, nil
//
//	pbMsg, ok := msg.(proto.Message)
//	if !ok {
//		return []byte{}, fmt.Errorf("msg is not protobuf message")
//	}
//	// data
//	data, err := proto.Marshal(pbMsg)
//	if err != nil {
//		return data, err
//	}
//	// 4byte = len(flag)[2byte] + len(msgID)[2byte]
//	buf := make([]byte, 4+len(data))
//	mId := uint16(msgID)
//	if p.byteOrder == binary.LittleEndian {
//		binary.LittleEndian.PutUint16(buf[0:2], 0)
//		binary.LittleEndian.PutUint16(buf[2:], mId)
//	} else {
//		binary.BigEndian.PutUint16(buf[0:2], 0)
//		binary.BigEndian.PutUint16(buf[2:], mId)
//	}
//	copy(buf[4:], data)
//	return buf, err
//}

// Unpack id + protobuf data
func (p *DefaultBufferPacker) Unpack(data []byte) (*Message, error) {
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

const MsgLengthSize = 4
const MinMsgDataLen = 6

func (p *DefaultBufferPacker) TestPack(msgId uint64, msg interface{}) []byte {
	pbMsg, ok := msg.(proto.Message)
	if !ok {
		//return []byte{}, fmt.Errorf("msg is not protobuf message")
	}
	// data
	data, _ := proto.Marshal(pbMsg)

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

	return buffer
}

// conn *NetPoint
func (p *DefaultBufferPacker) TestRead(c []byte) (*Message, error) {
	var b [MsgLengthSize]byte

	//_, err := io.ReadFull(conn, b[:])

	buffer := bytes.NewBuffer(c)
	_, err := io.ReadFull(buffer, b[:])

	// 包长度
	packLen := binary.LittleEndian.Uint32(b[:])
	//xsf_log.Debug("defaultPakcer Read", xsf_log.Int("packLen", int(packLen)))
	MaxMsgDataLen := 0xFFFFFFFF
	if packLen > uint32(MaxMsgDataLen) {
		return &Message{}, err
		//return "", nil, 0, 0, fmt.Errorf("message too long, len=%d", packLen)
	} else if packLen < MinMsgDataLen {
		return &Message{}, err
		//return "", nil, 0, 0, fmt.Errorf("message too short, len=%d", packLen)
	}

	totalMsg := make([]byte, packLen)
	_, err = io.ReadFull(buffer, totalMsg)
	//_, err = io.ReadFull(conn, totalMsg)
	// 消息ID
	msgID := binary.LittleEndian.Uint16(totalMsg[:2])
	//xsf_log.Debug("defaultPakcer Read", xsf_log.Int("bufMsgID", int(msgID)))

	//rawID := binary.LittleEndian.Uint32(totalMsg[2:6])
	//xsf_log.Debug("defaultPakcer Read", xsf_log.Uint32("rawID", rawID))

	msgData := make([]byte, packLen-6)
	copy(msgData, totalMsg[6:])

	//return "", msgData, msgID, rawID, err
	return &Message{ID: uint64(msgID), Data: msgData}, err
}

//func (p *BufferPacker) PackData(messageID uint32, data interface{}) ([]byte, error) {
//	pbMsg, ok := data.(proto.Message)
//	if !ok {
//		return []byte{}, fmt.Errorf("msg is not protobuf message")
//	}
//	// 序列化proto对象
//	dataBytes, err := proto.Marshal(pbMsg)
//	if err != nil {
//		return nil, err
//	}
//
//	// 创建一个缓冲区
//	buffer := make([]byte, 4+len(dataBytes))
//
//	// 将messageID写入缓冲区
//	binary.LittleEndian.PutUint16(buffer[:2], 0)
//	binary.LittleEndian.PutUint16(buffer[2:4], uint16(messageID))
//
//	// 将data写入缓冲区
//	copy(buffer[4:], dataBytes)
//
//	return buffer, nil
//}
//
//func UnpackData(conn net.Conn) ([]byte, error) {
//	// 读取前4个字节，判断messageID
//	messageIDBytes := make([]byte, 4)
//	_, err := io.ReadFull(conn, messageIDBytes)
//	if err != nil {
//		return nil, err
//	}
//	messageID := binary.LittleEndian.Uint16(messageIDBytes)
//	fmt.Println(messageID)
//	// 读取剩余的数据作为data
//	dataBytes, err := io.ReadAll(conn)
//	if err != nil {
//		return nil, err
//	}
//
//	return dataBytes, err
//}
//
//const headerSize = 4 // 消息头的大小，用于表示消息的长度
//
//// 打包消息
//func (p *BufferPacker) PackMessage(msg interface{}) ([]byte, error) {
//	pbMsg, ok := msg.(proto.Message)
//	if !ok {
//		log.Fatal("")
//	}
//	// data
//	message, err := proto.Marshal(pbMsg)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 创建一个字节数组，用于存放打包后的消息
//	packedMessage := make([]byte, headerSize+len(message))
//
//	// 将消息的长度转换为字节数组，并拷贝到打包后的消息中
//	binary.BigEndian.PutUint32(packedMessage[:headerSize], uint32(len(message)))
//
//	// 将消息内容拷贝到打包后的消息中
//	copy(packedMessage[headerSize:], message)
//
//	return packedMessage, nil
//}
//
//// 解包消息
//func (p *BufferPacker) UnpackMessage(conn net.Conn) ([]byte, error) {
//	// 读取消息头，获取消息的长度
//	header := make([]byte, headerSize)
//	if _, err := conn.Read(header); err != nil {
//		return nil, err
//	}
//	messageSize := binary.BigEndian.Uint32(header)
//
//	// 读取消息内容
//	message := make([]byte, messageSize)
//	if _, err := conn.Read(message); err != nil {
//		return nil, err
//	}
//
//	return message, nil
//}
