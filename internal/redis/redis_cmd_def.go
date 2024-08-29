package redis

import (
	"context"
	"github.com/evanyxw/monster-go/internal/common"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"github.com/evanyxw/monster-go/pkg/server"

	//xsf_server "xsf/server"

	"fmt"

	"github.com/mediocregopher/radix/v4"
)

type iRedisCmd interface {
	init(conn radix.Conn, ctx context.Context, rpc *rpc.Acceptor, keyID string, cfg *common.DBRedisCfg, args []interface{})
	do()
}

func getCmd(cmd string) iRedisCmd {
	switch cmd {
	case "SET":
		return new(redisCmd_SET)
	case "GET":
		return new(redisCmd_GET)
	case "HGET":
		return new(redisCmd_HGET)
	case "HSET":
		return new(redisCmd_HSET)
	case "HGETALL":
		return new(redisCmd_HGETALL)
	case "HMSET":
		return new(redisCmd_HMSET)
	case "HMSETBATCH":
		return new(redisCmd_HMSETBATCH)
	default:
		//xsf_log.Error("redis getCmd not support", xsf_log.String("cmd", cmd))
		return nil
	}
}

func getKey(keyID string, cfg *common.DBRedisCfg) string {
	tempKey := ""
	if len(cfg.FmtKey) > 0 {
		tempKey = fmt.Sprintf(cfg.FmtKey, keyID)
	} else {
		tempKey = cfg.Key
	}

	key := fmt.Sprintf("%v:%v", server.SID.ID, tempKey)
	//xsf_log.Info("redis get key=" + key)
	return key
}
