// Package rpc @Author evan_yxw
// @Date 2024/8/16 16:12:00
// @Desc
package rpc

import (
	"fmt"
	"github.com/evanyxw/monster-go/message/pb/xsf_pb"
	"github.com/golang/protobuf/proto"
	"reflect"
)

// 包级变量，初始化一次
var typeToMsgIDMap map[reflect.Type]int32

func init() {
	typeToMsgIDMap = make(map[reflect.Type]int32)

	for _, messageValue := range xsf_pb.SMSGID_value {
		messageName := xsf_pb.SMSGID_name[int32(messageValue)]
		messageType := proto.MessageType("NLD_PB." + messageName)
		if messageType != nil {
			typeToMsgIDMap[messageType] = int32(messageValue)
		}
	}

	for _, messageValue := range xsf_pb.MSGID_value {
		messageName := xsf_pb.MSGID_name[int32(messageValue)]
		messageType := proto.MessageType("NLD_PB." + messageName)
		if messageType != nil {
			typeToMsgIDMap[messageType] = int32(messageValue)
		}
	}
}

// 根据 proto.Message 获取对应的 msgID
func GetMsgID(message proto.Message) (uint64, error) {
	messageType := reflect.TypeOf(message)
	if msgID, found := typeToMsgIDMap[messageType]; found {
		return uint64(msgID), nil
	}
	return 0, fmt.Errorf("未找到对应的 msgID, 消息类型: %v", messageType)
}

func GetMessage(messageID uint64) (interface{}, error) {
	// 利用反射获取消息 ID 对应的协议结构体
	//messageType2 := reflect.TypeOf((*proto.Message)(nil)).Elem()
	for _, messageValue := range xsf_pb.SMSGID_value {
		if messageValue == int32(messageID) {
			messageName := xsf_pb.SMSGID_name[int32(messageValue)]
			messageType := proto.MessageType("NLD_PB." + messageName)
			if messageType == nil {
				return nil, fmt.Errorf("未找到对应的协议结构体")
			}

			instance := reflect.New(messageType.Elem()).Interface()
			msg, ok := instance.(proto.Message)
			if !ok {
				return nil, fmt.Errorf("实例化失败")
			}

			return msg, nil
		}
	}

	for _, messageValue := range xsf_pb.MSGID_value {
		if messageValue == int32(messageID) {
			messageName := xsf_pb.MSGID_name[int32(messageValue)]
			messageType := proto.MessageType("NLD_PB." + messageName)
			if messageType == nil {
				return nil, fmt.Errorf("未找到对应的协议结构体")
			}

			instance := reflect.New(messageType.Elem()).Interface()
			msg, ok := instance.(proto.Message)
			if !ok {
				return nil, fmt.Errorf("实例化失败")
			}

			return msg, nil
		}
	}

	return nil, fmt.Errorf("未找到对应的协议结构体")
}

func Import(data []byte, msg proto.Message) {
	proto.Unmarshal(data, msg)
}

//func GetMessage(messageID uint64) (proto.Message, error) {
//	// 创建一个 map 来映射 messageID 到消息类型
//	messageTypeMap := make(map[int32]reflect.Type)
//
//	for _, messageValue := range xsf_pb.SMSGID_value {
//		messageName := xsf_pb.SMSGID_name[int32(messageValue)]
//		messageType := proto.MessageType("NLD_PB." + messageName)
//		if messageType != nil {
//			messageTypeMap[int32(messageValue)] = messageType.Elem()
//		}
//	}
//
//	for _, messageValue := range xsf_pb.MSGID_value {
//		messageName := xsf_pb.MSGID_name[int32(messageValue)]
//		messageType := proto.MessageType("NLD_PB." + messageName)
//		if messageType != nil {
//			messageTypeMap[int32(messageValue)] = messageType.Elem()
//		}
//	}
//
//	// 查找消息类型并创建实例
//	if messageType, found := messageTypeMap[int32(messageID)]; found {
//		instance := reflect.New(messageType).Interface()
//		msg, ok := instance.(proto.Message)
//		if !ok {
//			return nil, fmt.Errorf("实例化失败")
//		}
//		return msg, nil
//	}
//
//	return nil, fmt.Errorf("未找到对应的协议结构体")
//}
