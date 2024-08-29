// ////////////////////////////////////////////////////////////////////////
//
// 文件：common/db/db_def.go
// 作者：xiexx
// 时间：2023/07/17
// 描述：db相关定义
// 说明：
//
// ////////////////////////////////////////////////////////////////////////
package db

import "github.com/evanyxw/monster-go/message/pb/xsf_pb"

var actorBaseProps []string

type AccountInfo struct {
	Id      string `bson:"id"`       // 账号唯一ID，或者sdk的token
	ActorID uint32 `bson:"actor_id"` // 角色ID
	Time    uint32 `bson:"time"`     // 创建时间
	GameID  uint32 `bson:"game_id"`  // 服务器ID
}

const (
	// REDIS_REQUEST_ID_BEGIN
	REDIS_GET_INIT_FLAG          = 10000 // 获取初始化标记
	REDIS_SET_INIT_FLAG          = 10001 // 设置初始化标记
	REDIS_GET_ACCOUNT_INFO       = 20000 // 获取账号信息
	REDIS_SET_ACCOUNT_INFO       = 20001 // 设置账号信息
	REDIS_BATCH_SET_ACCOUNT_INFO = 20002 // 批量设置账号信息
	// REDIS_REQUEST_ID_END
)

const (
	// MONGO_REQUEST_ID_BEGIN
	MONGO_GET_ALL_ACCOUNT        = 10000  // 查找所有账号数据
	MONGO_GET_MAX_ACTOR_ID       = 10001  // 查找最大角色id
	MONGO_GET_ACCOUNT            = 10002  // 查找账号
	MONGO_NEW_ACCOUNT            = 10003  // 新建账号
	MONGO_GET_ACTOR_BASE         = 20000  // 获取玩家基础数据
	MONGO_UPDATE_ACTOR_BASE      = 20001  // 更新玩家数据
	MONGO_GET_ALL_ACTORS         = 20002  // 查找所有玩家数据
	MONGO_GET_ACTOR_ITEM         = 30000  // 获取玩家物品数据
	MONGO_UPDATE_ACTOR_ITEM      = 30001  // 更新玩家物品数据
	MONGO_GET_ACTOR_CHARACTER    = 40000  // 获取玩家角色数据
	MONGO_UPDATE_ACTOR_CHARACTER = 40001  // 更新玩家角色数据
	MONGO_GET_ACTOR_LOGIC        = 50000  // 获取玩家逻辑数据
	MONGO_UPDATE_ACTOR_LOGIC     = 50001  // 更新玩家逻辑数据
	MONGO_GET_ACTOR_TASK         = 60000  // 获取玩家任务数据
	MONGO_UPDATE_ACTOR_TASK      = 60001  // 更新玩家任务数据
	MONGO_GET_ACTOR_PVP          = 70000  // 获取玩家PVP数据
	MONGO_UPDATE_ACTOR_PVP       = 70001  // 更新玩家PVP数据
	MONGO_GET_ALL_ACTORS_PVP     = 70002  // 查找玩家PVP数据
	MONGO_GET_ACTOR_LEVEL        = 80000  // 获取玩家关卡数据
	MONGO_UPDATE_ACTOR_LEVEL     = 80001  // 更新玩家关卡数据
	MONGO_GET_ACTOR_PURCHASE     = 90000  // 获取玩家现金购买数据
	MONGO_UPDATE_ACTOR_PURCHASE  = 90001  // 更新玩家购买数据
	MONGO_GET_SERVER_DATA        = 91000  // 获取服务器数据
	MONGO_SET_SERVER_DATA        = 91001  // 设置服务器数据
	MONGO_GET_ALL_MAIL           = 100000 // 查找某个玩家的所有邮件
	MONGO_NEW_MAIL               = 100001 // 新建一封邮件
	MONGO_UPDATE_MAIL_STATUS     = 100002 // 更新一封邮件状态
	MONGO_DELETE_MAIL            = 100003 // 删除一封邮件
	MONGO_NEW_PURCHASE_ERROR     = 200000 // 新建一个购买错误记录
	// MONGO_REQUEST_ID_END
)

func init() {
	actorBaseProps = make([]string, xsf_pb.ActorProp_AP_Max)
	actorBaseProps[xsf_pb.ActorProp_AccountID] = "account_id"
	actorBaseProps[xsf_pb.ActorProp_ActorID] = "actor_id"
	actorBaseProps[xsf_pb.ActorProp_Nickname] = "nickname"
	actorBaseProps[xsf_pb.ActorProp_HeadIcon] = "head_icon"
	actorBaseProps[xsf_pb.ActorProp_CreateTime] = "create_time"
	actorBaseProps[xsf_pb.ActorProp_LastLogin] = "last_login"
	actorBaseProps[xsf_pb.ActorProp_SetNickTime] = "set_nick_time"
	actorBaseProps[xsf_pb.ActorProp_LogoutTime] = "logout_time"
	actorBaseProps[xsf_pb.ActorProp_Level] = "level"
	actorBaseProps[xsf_pb.ActorProp_Exp] = "exp"
	actorBaseProps[xsf_pb.ActorProp_LocalAccount] = "local_account"
}

func GetActorBaseProps() []string {
	return actorBaseProps
}

type MailData struct {
	Id       uint64 `bson:"id"`        // 唯一ID，或者sdk的token
	ActorID  uint32 `bson:"actor_id"`  // 角色ID
	SendTime uint32 `bson:"send_time"` // 创建时间
	Status   uint32 `bson:"status"`    // 邮件状态
	MailID   uint32 `bson:"mail_id"`
	MailData []byte `bson:"mail_data"`
	Title    string `bson:"title"`
	Body     string `bson:"body"`
}

func (md *MailData) Clone() *MailData {
	newData := &MailData{
		Id:       md.Id,
		ActorID:  md.ActorID,
		SendTime: md.SendTime,
		Status:   md.Status,
		MailID:   md.MailID,
		Title:    md.Title,
		Body:     md.Body,
	}

	newData.MailData = make([]byte, len(md.MailData))
	copy(newData.MailData, md.MailData)

	return newData
}
