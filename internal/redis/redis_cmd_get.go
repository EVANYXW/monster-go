package redis

import (
	"context"
	"fmt"
	"github.com/evanyxw/monster-go/internal/common"
	"github.com/evanyxw/monster-go/pkg/rpc"

	//xsf_log "xsf/log"
	//xsf_rpc "xsf/rpc"

	"github.com/mediocregopher/radix/v4"
)

type redisCmd_GET struct {
	conn radix.Conn
	ctx  context.Context
	rpc  *rpc.Acceptor

	rid     uint32
	key     string
	outType string
}

func (rc *redisCmd_GET) init(conn radix.Conn, ctx context.Context, rpc *rpc.Acceptor, keyID string, cfg *common.DBRedisCfg, args []interface{}) {
	rc.conn = conn
	rc.ctx = ctx
	rc.rpc = rpc
	rc.rid = cfg.Id

	rc.key = getKey(keyID, cfg)

	rc.outType = cfg.OutParams[0].Type
}

func (rc *redisCmd_GET) do() {
	var err error
	switch rc.outType {
	case "uint32":
		{
			var uOut uint32
			err = rc.conn.Do(rc.ctx, radix.FlatCmd(&uOut, "GET", rc.key))
			if err == nil {
				rc.rpc.Go(RPC_NAME_RESULT, rc.rid, REDIS_RES_OK, uOut)
				return
			}
		}

	case "float32":
		{
			var fOut float32
			err = rc.conn.Do(rc.ctx, radix.FlatCmd(&fOut, "GET", rc.key))
			if err == nil {
				rc.rpc.Go(RPC_NAME_RESULT, rc.rid, REDIS_RES_OK, fOut)
				return
			}
		}

	case "int32":
		{
			var iOut int32
			err = rc.conn.Do(rc.ctx, radix.FlatCmd(&iOut, "GET", rc.key))
			if err == nil {
				rc.rpc.Go(RPC_NAME_RESULT, rc.rid, REDIS_RES_OK, iOut)
				return
			}
		}

	case "string":
		{
			var sOut string
			err = rc.conn.Do(rc.ctx, radix.FlatCmd(&sOut, "GET", rc.key))
			if err == nil {
				rc.rpc.Go(RPC_NAME_RESULT, rc.rid, REDIS_RES_OK, sOut)
				return
			}
		}

	default:
		//xsf_log.Error("redisCmd_GET Do out type not support", xsf_log.String("type", rc.outType))
		err = fmt.Errorf("redisCmd_GET Do out type not support, type=%s", rc.outType)
	}

	rc.rpc.Go(RPC_NAME_RESULT, rc.rid, REDIS_RES_ERROR, err)
}
