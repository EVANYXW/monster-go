// ////////////////////////////////////////////////////////////////////////
//
// 文件：common/db/mongo_func.go
// 作者：xiexx
// 时间：2023/07/17
// 描述：mongo 相关函数
// 说明：
//
// ////////////////////////////////////////////////////////////////////////
package db

import (
	"github.com/evanyxw/monster-go/internal/servers/login/config"
	"github.com/evanyxw/monster-go/pkg/rpc"
	//"go.mongodb.org/mongo-driver/bson"
)

// REQUEST_FUNC_BEGEIN
// 查找所有账号数据
func MG_GetAllAccount(keyID uint32, rpc *rpc.Acceptor, page uint32) {
	cfg := config.XSFSchema().DBMongo.Get(MONGO_GET_ALL_ACCOUNT)
	if cfg == nil {
		//xsf_log.Error("MG_GetAllAccount cfg == nil MONGO_GET_ALL_ACCOUNT")
		return
	}

	//filter := bson.M{}

	//xsf_mongo.FindAll(keyID, cfg, rpc, page, filter)
}

//
//// 查找最大角色id
//func MG_GetMaxActorId(keyID uint32, rpc *rpc.Acceptor) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_GET_MAX_ACTOR_ID)
//	if cfg == nil {
//		xsf_log.Error("MG_GetMaxActorId cfg == nil MONGO_GET_MAX_ACTOR_ID")
//		return
//	}
//
//	filter := bson.M{}
//
//	xsf_mongo.FindOneSort(keyID, cfg, rpc, filter)
//}
//
//// 查找账号
//func MG_GetAccount(keyID uint32, rpc *rpc.Acceptor, F_id string) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_GET_ACCOUNT)
//	if cfg == nil {
//		xsf_log.Error("MG_GetAccount cfg == nil MONGO_GET_ACCOUNT")
//		return
//	}
//
//	filter := bson.M{
//		cfg.Filters[0]: F_id,
//	}
//
//	xsf_mongo.FindOne(keyID, cfg, rpc, filter)
//}
//
//// 新建账号
//func MG_NewAccount(keyID uint32, rpc *rpc.Acceptor, D_id string, D_actor_id uint32, D_time uint32, D_game_id uint32) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_NEW_ACCOUNT)
//	if cfg == nil {
//		xsf_log.Error("MG_NewAccount cfg == nil MONGO_NEW_ACCOUNT")
//		return
//	}
//
//	data := bson.M{
//		cfg.Datas[0]: D_id,
//		cfg.Datas[1]: D_actor_id,
//		cfg.Datas[2]: D_time,
//		cfg.Datas[3]: D_game_id,
//	}
//
//	xsf_mongo.Insert(keyID, cfg, rpc, data)
//}
//
//// 获取玩家基础数据
//func MG_GetActorBase(keyID uint32, rpc *rpc.Acceptor, F_actor_id uint32) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_GET_ACTOR_BASE)
//	if cfg == nil {
//		xsf_log.Error("MG_GetActorBase cfg == nil MONGO_GET_ACTOR_BASE")
//		return
//	}
//
//	filter := bson.M{
//		cfg.Filters[0]: F_actor_id,
//	}
//
//	xsf_mongo.FindOne(keyID, cfg, rpc, filter)
//}
//
//// 更新玩家数据
//func MG_UpdateActorBase(keyID uint32, rpc *rpc.Acceptor, F_actor_id uint32, D_actor_base map[string]interface{}, IsUpsert bool) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_UPDATE_ACTOR_BASE)
//	if cfg == nil {
//		xsf_log.Error("MG_UpdateActorBase cfg == nil MONGO_UPDATE_ACTOR_BASE")
//		return
//	}
//
//	filter := bson.M{
//		cfg.Filters[0]: F_actor_id,
//	}
//
//	xsf_mongo.UpdateOne(keyID, cfg, rpc, IsUpsert, filter, D_actor_base)
//}
//
//// 查找所有玩家数据
//func MG_GetAllActors(keyID uint32, rpc *rpc.Acceptor, page uint32) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_GET_ALL_ACTORS)
//	if cfg == nil {
//		xsf_log.Error("MG_GetAllActors cfg == nil MONGO_GET_ALL_ACTORS")
//		return
//	}
//
//	filter := bson.M{}
//
//	xsf_mongo.FindAll(keyID, cfg, rpc, page, filter)
//}
//
//// 获取玩家物品数据
//func MG_GetActorItem(keyID uint32, rpc *rpc.Acceptor, F_actor_id uint32) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_GET_ACTOR_ITEM)
//	if cfg == nil {
//		xsf_log.Error("MG_GetActorItem cfg == nil MONGO_GET_ACTOR_ITEM")
//		return
//	}
//
//	filter := bson.M{
//		cfg.Filters[0]: F_actor_id,
//	}
//
//	xsf_mongo.FindOne(keyID, cfg, rpc, filter)
//}
//
//// 更新玩家物品数据
//func MG_UpdateActorItem(keyID uint32, rpc *rpc.Acceptor, F_actor_id uint32, D_actor_id uint32, D_data []byte, D_update_time uint32, IsUpsert bool) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_UPDATE_ACTOR_ITEM)
//	if cfg == nil {
//		xsf_log.Error("MG_UpdateActorItem cfg == nil MONGO_UPDATE_ACTOR_ITEM")
//		return
//	}
//
//	filter := bson.M{
//		cfg.Filters[0]: F_actor_id,
//	}
//
//	data := bson.M{
//		cfg.Datas[0]: D_actor_id,
//		cfg.Datas[1]: D_data,
//		cfg.Datas[2]: D_update_time,
//	}
//	xsf_mongo.UpdateOne(keyID, cfg, rpc, IsUpsert, filter, data)
//}
//
//// 获取玩家角色数据
//func MG_GetActorCharacter(keyID uint32, rpc *rpc.Acceptor, F_actor_id uint32) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_GET_ACTOR_CHARACTER)
//	if cfg == nil {
//		xsf_log.Error("MG_GetActorCharacter cfg == nil MONGO_GET_ACTOR_CHARACTER")
//		return
//	}
//
//	filter := bson.M{
//		cfg.Filters[0]: F_actor_id,
//	}
//
//	xsf_mongo.FindOne(keyID, cfg, rpc, filter)
//}
//
//// 更新玩家角色数据
//func MG_UpdateActorCharacter(keyID uint32, rpc *rpc.Acceptor, F_actor_id uint32, D_actor_id uint32, D_data []byte, D_update_time uint32, IsUpsert bool) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_UPDATE_ACTOR_CHARACTER)
//	if cfg == nil {
//		xsf_log.Error("MG_UpdateActorCharacter cfg == nil MONGO_UPDATE_ACTOR_CHARACTER")
//		return
//	}
//
//	filter := bson.M{
//		cfg.Filters[0]: F_actor_id,
//	}
//
//	data := bson.M{
//		cfg.Datas[0]: D_actor_id,
//		cfg.Datas[1]: D_data,
//		cfg.Datas[2]: D_update_time,
//	}
//	xsf_mongo.UpdateOne(keyID, cfg, rpc, IsUpsert, filter, data)
//}
//
//// 获取玩家逻辑数据
//func MG_GetActorLogic(keyID uint32, rpc *rpc.Acceptor, F_actor_id uint32) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_GET_ACTOR_LOGIC)
//	if cfg == nil {
//		xsf_log.Error("MG_GetActorLogic cfg == nil MONGO_GET_ACTOR_LOGIC")
//		return
//	}
//
//	filter := bson.M{
//		cfg.Filters[0]: F_actor_id,
//	}
//
//	xsf_mongo.FindOne(keyID, cfg, rpc, filter)
//}
//
//// 更新玩家逻辑数据
//func MG_UpdateActorLogic(keyID uint32, rpc *rpc.Acceptor, F_actor_id uint32, D_actor_id uint32, D_data []byte, D_update_time uint32, IsUpsert bool) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_UPDATE_ACTOR_LOGIC)
//	if cfg == nil {
//		xsf_log.Error("MG_UpdateActorLogic cfg == nil MONGO_UPDATE_ACTOR_LOGIC")
//		return
//	}
//
//	filter := bson.M{
//		cfg.Filters[0]: F_actor_id,
//	}
//
//	data := bson.M{
//		cfg.Datas[0]: D_actor_id,
//		cfg.Datas[1]: D_data,
//		cfg.Datas[2]: D_update_time,
//	}
//	xsf_mongo.UpdateOne(keyID, cfg, rpc, IsUpsert, filter, data)
//}
//
//// 获取玩家任务数据
//func MG_GetActorTask(keyID uint32, rpc *rpc.Acceptor, F_actor_id uint32) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_GET_ACTOR_TASK)
//	if cfg == nil {
//		xsf_log.Error("MG_GetActorTask cfg == nil MONGO_GET_ACTOR_TASK")
//		return
//	}
//
//	filter := bson.M{
//		cfg.Filters[0]: F_actor_id,
//	}
//
//	xsf_mongo.FindOne(keyID, cfg, rpc, filter)
//}
//
//// 更新玩家任务数据
//func MG_UpdateActorTask(keyID uint32, rpc *rpc.Acceptor, F_actor_id uint32, D_actor_id uint32, D_data []byte, D_update_time uint32, IsUpsert bool) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_UPDATE_ACTOR_TASK)
//	if cfg == nil {
//		xsf_log.Error("MG_UpdateActorTask cfg == nil MONGO_UPDATE_ACTOR_TASK")
//		return
//	}
//
//	filter := bson.M{
//		cfg.Filters[0]: F_actor_id,
//	}
//
//	data := bson.M{
//		cfg.Datas[0]: D_actor_id,
//		cfg.Datas[1]: D_data,
//		cfg.Datas[2]: D_update_time,
//	}
//	xsf_mongo.UpdateOne(keyID, cfg, rpc, IsUpsert, filter, data)
//}
//
//// 获取玩家PVP数据
//func MG_GetActorPvp(keyID uint32, rpc *rpc.Acceptor, F_actor_id uint32) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_GET_ACTOR_PVP)
//	if cfg == nil {
//		xsf_log.Error("MG_GetActorPvp cfg == nil MONGO_GET_ACTOR_PVP")
//		return
//	}
//
//	filter := bson.M{
//		cfg.Filters[0]: F_actor_id,
//	}
//
//	xsf_mongo.FindOne(keyID, cfg, rpc, filter)
//}
//
//// 更新玩家PVP数据
//func MG_UpdateActorPvp(keyID uint32, rpc *rpc.Acceptor, F_actor_id uint32, D_actor_id uint32, D_data []byte, D_update_time uint32, IsUpsert bool) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_UPDATE_ACTOR_PVP)
//	if cfg == nil {
//		xsf_log.Error("MG_UpdateActorPvp cfg == nil MONGO_UPDATE_ACTOR_PVP")
//		return
//	}
//
//	filter := bson.M{
//		cfg.Filters[0]: F_actor_id,
//	}
//
//	data := bson.M{
//		cfg.Datas[0]: D_actor_id,
//		cfg.Datas[1]: D_data,
//		cfg.Datas[2]: D_update_time,
//	}
//	xsf_mongo.UpdateOne(keyID, cfg, rpc, IsUpsert, filter, data)
//}
//
//// 查找玩家PVP数据
//func MG_GetAllActorsPvp(keyID uint32, rpc *rpc.Acceptor, page uint32) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_GET_ALL_ACTORS_PVP)
//	if cfg == nil {
//		xsf_log.Error("MG_GetAllActorsPvp cfg == nil MONGO_GET_ALL_ACTORS_PVP")
//		return
//	}
//
//	filter := bson.M{}
//
//	xsf_mongo.FindAll(keyID, cfg, rpc, page, filter)
//}
//
//// 获取玩家关卡数据
//func MG_GetActorLevel(keyID uint32, rpc *rpc.Acceptor, F_actor_id uint32) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_GET_ACTOR_LEVEL)
//	if cfg == nil {
//		xsf_log.Error("MG_GetActorLevel cfg == nil MONGO_GET_ACTOR_LEVEL")
//		return
//	}
//
//	filter := bson.M{
//		cfg.Filters[0]: F_actor_id,
//	}
//
//	xsf_mongo.FindOne(keyID, cfg, rpc, filter)
//}
//
//// 更新玩家关卡数据
//func MG_UpdateActorLevel(keyID uint32, rpc *rpc.Acceptor, F_actor_id uint32, D_actor_id uint32, D_data []byte, D_update_time uint32, IsUpsert bool) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_UPDATE_ACTOR_LEVEL)
//	if cfg == nil {
//		xsf_log.Error("MG_UpdateActorLevel cfg == nil MONGO_UPDATE_ACTOR_LEVEL")
//		return
//	}
//
//	filter := bson.M{
//		cfg.Filters[0]: F_actor_id,
//	}
//
//	data := bson.M{
//		cfg.Datas[0]: D_actor_id,
//		cfg.Datas[1]: D_data,
//		cfg.Datas[2]: D_update_time,
//	}
//	xsf_mongo.UpdateOne(keyID, cfg, rpc, IsUpsert, filter, data)
//}
//
//// 获取玩家现金购买数据
//func MG_GetActorPurchase(keyID uint32, rpc *rpc.Acceptor, F_actor_id uint32) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_GET_ACTOR_PURCHASE)
//	if cfg == nil {
//		xsf_log.Error("MG_GetActorPurchase cfg == nil MONGO_GET_ACTOR_PURCHASE")
//		return
//	}
//
//	filter := bson.M{
//		cfg.Filters[0]: F_actor_id,
//	}
//
//	xsf_mongo.FindOne(keyID, cfg, rpc, filter)
//}
//
//// 更新玩家购买数据
//func MG_UpdateActorPurchase(keyID uint32, rpc *rpc.Acceptor, F_actor_id uint32, D_actor_id uint32, D_data []byte, D_update_time uint32, IsUpsert bool) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_UPDATE_ACTOR_PURCHASE)
//	if cfg == nil {
//		xsf_log.Error("MG_UpdateActorPurchase cfg == nil MONGO_UPDATE_ACTOR_PURCHASE")
//		return
//	}
//
//	filter := bson.M{
//		cfg.Filters[0]: F_actor_id,
//	}
//
//	data := bson.M{
//		cfg.Datas[0]: D_actor_id,
//		cfg.Datas[1]: D_data,
//		cfg.Datas[2]: D_update_time,
//	}
//	xsf_mongo.UpdateOne(keyID, cfg, rpc, IsUpsert, filter, data)
//}
//
//// 获取服务器数据
//func MG_GetServerData(keyID uint32, rpc *rpc.Acceptor, F_id uint32) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_GET_SERVER_DATA)
//	if cfg == nil {
//		xsf_log.Error("MG_GetServerData cfg == nil MONGO_GET_SERVER_DATA")
//		return
//	}
//
//	filter := bson.M{
//		cfg.Filters[0]: F_id,
//	}
//
//	xsf_mongo.FindOne(keyID, cfg, rpc, filter)
//}
//
//// 设置服务器数据
//func MG_SetServerData(keyID uint32, rpc *rpc.Acceptor, F_id uint32, D_id uint32, D_start_time uint32, IsUpsert bool) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_SET_SERVER_DATA)
//	if cfg == nil {
//		xsf_log.Error("MG_SetServerData cfg == nil MONGO_SET_SERVER_DATA")
//		return
//	}
//
//	filter := bson.M{
//		cfg.Filters[0]: F_id,
//	}
//
//	data := bson.M{
//		cfg.Datas[0]: D_id,
//		cfg.Datas[1]: D_start_time,
//	}
//	xsf_mongo.UpdateOne(keyID, cfg, rpc, IsUpsert, filter, data)
//}
//
//// 查找某个玩家的所有邮件
//func MG_GetAllMail(keyID uint32, rpc *rpc.Acceptor, page uint32, F_actor_id uint32, F_send_time map[string]interface{}) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_GET_ALL_MAIL)
//	if cfg == nil {
//		xsf_log.Error("MG_GetAllMail cfg == nil MONGO_GET_ALL_MAIL")
//		return
//	}
//
//	filter := bson.M{
//		cfg.Filters[0]: F_actor_id,
//		cfg.Filters[1]: F_send_time,
//	}
//
//	xsf_mongo.FindAll(keyID, cfg, rpc, page, filter)
//}
//
//// 新建一封邮件
//func MG_NewMail(keyID uint32, rpc *rpc.Acceptor, D_id uint64, D_actor_id uint32, D_send_time uint32, D_status uint32, D_mail_id uint32, D_mail_data []byte, D_title string, D_body string) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_NEW_MAIL)
//	if cfg == nil {
//		xsf_log.Error("MG_NewMail cfg == nil MONGO_NEW_MAIL")
//		return
//	}
//
//	data := bson.M{
//		cfg.Datas[0]: D_id,
//		cfg.Datas[1]: D_actor_id,
//		cfg.Datas[2]: D_send_time,
//		cfg.Datas[3]: D_status,
//		cfg.Datas[4]: D_mail_id,
//		cfg.Datas[5]: D_mail_data,
//		cfg.Datas[6]: D_title,
//		cfg.Datas[7]: D_body,
//	}
//
//	xsf_mongo.Insert(keyID, cfg, rpc, data)
//}
//
//// 更新一封邮件状态
//func MG_UpdateMailStatus(keyID uint32, rpc *rpc.Acceptor, F_id uint64, D_status uint32, IsUpsert bool) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_UPDATE_MAIL_STATUS)
//	if cfg == nil {
//		xsf_log.Error("MG_UpdateMailStatus cfg == nil MONGO_UPDATE_MAIL_STATUS")
//		return
//	}
//
//	filter := bson.M{
//		cfg.Filters[0]: F_id,
//	}
//
//	data := bson.M{
//		cfg.Datas[0]: D_status,
//	}
//	xsf_mongo.UpdateOne(keyID, cfg, rpc, IsUpsert, filter, data)
//}
//
//// 删除一封邮件
//func MG_DeleteMail(keyID uint32, rpc *rpc.Acceptor, F_id uint64) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_DELETE_MAIL)
//	if cfg == nil {
//		xsf_log.Error("MG_DeleteMail cfg == nil MONGO_DELETE_MAIL")
//		return
//	}
//
//	filter := bson.M{
//		cfg.Filters[0]: F_id,
//	}
//
//	xsf_mongo.DeleteOne(keyID, cfg, rpc, filter)
//}
//
//// 新建一个购买错误记录
//func MG_NewPurchaseError(keyID uint32, rpc *rpc.Acceptor, D_actor_id uint32, D_error_code uint32, D_store_code uint32, D_error_msg string, D_product_id string, D_is_offline bool) {
//	cfg := xsf_scp.XSFSchema().DBMongo.Get(MONGO_NEW_PURCHASE_ERROR)
//	if cfg == nil {
//		xsf_log.Error("MG_NewPurchaseError cfg == nil MONGO_NEW_PURCHASE_ERROR")
//		return
//	}
//
//	data := bson.M{
//		cfg.Datas[0]: D_actor_id,
//		cfg.Datas[1]: D_error_code,
//		cfg.Datas[2]: D_store_code,
//		cfg.Datas[3]: D_error_msg,
//		cfg.Datas[4]: D_product_id,
//		cfg.Datas[5]: D_is_offline,
//	}
//
//	xsf_mongo.Insert(keyID, cfg, rpc, data)
//}

// REQUEST_FUNC_END
