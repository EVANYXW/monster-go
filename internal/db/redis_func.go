// ////////////////////////////////////////////////////////////////////////
//
// 文件：common/db/redis_func.go
// 作者：xiexx
// 时间：2023/07/17
// 描述：redis 相关函数
// 说明：
//
// ////////////////////////////////////////////////////////////////////////
package db

import (
	"github.com/evanyxw/monster-go/internal/redis"
	"github.com/evanyxw/monster-go/internal/servers/login/config"
	"github.com/evanyxw/monster-go/pkg/rpc"
	//xsf_log "xsf/log"
	//xsf_redis "xsf/redis"
)

// REQUEST_FUNC_BEGEIN
// 获取初始化标记
func RD_GetInitFlag(keyID uint32, keyStr string, rpc *rpc.Acceptor) {
	cfg := config.XSFSchema().DBRedis.Get(REDIS_GET_INIT_FLAG)
	if cfg == nil {
		//xsf_log.Error("RD_GetInitFlag cfg == nil REDIS_GET_INIT_FLAG")
		return
	}

	//xsf_redis.Request(keyID, keyStr, cfg, rpc)
}

// // 设置初始化标记
//
//	func RD_SetInitFlag(keyID uint32, keyStr string, rpc *xsf_rpc.Acceptor, flag uint32) {
//		cfg := xsf_scp.XSFSchema().DBRedis.Get(REDIS_SET_INIT_FLAG)
//		if cfg == nil {
//			xsf_log.Error("RD_SetInitFlag cfg == nil REDIS_SET_INIT_FLAG")
//			return
//		}
//
//		xsf_redis.Request(keyID, keyStr, cfg, rpc, flag)
//	}
//
// 获取账号信息
func RD_GetAccountInfo(keyID uint32, keyStr string, rpc *rpc.Acceptor) {
	cfg := config.XSFSchema().DBRedis.Get(REDIS_GET_ACCOUNT_INFO)
	if cfg == nil {
		//xsf_log.Error("RD_GetAccountInfo cfg == nil REDIS_GET_ACCOUNT_INFO")
		return
	}

	redis.Request(keyID, keyStr, cfg, rpc)
}

//
//// 设置账号信息
//func RD_SetAccountInfo(keyID uint32, keyStr string, rpc *xsf_rpc.Acceptor, account *AccountInfo) {
//	cfg := xsf_scp.XSFSchema().DBRedis.Get(REDIS_SET_ACCOUNT_INFO)
//	if cfg == nil {
//		xsf_log.Error("RD_SetAccountInfo cfg == nil REDIS_SET_ACCOUNT_INFO")
//		return
//	}
//
//	xsf_redis.Request(keyID, keyStr, cfg, rpc, account)
//}
//
//// 批量设置账号信息
//func RD_BatchSetAccountInfo(keyID uint32, keyStr string, rpc *xsf_rpc.Acceptor, account map[string]interface{}) {
//	cfg := xsf_scp.XSFSchema().DBRedis.Get(REDIS_BATCH_SET_ACCOUNT_INFO)
//	if cfg == nil {
//		xsf_log.Error("RD_BatchSetAccountInfo cfg == nil REDIS_BATCH_SET_ACCOUNT_INFO")
//		return
//	}
//
//	xsf_redis.Request(keyID, keyStr, cfg, rpc, account)
//}

// REQUEST_FUNC_END
