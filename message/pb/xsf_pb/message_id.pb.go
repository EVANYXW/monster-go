// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v4.23.4
// source: message_id.proto

package xsf_pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MSGID int32

const (
	MSGID_None MSGID = 0
	// MESSAGE_ID_START
	MSGID_Clt_Gt_Handshake           MSGID = 1  // client --> Gate 握手
	MSGID_Gt_Clt_Handshake           MSGID = 2  // gate --> client 握手反馈
	MSGID_Clt_Gt_Heartbeat           MSGID = 3  // client --> gate 心跳
	MSGID_Gt_Clt_Heartbeat           MSGID = 4  // gate --> client 心跳反馈
	MSGID_Gt_Clt_Disconnect          MSGID = 5  // gate --> client 服务器断开连接
	MSGID_Clt_L_Login                MSGID = 6  // client --> login 登录
	MSGID_L_Clt_LoginResult          MSGID = 7  // login --> client 登录结果
	MSGID_G_Clt_LoginResult          MSGID = 9  // game --> client 登录结果
	MSGID_G_Clt_OpResult             MSGID = 10 // game --> client 通用操作反馈
	MSGID_G_Clt_ActorBase            MSGID = 11 // game --> client 基础数据
	MSGID_G_Clt_ActorItem            MSGID = 12 // game --> client 物品数据
	MSGID_G_Clt_ActorCharacter       MSGID = 13 // game --> client 角色数据
	MSGID_Clt_G_CharacterUpgrade     MSGID = 14 // client --> game 角色提升
	MSGID_G_Clt_ActorLogic           MSGID = 15 // game --> client 逻辑计数数据
	MSGID_G_Clt_ActorTask            MSGID = 16 // game --> client 任务数据
	MSGID_Clt_G_SetNickname          MSGID = 17 // client --> game 设置昵称
	MSGID_G_Clt_ActorLevel           MSGID = 18 // game --> client 关卡数据
	MSGID_Clt_G_EnterLevel           MSGID = 19 // client --> game 进入关卡
	MSGID_Clt_G_ExitLevel            MSGID = 20 // client --> game 退出关卡个
	MSGID_G_Clt_ActorRaffle          MSGID = 21 // game --> client 抽奖数据
	MSGID_Clt_G_Raffle               MSGID = 22 // client --> game 抽奖
	MSGID_G_Clt_ItemList             MSGID = 23 // game --> client 物品列表
	MSGID_Clt_G_GroupReward          MSGID = 24 // client --> game 领取章节宝箱
	MSGID_Clt_G_VehicleUpgrade       MSGID = 25 // client --> game 载具升级
	MSGID_G_Clt_ActorTeam            MSGID = 26 // game --> client 玩家编组数据
	MSGID_Clt_G_CreateTeam           MSGID = 27 // client --> game 创建队伍
	MSGID_Clt_G_UseItem              MSGID = 28 // client --> game 使用唯品
	MSGID_Clt_G_TestInit             MSGID = 29 // client --> game 测试初始化
	MSGID_Clt_G_GiftCode             MSGID = 30 // client --> game 礼包码兑换
	MSGID_Ml_Clt_MailInfo            MSGID = 31 // mail --> client 邮件简要信息
	MSGID_Clt_Ml_MailRequest         MSGID = 32 // client --> mail 请求邮件数据
	MSGID_Ml_Clt_MailData            MSGID = 33 // mail --> client 邮件数据
	MSGID_Clt_Ml_OpenMail            MSGID = 34 // client --> mail 打开邮件
	MSGID_Ml_Clt_NewGlobalMail       MSGID = 35 // mail --> client 有新的全局邮件
	MSGID_G_Clt_LoginKey             MSGID = 36 // game --> client 登录key更新
	MSGID_Clt_G_Relogin              MSGID = 37 // client --> game 重新登录
	MSGID_Clt_Ml_ClaimMailItems      MSGID = 38 // client --> mail 领取邮件物品
	MSGID_Clt_Ml_DeleteMail          MSGID = 39 // client --> mail 删除邮件
	MSGID_Clt_G_TaskClaim            MSGID = 40 // client --> game 领取任务奖励
	MSGID_Clt_G_GetShopContent       MSGID = 41 // client --> game 获取商店数据
	MSGID_G_Clt_ShopDatas            MSGID = 42 // game --> client 商店数据
	MSGID_Clt_G_ShopPurchase         MSGID = 43 // client --> game 购买商品
	MSGID_G_Clt_ActorPlayerSkill     MSGID = 44 // game --> client 指挥官技能数据
	MSGID_Clt_G_PlayerSkillUpgrade   MSGID = 45 // client --> game 升级指挥官技能
	MSGID_Clt_G_SetDefenceTeamForPvp MSGID = 46 // client --> game 设置PVP防守阵营
	MSGID_Clt_Mg_PvpInformation      MSGID = 47 // client --> manager 请求PVP数据
	MSGID_G_Clt_PvpInformation       MSGID = 48 // game --> client PVP数据反馈
	MSGID_Clt_G_PvpStart             MSGID = 49 // client --> game 开始PVP战斗
	MSGID_Clt_G_PvpResult            MSGID = 50 // client --> game PVP战斗结果
	MSGID_G_Clt_ActorPvp             MSGID = 51 // game --> client 玩家pvp数据
	MSGID_G_Clt_PvpDefenceData       MSGID = 52 // game --> client 防守数据
	MSGID_Clt_Mg_PvpOrderList        MSGID = 53 // client --> manager 请求pvp排行榜
	MSGID_Mg_Clt_PvpOrderListData    MSGID = 54 // manager --> client pvp排行榜数据
	MSGID_Clt_G_PatioExchange        MSGID = 55 // client --> game 天井兑换
	MSGID_Clt_G_InAppPurchase        MSGID = 56 // client --> game 请求内购
	MSGID_G_Clt_ActorPurchase        MSGID = 57 // game --> client 玩家内购数据
	MSGID_Clt_G_ExchangeItem         MSGID = 58 // client --> game 道具兑换道具
	MSGID_Clt_G_NewGuideProgress     MSGID = 59 // client --> game 新手引导进度
	MSGID_Clt_G_SetHeadicon          MSGID = 60 // client --> game 设置形象id
	MSGID_Clt_G_BuyAlbum             MSGID = 61 // client --> game 购买音乐专辑
	MSGID_Clt_G_AdultAuth            MSGID = 62 // client --> game 请求防沉迷认证
	MSGID_MSGID_Max                  MSGID = 100
)

