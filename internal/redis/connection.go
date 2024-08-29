package redis

import (
	"context"
	"fmt"
	"github.com/evanyxw/monster-go/internal/common"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/evanyxw/monster-go/pkg/rpc"
	"sync"

	"github.com/mediocregopher/radix/v4"
	//xsf_log "xsf/log"
)

const (
	redis_ARG_KEY_ID = 0 // 唯一ID
	redis_ARG_SCP    = 1
	redis_ARG_RPC    = 2
	redis_ARG_START  = 3
)

type SCDBInfo struct {
	ID    uint32
	Desc  string
	IP    string
	User  string
	Pwd   string
	Port  uint32
	Count int
}

type redisConn struct {
	closeSig chan bool
	wg       sync.WaitGroup

	acceptor *rpc.Acceptor

	conn   radix.Conn
	ctx    context.Context
	cancel context.CancelFunc
}

func (rc *redisConn) init(info *SCDBInfo) bool {
	rc.closeSig = make(chan bool, 1)
	rc.acceptor = rpc.NewAcceptor(10000)
	rc.acceptor.Regist(rpc_NAME_REQUEST, rc.onRequest)

	rc.ctx, rc.cancel = context.WithCancel(context.Background())

	var dialer radix.Dialer
	dialer.AuthPass = info.Pwd
	url := fmt.Sprintf("%s:%d", info.IP, info.Port)
	conn, err := dialer.Dial(rc.ctx, "tcp", url)
	if err != nil {
		//xsf_log.Panic("redisConn Init error", xsf_log.NamedError("error", err))
		return false
	}

	rc.conn = conn

	return true
}

func (rc *redisConn) close() {
	rc.closeSig <- true
}

func (rc *redisConn) release() {
	rc.cancel()
	rc.acceptor.Close()
	rc.wg.Wait()
}

func (rc *redisConn) start() {
	async.Go(func() {
	OUTLABEL:
		for {
			select {
			case <-rc.closeSig:
				return

			case callMsg, ok := <-rc.acceptor.ChanCall:
				if !ok {
					break OUTLABEL
					//return
				}
				if callMsg == nil {
					continue
				}
				rc.acceptor.Execute(callMsg)
			}
		}
	})

}

func (rc *redisConn) request(keyID string, cfg *common.DBRedisCfg, rpc *rpc.Acceptor, args interface{}) {
	rc.acceptor.Go(rpc_NAME_REQUEST, keyID, cfg, rpc, args)
}

func (rc *redisConn) onRequest(args []interface{}) {
	keyID := args[redis_ARG_KEY_ID].(string)
	cfg := args[redis_ARG_SCP].(*common.DBRedisCfg)
	rpc := args[redis_ARG_RPC].(*rpc.Acceptor)

	cmd := getCmd(cfg.Op)
	if cmd == nil {
		//xsf_log.Error("RedisConn onRequest cmd == nil", xsf_log.String("cmd", cfg.Op))
		return
	}

	cmd.init(rc.conn, rc.ctx, rpc, keyID, cfg, args)
	cmd.do()
}
