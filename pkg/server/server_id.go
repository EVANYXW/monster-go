package server

import (
	"encoding/binary"
)

type ServerID struct {
	ID    uint16 // 唯一ID
	Type  uint8  // 服务器类型
	Index uint8  // 索引
}

type ClientID struct {
	ID   uint16 // 唯一索引ID
	Gate uint8  // 网关ID
	Key  uint8  // 标识ID
}

var (
	ServerTag  string = "NLD" // 服务器启动tag
	ServerName string
	ServerDesc string
	AutoStart  bool        // 是否自动开服
	DebugRun   bool = true // 是否是调试模式运行

	RunDirPrefix string // 运行时目录前缀

	LogIndex int32

	ID    uint32         = 0 // 服务器ID
	SID   ServerID           // 服务器复合ID
	Ports [EP_Max]uint32     // 端口信息

	//Status xsf_util.AtomicInt32
)

func UpdateSID() {
	ID2Sid(ID, &SID)
}

func UpdateID() {
	ID = Sid2ID(&SID)
}

// 整形ID转复合ID
func ID2Sid(id uint32, pSID *ServerID) {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(id))

	var idPart [2]byte
	idPart[0] = buf[2]
	idPart[1] = buf[3]
	pSID.ID = binary.LittleEndian.Uint16(idPart[:])

	pSID.Type = buf[1]
	pSID.Index = buf[0]
}

// 复合ID转整形ID
func Sid2ID(pSID *ServerID) uint32 {
	var idPart [2]byte
	binary.LittleEndian.PutUint16(idPart[:], pSID.ID)

	var buf [4]byte
	buf[0] = pSID.Index
	buf[1] = pSID.Type
	buf[2] = idPart[0]
	buf[3] = idPart[1]

	return binary.LittleEndian.Uint32(buf[:])
}

// 整形ID转客户端ID
func ID2Cid(id uint32, pCID *ClientID) {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], id)

	var idPart [2]byte
	idPart[0] = buf[0]
	idPart[1] = buf[1]
	pCID.ID = binary.BigEndian.Uint16(idPart[:])
	pCID.Gate = buf[2]
	pCID.Key = buf[3]
}

// 客户端ID转整形ID
func Cid2ID(pCID *ClientID) uint32 {
	var idPart [2]byte
	binary.BigEndian.PutUint16(idPart[:], pCID.ID)

	var buf [4]byte
	buf[0] = idPart[0]
	buf[1] = idPart[1]
	buf[2] = pCID.Gate
	buf[3] = pCID.Key

	return binary.BigEndian.Uint32(buf[:])
}
