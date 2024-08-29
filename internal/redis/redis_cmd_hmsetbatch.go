package redis

import (
	"context"
	"fmt"
	"github.com/evanyxw/monster-go/internal/common"
	"github.com/evanyxw/monster-go/pkg/rpc"
	//xsf_rpc "xsf/rpc"

	"github.com/mediocregopher/radix/v4"
)

type redisCmd_HMSETBATCH struct {
	conn radix.Conn
	ctx  context.Context
	rpc  *rpc.Acceptor

	rid    uint32
	key    string
	fmtKey string
	value  interface{}
}

func (rc *redisCmd_HMSETBATCH) init(conn radix.Conn, ctx context.Context, rpc *rpc.Acceptor, keyID string, cfg *common.DBRedisCfg,

	// args[0] 为需要设置的值
	args []interface{}) {

	rc.conn = conn
	rc.ctx = ctx
	rc.rpc = rpc
	rc.rid = cfg.Id
	rc.fmtKey = cfg.FmtKey
	paramIndex := redis_ARG_START

	rc.key = getKey(keyID, cfg)

	rc.value = args[paramIndex].([]interface{})[0]
}

func (rc *redisCmd_HMSETBATCH) do() {
	pl := radix.NewPipeline()

	for key, value := range rc.value.(map[string]interface{}) {
		rk := getRedisKey(key, rc.fmtKey)
		pl.Append(radix.FlatCmd(nil, "HMSET", rk, value))
	}
	err := rc.conn.Do(rc.ctx, pl)
	if err != nil {
		rc.rpc.Go(RPC_NAME_RESULT, rc.rid, REDIS_RES_ERROR, err)
	} else {
		rc.rpc.Go(RPC_NAME_RESULT, rc.rid, REDIS_RES_OK)
	}
}

func getRedisKey(keyId string, fmtKey string) string {
	return fmt.Sprintf(fmtKey, keyId)
}
