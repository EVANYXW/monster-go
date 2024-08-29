package redis

import (
	"context"
	"github.com/evanyxw/monster-go/internal/common"
	"github.com/evanyxw/monster-go/pkg/rpc"

	//xsf_rpc "xsf/rpc"

	"github.com/mediocregopher/radix/v4"
)

type redisCmd_HSET struct {
	conn radix.Conn
	ctx  context.Context
	rpc  *rpc.Acceptor

	rid    uint32
	key    string
	subKey string
	value  interface{}
}

func (rc *redisCmd_HSET) init(conn radix.Conn, ctx context.Context, rpc *rpc.Acceptor, keyID string, cfg *common.DBRedisCfg,

	// args[0] 为需要设置的值
	args []interface{}) {

	rc.conn = conn
	rc.ctx = ctx
	rc.rpc = rpc
	rc.rid = cfg.Id

	paramIndex := redis_ARG_START

	rc.key = getKey(keyID, cfg)
	rc.subKey = cfg.InParams[0].Key
	rc.value = args[paramIndex]
}

func (rc *redisCmd_HSET) do() {
	err := rc.conn.Do(rc.ctx, radix.FlatCmd(nil, "HSET", rc.key, rc.subKey, rc.value))
	if err != nil {
		rc.rpc.Go(RPC_NAME_RESULT, rc.rid, REDIS_RES_ERROR, err)
	} else {
		rc.rpc.Go(RPC_NAME_RESULT, rc.rid, REDIS_RES_OK)
	}
}