// Enum value maps for MSGID.
var (
	MSGID_name = map[int32]string{
		0:   "None",
		1:   "Clt_Gt_Handshake",
		2:   "Gt_Clt_Handshake",
		3:   "Clt_Gt_Heartbeat",
		4:   "Gt_Clt_Heartbeat",
		5:   "Gt_Clt_Disconnect",
		6:   "Clt_L_Login",
		7:   "L_Clt_LoginResult",
		9:   "G_Clt_LoginResult",
		10:  "G_Clt_OpResult",
		11:  "G_Clt_ActorBase",
		12:  "G_Clt_ActorItem",
		13:  "G_Clt_ActorCharacter",
		14:  "Clt_G_CharacterUpgrade",
		15:  "G_Clt_ActorLogic",
		16:  "G_Clt_ActorTask",
		17:  "Clt_G_SetNickname",
		18:  "G_Clt_ActorLevel",
		19:  "Clt_G_EnterLevel",
		20:  "Clt_G_ExitLevel",
		21:  "G_Clt_ActorRaffle",
		22:  "Clt_G_Raffle",
		23:  "G_Clt_ItemList",
		24:  "Clt_G_GroupReward",
		25:  "Clt_G_VehicleUpgrade",
		26:  "G_Clt_ActorTeam",
		27:  "Clt_G_CreateTeam",
		28:  "Clt_G_UseItem",
		29:  "Clt_G_TestInit",
		30:  "Clt_G_GiftCode",
		31:  "Ml_Clt_MailInfo",
		32:  "Clt_Ml_MailRequest",
		33:  "Ml_Clt_MailData",
		34:  "Clt_Ml_OpenMail",
		35:  "Ml_Clt_NewGlobalMail",
		36:  "G_Clt_LoginKey",
		37:  "Clt_G_Relogin",
		38:  "Clt_Ml_ClaimMailItems",
		39:  "Clt_Ml_DeleteMail",
		40:  "Clt_G_TaskClaim",
		41:  "Clt_G_GetShopContent",
		42:  "G_Clt_ShopDatas",
		43:  "Clt_G_ShopPurchase",
		44:  "G_Clt_ActorPlayerSkill",
		45:  "Clt_G_PlayerSkillUpgrade",
		46:  "Clt_G_SetDefenceTeamForPvp",
		47:  "Clt_Mg_PvpInformation",
		48:  "G_Clt_PvpInformation",
		49:  "Clt_G_PvpStart",
		50:  "Clt_G_PvpResult",
		51:  "G_Clt_ActorPvp",
		52:  "G_Clt_PvpDefenceData",
		53:  "Clt_Mg_PvpOrderList",
		54:  "Mg_Clt_PvpOrderListData",
		55:  "Clt_G_PatioExchange",
		56:  "Clt_G_InAppPurchase",
		57:  "G_Clt_ActorPurchase",
		58:  "Clt_G_ExchangeItem",
		59:  "Clt_G_NewGuideProgress",
		60:  "Clt_G_SetHeadicon",
		61:  "Clt_G_BuyAlbum",
		62:  "Clt_G_AdultAuth",
		100: "MSGID_Max",
	}
	MSGID_value = map[string]int32{
		"None":                       0,
		"Clt_Gt_Handshake":           1,
		"Gt_Clt_Handshake":           2,
		"Clt_Gt_Heartbeat":           3,
		"Gt_Clt_Heartbeat":           4,
		"Gt_Clt_Disconnect":          5,
		"Clt_L_Login":                6,
		"L_Clt_LoginResult":          7,
		"G_Clt_LoginResult":          9,
		"G_Clt_OpResult":             10,
		"G_Clt_ActorBase":            11,
		"G_Clt_ActorItem":            12,
		"G_Clt_ActorCharacter":       13,
		"Clt_G_CharacterUpgrade":     14,
		"G_Clt_ActorLogic":           15,
		"G_Clt_ActorTask":            16,
		"Clt_G_SetNickname":          17,
		"G_Clt_ActorLevel":           18,
		"Clt_G_EnterLevel":           19,
		"Clt_G_ExitLevel":            20,
		"G_Clt_ActorRaffle":          21,
		"Clt_G_Raffle":               22,
		"G_Clt_ItemList":             23,
		"Clt_G_GroupReward":          24,
		"Clt_G_VehicleUpgrade":       25,
		"G_Clt_ActorTeam":            26,
		"Clt_G_CreateTeam":           27,
		"Clt_G_UseItem":              28,
		"Clt_G_TestInit":             29,
		"Clt_G_GiftCode":             30,
		"Ml_Clt_MailInfo":            31,
		"Clt_Ml_MailRequest":         32,
		"Ml_Clt_MailData":            33,
		"Clt_Ml_OpenMail":            34,
		"Ml_Clt_NewGlobalMail":       35,
		"G_Clt_LoginKey":             36,
		"Clt_G_Relogin":              37,
		"Clt_Ml_ClaimMailItems":      38,
		"Clt_Ml_DeleteMail":          39,
		"Clt_G_TaskClaim":            40,
		"Clt_G_GetShopContent":       41,
		"G_Clt_ShopDatas":            42,
		"Clt_G_ShopPurchase":         43,
		"G_Clt_ActorPlayerSkill":     44,
		"Clt_G_PlayerSkillUpgrade":   45,
		"Clt_G_SetDefenceTeamForPvp": 46,
		"Clt_Mg_PvpInformation":      47,
		"G_Clt_PvpInformation":       48,
		"Clt_G_PvpStart":             49,
		"Clt_G_PvpResult":            50,
		"G_Clt_ActorPvp":             51,
		"G_Clt_PvpDefenceData":       52,
		"Clt_Mg_PvpOrderList":        53,
		"Mg_Clt_PvpOrderListData":    54,
		"Clt_G_PatioExchange":        55,
		"Clt_G_InAppPurchase":        56,
		"G_Clt_ActorPurchase":        57,
		"Clt_G_ExchangeItem":         58,
		"Clt_G_NewGuideProgress":     59,
		"Clt_G_SetHeadicon":          60,
		"Clt_G_BuyAlbum":             61,
		"Clt_G_AdultAuth":            62,
		"MSGID_Max":                  100,
	}
)

