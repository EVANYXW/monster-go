package redis

import (
	"github.com/evanyxw/monster-go/internal/common"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
	"github.com/evanyxw/monster-go/pkg/rpc"
)

type RedisKernel struct {
	msgHandler module.Handler
	conns      []*redisConn
	total      int
	id         uint32
}

var (
	//redisModule *xsf_module.Module
	redisKernl *RedisKernel
)

func NewRedisKernel(msgHandler module.Handler, rds []SCDBInfo) *RedisKernel {
	rk := &RedisKernel{
		msgHandler: msgHandler,
	}
	rk.total = len(rds)
	rk.conns = make([]*redisConn, rk.total)
	for i := 0; i < rk.total; i++ {
		config := &(rds[i])
		//if rk.id == config.ID {
		for J := 0; J < int(config.Count); J++ {
			rk.conns[i] = new(redisConn)
			rk.conns[i].init(config)
		}
		//}
	}
	redisKernl = rk
	return rk
}

func (rk *RedisKernel) Init(baseModule *module.BaseModule) bool {
	rk.msgHandler.OnInit(baseModule)
	return true
}

func (rk *RedisKernel) DoRegister() {

}

func (rk *RedisKernel) GetNPManager() network.INPManager {
	return nil
}

func (rk *RedisKernel) GetStatus() int {
	return 0
}

func (rk *RedisKernel) DoRun() {
	//rk.msgHandler.Start()
	for i := 0; i < rk.total; i++ {
		rk.conns[i].start()
	}
}

func (rk *RedisKernel) DoWaitStart() {

}

func (rk *RedisKernel) DoRelease() {
	for i := 0; i < rk.total; i++ {
		rk.conns[i].release()
	}
}

func (rk *RedisKernel) Update() {

}

func (rk *RedisKernel) OnOk() {
	//rk.msgHandler.OnOk()
}

func (rk *RedisKernel) OnStartClose() {

}

func (rk *RedisKernel) DoClose() {

	for i := 0; i < rk.total; i++ {
		rk.conns[i].close()
	}
}

func (rk *RedisKernel) OnStartCheck() int {
	return 0
}

func (rk *RedisKernel) OnCloseCheck() int {
	return 0
}

func (rk *RedisKernel) GetNoWaitStart() bool {
	return true
}

func (rk *RedisKernel) MessageHandler(packet *network.Packet) {

}

func (rk *RedisKernel) OnRpcNetAccept(args []interface{}) {

}

func (rk *RedisKernel) OnRpcNetConnected(args []interface{}) {

}

func (rk *RedisKernel) OnRpcNetError(args []interface{}) {

}

func (rk *RedisKernel) OnRpcNetData(args []interface{}) {

}

func (rk *RedisKernel) OnRpcNetMessage(args []interface{}) {

}

func Request(keyID uint32, keyStr string, cfg *common.DBRedisCfg, rpc *rpc.Acceptor, args ...interface{}) {
	index := keyID % uint32(redisKernl.total)

	redisKernl.conns[index].request(keyStr, cfg, rpc, args)
}
