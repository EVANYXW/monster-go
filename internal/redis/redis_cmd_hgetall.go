package redis

import (
	"context"
	"github.com/evanyxw/monster-go/internal/common"
	"github.com/evanyxw/monster-go/pkg/rpc"
	//xsf_rpc "xsf/rpc"

	"github.com/mediocregopher/radix/v4"
)

type redisCmd_HGETALL struct {
	conn radix.Conn
	ctx  context.Context
	rpc  *rpc.Acceptor

	rid uint32
	key string
}

func (rc *redisCmd_HGETALL) init(conn radix.Conn, ctx context.Context, rpc *rpc.Acceptor, keyID string, cfg *common.DBRedisCfg,
	// args[0] 为需要获取值的对象指针
	args []interface{}) {
	rc.conn = conn
	rc.ctx = ctx
	rc.rpc = rpc
	rc.rid = cfg.Id

	rc.key = getKey(keyID, cfg)
}

func (rc *redisCmd_HGETALL) do() {
	var err error
	var outValue map[string]string
	err = rc.conn.Do(rc.ctx, radix.Cmd(&outValue, "HGETALL", rc.key))
	if err == nil {
		//rc.rpc.Go(RPC_NAME_RESULT, rc.rid, REDIS_RES_OK, &outValue)
		return
	}

	rc.rpc.Go(RPC_NAME_RESULT, rc.rid, REDIS_RES_ERROR, err)
}