func (x MSGID) Enum() *MSGID {
	p := new(MSGID)
	*p = x
	return p
}

func (x MSGID) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MSGID) Descriptor() protoreflect.EnumDescriptor {
	return file_message_id_proto_enumTypes[0].Descriptor()
}

func (MSGID) Type() protoreflect.EnumType {
	return &file_message_id_proto_enumTypes[0]
}

func (x MSGID) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MSGID.Descriptor instead.
func (MSGID) EnumDescriptor() ([]byte, []int) {
	return file_message_id_proto_rawDescGZIP(), []int{0}
}

var File_message_id_proto protoreflect.FileDescriptor

var file_message_id_proto_rawDesc = []byte{
	0x0a, 0x10, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x07, 0x4e, 0x4c, 0x44, 0x5f, 0x4d, 0x49, 0x44, 0x2a, 0x9a, 0x0b, 0x0a, 0x05,
	0x4d, 0x53, 0x47, 0x49, 0x44, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x6f, 0x6e, 0x65, 0x10, 0x00, 0x12,
	0x14, 0x0a, 0x10, 0x43, 0x6c, 0x74, 0x5f, 0x47, 0x74, 0x5f, 0x48, 0x61, 0x6e, 0x64, 0x73, 0x68,
	0x61, 0x6b, 0x65, 0x10, 0x01, 0x12, 0x14, 0x0a, 0x10, 0x47, 0x74, 0x5f, 0x43, 0x6c, 0x74, 0x5f,
	0x48, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x10, 0x02, 0x12, 0x14, 0x0a, 0x10, 0x43,
	0x6c, 0x74, 0x5f, 0x47, 0x74, 0x5f, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x10,
	0x03, 0x12, 0x14, 0x0a, 0x10, 0x47, 0x74, 0x5f, 0x43, 0x6c, 0x74, 0x5f, 0x48, 0x65, 0x61, 0x72,
	0x74, 0x62, 0x65, 0x61, 0x74, 0x10, 0x04, 0x12, 0x15, 0x0a, 0x11, 0x47, 0x74, 0x5f, 0x43, 0x6c,
	0x74, 0x5f, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x10, 0x05, 0x12, 0x0f,
	0x0a, 0x0b, 0x43, 0x6c, 0x74, 0x5f, 0x4c, 0x5f, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x10, 0x06, 0x12,
	0x15, 0x0a, 0x11, 0x4c, 0x5f, 0x43, 0x6c, 0x74, 0x5f, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x10, 0x07, 0x12, 0x15, 0x0a, 0x11, 0x47, 0x5f, 0x43, 0x6c, 0x74, 0x5f,
	0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x10, 0x09, 0x12, 0x12, 0x0a,
	0x0e, 0x47, 0x5f, 0x43, 0x6c, 0x74, 0x5f, 0x4f, 0x70, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x10,
	0x0a, 0x12, 0x13, 0x0a, 0x0f, 0x47, 0x5f, 0x43, 0x6c, 0x74, 0x5f, 0x41, 0x63, 0x74, 0x6f, 0x72,
	0x42, 0x61, 0x73, 0x65, 0x10, 0x0b, 0x12, 0x13, 0x0a, 0x0f, 0x47, 0x5f, 0x43, 0x6c, 0x74, 0x5f,
	0x41, 0x63, 0x74, 0x6f, 0x72, 0x49, 0x74, 0x65, 0x6d, 0x10, 0x0c, 0x12, 0x18, 0x0a, 0x14, 0x47,
	0x5f, 0x43, 0x6c, 0x74, 0x5f, 0x41, 0x63, 0x74, 0x6f, 0x72, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63,
	0x74, 0x65, 0x72, 0x10, 0x0d, 0x12, 0x1a, 0x0a, 0x16, 0x43, 0x6c, 0x74, 0x5f, 0x47, 0x5f, 0x43,
	0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x55, 0x70, 0x67, 0x72, 0x61, 0x64, 0x65, 0x10,
	0x0e, 0x12, 0x14, 0x0a, 0x10, 0x47, 0x5f, 0x43, 0x6c, 0x74, 0x5f, 0x41, 0x63, 0x74, 0x6f, 0x72,
	0x4c, 0x6f, 0x67, 0x69, 0x63, 0x10, 0x0f, 0x12, 0x13, 0x0a, 0x0f, 0x47, 0x5f, 0x43, 0x6c, 0x74,
	0x5f, 0x41, 0x63, 0x74, 0x6f, 0x72, 0x54, 0x61, 0x73, 0x6b, 0x10, 0x10, 0x12, 0x15, 0x0a, 0x11,
	0x43, 0x6c, 0x74, 0x5f, 0x47, 0x5f, 0x53, 0x65, 0x74, 0x4e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d,
	0x65, 0x10, 0x11, 0x12, 0x14, 0x0a, 0x10, 0x47, 0x5f, 0x43, 0x6c, 0x74, 0x5f, 0x41, 0x63, 0x74,
	0x6f, 0x72, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x10, 0x12, 0x12, 0x14, 0x0a, 0x10, 0x43, 0x6c, 0x74,
	0x5f, 0x47, 0x5f, 0x45, 0x6e, 0x74, 0x65, 0x72, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x10, 0x13, 0x12,
	0x13, 0x0a, 0x0f, 0x43, 0x6c, 0x74, 0x5f, 0x47, 0x5f, 0x45, 0x78, 0x69, 0x74, 0x4c, 0x65, 0x76,
	0x65, 0x6c, 0x10, 0x14, 0x12, 0x15, 0x0a, 0x11, 0x47, 0x5f, 0x43, 0x6c, 0x74, 0x5f, 0x41, 0x63,
	0x74, 0x6f, 0x72, 0x52, 0x61, 0x66, 0x66, 0x6c, 0x65, 0x10, 0x15, 0x12, 0x10, 0x0a, 0x0c, 0x43,
	0x6c, 0x74, 0x5f, 0x47, 0x5f, 0x52, 0x61, 0x66, 0x66, 0x6c, 0x65, 0x10, 0x16, 0x12, 0x12, 0x0a,
	0x0e, 0x47, 0x5f, 0x43, 0x6c, 0x74, 0x5f, 0x49, 0x74, 0x65, 0x6d, 0x4c, 0x69, 0x73, 0x74, 0x10,
	0x17, 0x12, 0x15, 0x0a, 0x11, 0x43, 0x6c, 0x74, 0x5f, 0x47, 0x5f, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x52, 0x65, 0x77, 0x61, 0x72, 0x64, 0x10, 0x18, 0x12, 0x18, 0x0a, 0x14, 0x43, 0x6c, 0x74, 0x5f,
	0x47, 0x5f, 0x56, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x55, 0x70, 0x67, 0x72, 0x61, 0x64, 0x65,
	0x10, 0x19, 0x12, 0x13, 0x0a, 0x0f, 0x47, 0x5f, 0x43, 0x6c, 0x74, 0x5f, 0x41, 0x63, 0x74, 0x6f,
	0x72, 0x54, 0x65, 0x61, 0x6d, 0x10, 0x1a, 0x12, 0x14, 0x0a, 0x10, 0x43, 0x6c, 0x74, 0x5f, 0x47,
	0x5f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x65, 0x61, 0x6d, 0x10, 0x1b, 0x12, 0x11, 0x0a,
	0x0d, 0x43, 0x6c, 0x74, 0x5f, 0x47, 0x5f, 0x55, 0x73, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x10, 0x1c,
	0x12, 0x12, 0x0a, 0x0e, 0x43, 0x6c, 0x74, 0x5f, 0x47, 0x5f, 0x54, 0x65, 0x73, 0x74, 0x49, 0x6e,
	0x69, 0x74, 0x10, 0x1d, 0x12, 0x12, 0x0a, 0x0e, 0x43, 0x6c, 0x74, 0x5f, 0x47, 0x5f, 0x47, 0x69,
	0x66, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x10, 0x1e, 0x12, 0x13, 0x0a, 0x0f, 0x4d, 0x6c, 0x5f, 0x43,
	0x6c, 0x74, 0x5f, 0x4d, 0x61, 0x69, 0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x10, 0x1f, 0x12, 0x16, 0x0a,
	0x12, 0x43, 0x6c, 0x74, 0x5f, 0x4d, 0x6c, 0x5f, 0x4d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x10, 0x20, 0x12, 0x13, 0x0a, 0x0f, 0x4d, 0x6c, 0x5f, 0x43, 0x6c, 0x74, 0x5f,
	0x4d, 0x61, 0x69, 0x6c, 0x44, 0x61, 0x74, 0x61, 0x10, 0x21, 0x12, 0x13, 0x0a, 0x0f, 0x43, 0x6c,
	0x74, 0x5f, 0x4d, 0x6c, 0x5f, 0x4f, 0x70, 0x65, 0x6e, 0x4d, 0x61, 0x69, 0x6c, 0x10, 0x22, 0x12,
	0x18, 0x0a, 0x14, 0x4d, 0x6c, 0x5f, 0x43, 0x6c, 0x74, 0x5f, 0x4e, 0x65, 0x77, 0x47, 0x6c, 0x6f,
	0x62, 0x61, 0x6c, 0x4d, 0x61, 0x69, 0x6c, 0x10, 0x23, 0x12, 0x12, 0x0a, 0x0e, 0x47, 0x5f, 0x43,
	0x6c, 0x74, 0x5f, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x4b, 0x65, 0x79, 0x10, 0x24, 0x12, 0x11, 0x0a,
	0x0d, 0x43, 0x6c, 0x74, 0x5f, 0x47, 0x5f, 0x52, 0x65, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x10, 0x25,
	0x12, 0x19, 0x0a, 0x15, 0x43, 0x6c, 0x74, 0x5f, 0x4d, 0x6c, 0x5f, 0x43, 0x6c, 0x61, 0x69, 0x6d,
	0x4d, 0x61, 0x69, 0x6c, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x10, 0x26, 0x12, 0x15, 0x0a, 0x11, 0x43,
	0x6c, 0x74, 0x5f, 0x4d, 0x6c, 0x5f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4d, 0x61, 0x69, 0x6c,
	0x10, 0x27, 0x12, 0x13, 0x0a, 0x0f, 0x43, 0x6c, 0x74, 0x5f, 0x47, 0x5f, 0x54, 0x61, 0x73, 0x6b,
	0x43, 0x6c, 0x61, 0x69, 0x6d, 0x10, 0x28, 0x12, 0x18, 0x0a, 0x14, 0x43, 0x6c, 0x74, 0x5f, 0x47,
	0x5f, 0x47, 0x65, 0x74, 0x53, 0x68, 0x6f, 0x70, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x10,
	0x29, 0x12, 0x13, 0x0a, 0x0f, 0x47, 0x5f, 0x43, 0x6c, 0x74, 0x5f, 0x53, 0x68, 0x6f, 0x70, 0x44,
	0x61, 0x74, 0x61, 0x73, 0x10, 0x2a, 0x12, 0x16, 0x0a, 0x12, 0x43, 0x6c, 0x74, 0x5f, 0x47, 0x5f,
	0x53, 0x68, 0x6f, 0x70, 0x50, 0x75, 0x72, 0x63, 0x68, 0x61, 0x73, 0x65, 0x10, 0x2b, 0x12, 0x1a,
	0x0a, 0x16, 0x47, 0x5f, 0x43, 0x6c, 0x74, 0x5f, 0x41, 0x63, 0x74, 0x6f, 0x72, 0x50, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x53, 0x6b, 0x69, 0x6c, 0x6c, 0x10, 0x2c, 0x12, 0x1c, 0x0a, 0x18, 0x43, 0x6c,
	0x74, 0x5f, 0x47, 0x5f, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x53, 0x6b, 0x69, 0x6c, 0x6c, 0x55,
	0x70, 0x67, 0x72, 0x61, 0x64, 0x65, 0x10, 0x2d, 0x12, 0x1e, 0x0a, 0x1a, 0x43, 0x6c, 0x74, 0x5f,
	0x47, 0x5f, 0x53, 0x65, 0x74, 0x44, 0x65, 0x66, 0x65, 0x6e, 0x63, 0x65, 0x54, 0x65, 0x61, 0x6d,
	0x46, 0x6f, 0x72, 0x50, 0x76, 0x70, 0x10, 0x2e, 0x12, 0x19, 0x0a, 0x15, 0x43, 0x6c, 0x74, 0x5f,
	0x4d, 0x67, 0x5f, 0x50, 0x76, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x10, 0x2f, 0x12, 0x18, 0x0a, 0x14, 0x47, 0x5f, 0x43, 0x6c, 0x74, 0x5f, 0x50, 0x76, 0x70,
	0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x10, 0x30, 0x12, 0x12, 0x0a,
	0x0e, 0x43, 0x6c, 0x74, 0x5f, 0x47, 0x5f, 0x50, 0x76, 0x70, 0x53, 0x74, 0x61, 0x72, 0x74, 0x10,
	0x31, 0x12, 0x13, 0x0a, 0x0f, 0x43, 0x6c, 0x74, 0x5f, 0x47, 0x5f, 0x50, 0x76, 0x70, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x10, 0x32, 0x12, 0x12, 0x0a, 0x0e, 0x47, 0x5f, 0x43, 0x6c, 0x74, 0x5f,
	0x41, 0x63, 0x74, 0x6f, 0x72, 0x50, 0x76, 0x70, 0x10, 0x33, 0x12, 0x18, 0x0a, 0x14, 0x47, 0x5f,
	0x43, 0x6c, 0x74, 0x5f, 0x50, 0x76, 0x70, 0x44, 0x65, 0x66, 0x65, 0x6e, 0x63, 0x65, 0x44, 0x61,
	0x74, 0x61, 0x10, 0x34, 0x12, 0x17, 0x0a, 0x13, 0x43, 0x6c, 0x74, 0x5f, 0x4d, 0x67, 0x5f, 0x50,
	0x76, 0x70, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x4c, 0x69, 0x73, 0x74, 0x10, 0x35, 0x12, 0x1b, 0x0a,
	0x17, 0x4d, 0x67, 0x5f, 0x43, 0x6c, 0x74, 0x5f, 0x50, 0x76, 0x70, 0x4f, 0x72, 0x64, 0x65, 0x72,
	0x4c, 0x69, 0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x10, 0x36, 0x12, 0x17, 0x0a, 0x13, 0x43, 0x6c,
	0x74, 0x5f, 0x47, 0x5f, 0x50, 0x61, 0x74, 0x69, 0x6f, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x10, 0x37, 0x12, 0x17, 0x0a, 0x13, 0x43, 0x6c, 0x74, 0x5f, 0x47, 0x5f, 0x49, 0x6e, 0x41,
	0x70, 0x70, 0x50, 0x75, 0x72, 0x63, 0x68, 0x61, 0x73, 0x65, 0x10, 0x38, 0x12, 0x17, 0x0a, 0x13,
	0x47, 0x5f, 0x43, 0x6c, 0x74, 0x5f, 0x41, 0x63, 0x74, 0x6f, 0x72, 0x50, 0x75, 0x72, 0x63, 0x68,
	0x61, 0x73, 0x65, 0x10, 0x39, 0x12, 0x16, 0x0a, 0x12, 0x43, 0x6c, 0x74, 0x5f, 0x47, 0x5f, 0x45,
	0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x10, 0x3a, 0x12, 0x1a, 0x0a,
	0x16, 0x43, 0x6c, 0x74, 0x5f, 0x47, 0x5f, 0x4e, 0x65, 0x77, 0x47, 0x75, 0x69, 0x64, 0x65, 0x50,
	0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x10, 0x3b, 0x12, 0x15, 0x0a, 0x11, 0x43, 0x6c, 0x74,
	0x5f, 0x47, 0x5f, 0x53, 0x65, 0x74, 0x48, 0x65, 0x61, 0x64, 0x69, 0x63, 0x6f, 0x6e, 0x10, 0x3c,
	0x12, 0x12, 0x0a, 0x0e, 0x43, 0x6c, 0x74, 0x5f, 0x47, 0x5f, 0x42, 0x75, 0x79, 0x41, 0x6c, 0x62,
	0x75, 0x6d, 0x10, 0x3d, 0x12, 0x13, 0x0a, 0x0f, 0x43, 0x6c, 0x74, 0x5f, 0x47, 0x5f, 0x41, 0x64,
	0x75, 0x6c, 0x74, 0x41, 0x75, 0x74, 0x68, 0x10, 0x3e, 0x12, 0x0d, 0x0a, 0x09, 0x4d, 0x53, 0x47,
	0x49, 0x44, 0x5f, 0x4d, 0x61, 0x78, 0x10, 0x64, 0x42, 0x0b, 0x5a, 0x09, 0x70, 0x62, 0x2f, 0x78,
	0x73, 0x66, 0x5f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_message_id_proto_rawDescOnce sync.Once
	file_message_id_proto_rawDescData = file_message_id_proto_rawDesc
)

func file_message_id_proto_rawDescGZIP() []byte {
	file_message_id_proto_rawDescOnce.Do(func() {
		file_message_id_proto_rawDescData = protoimpl.X.CompressGZIP(file_message_id_proto_rawDescData)
	})
	return file_message_id_proto_rawDescData
}

var file_message_id_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_message_id_proto_goTypes = []interface{}{
	(MSGID)(0), // 0: NLD_MID.MSGID
}
var file_message_id_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_message_id_proto_init() }
func file_message_id_proto_init() {
	if File_message_id_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_message_id_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_message_id_proto_goTypes,
		DependencyIndexes: file_message_id_proto_depIdxs,
		EnumInfos:         file_message_id_proto_enumTypes,
	}.Build()
	File_message_id_proto = out.File
	file_message_id_proto_rawDesc = nil
	file_message_id_proto_goTypes = nil
	file_message_id_proto_depIdxs = nil
}
