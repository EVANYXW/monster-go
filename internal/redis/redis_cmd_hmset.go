package redis

import (
	"context"
	"github.com/evanyxw/monster-go/internal/common"
	"github.com/evanyxw/monster-go/pkg/rpc"

	//xsf_rpc "xsf/rpc"

	"github.com/mediocregopher/radix/v4"
)

/*
// performs "HMSET" "foohash" "a" "1" "b" "2" "c" "3"
// m := map[string]int{"a": 1, "b": 2, "c": 3}

type ExampleStruct struct {
	Foo string // The key "Foo" will be used.
	Bar string `redis:"BAR"` // The key "BAR" will be used
	Baz string `redis:"-"`   // This field will be skipped
}

type OuterExampleStruct struct {
	// adds fields "Foo" and "BAR" to OuterExampleStruct
	ExampleStruct
	Biz int
}

s := OuterExampleStruct{
	ExampleStruct: ExampleStruct{
		Foo: "1",
		Bar: "2",
		Baz: "3",
	},
	Biz: 4,
}
*/

type redisCmd_HMSET struct {
	conn radix.Conn
	ctx  context.Context
	rpc  *rpc.Acceptor

	rid   uint32
	key   string
	value interface{}
}

func (rc *redisCmd_HMSET) init(conn radix.Conn, ctx context.Context, rpc *rpc.Acceptor, keyID string, cfg *common.DBRedisCfg,

	// args[0] 为需要设置的值, 可以是map，也可以是struct

	args []interface{}) {

	rc.conn = conn
	rc.ctx = ctx
	rc.rpc = rpc
	rc.rid = cfg.Id

	rc.key = getKey(keyID, cfg)
	rc.value = args[redis_ARG_START]
	//xsf_log.Debug(fmt.Sprintf("args:%v", args))
}

func (rc *redisCmd_HMSET) do() {
	//xsf_log.Debug(fmt.Sprintf("value:%v", rc.value))
	err := rc.conn.Do(rc.ctx, radix.FlatCmd(nil, "HMSET", rc.key, rc.value))
	if err != nil {
		rc.rpc.Go(RPC_NAME_RESULT, rc.rid, REDIS_RES_ERROR, err)
	} else {
		rc.rpc.Go(RPC_NAME_RESULT, rc.rid, REDIS_RES_OK)
	}
}
