package module

import (
	"github.com/evanyxw/monster-go/internal/redis"
	"github.com/evanyxw/monster-go/internal/servers/common/handler"
	"github.com/evanyxw/monster-go/pkg/logger"
	"github.com/evanyxw/monster-go/pkg/module"
	"github.com/evanyxw/monster-go/pkg/network"
)

// 客户端消息接受体

type RedisClient struct {
	kernel module.IModuleKernel
	id     int32
}

func NewRedisClient(id int32) *RedisClient {
	var redisInfo []redis.SCDBInfo
	redisInfo = []redis.SCDBInfo{
		{
			ID:    1,
			Desc:  "monster",
			IP:    "127.0.0.1",
			User:  "root",
			Pwd:   "",
			Port:  6379,
			Count: 1,
		},
	}
	h := handler.NewCommonMsgHandler()
	r := &RedisClient{
		id:     id,
		kernel: redis.NewRedisKernel(h, redisInfo),
	}
	//module.NewBaseModule(id, r)
	//h.Init(baseModule) //fixMe 这个看能否改为kernel 里去调用
	return r
}

func (r *RedisClient) Init(baseModule *module.BaseModule) bool {
	r.kernel.Init(baseModule)
	return true
}

// DoRun BaseModule 调用
func (r *RedisClient) DoRun() {
	r.kernel.DoRun()
}

func (r *RedisClient) DoWaitStart() {
	r.kernel.DoWaitStart()
}

func (r *RedisClient) DoRelease() {
	r.kernel.DoRelease()
}

func (r *RedisClient) OnOk() {

}

func (r *RedisClient) OnStartCheck() int {
	return module.ModuleRunCode_Ok
}

func (r *RedisClient) OnCloseCheck() int {
	return r.kernel.OnCloseCheck()
}

func (r *RedisClient) GetID() int32 {
	return r.id
}

func (r *RedisClient) GetKernel() module.IModuleKernel {
	return r.kernel
}

func (r *RedisClient) Update() {

}

func (r *RedisClient) DoRegister() {
	r.kernel.DoRegister()
}

func (r *RedisClient) OnNetError(np *network.NetPoint) {
	logger.Debug("center onNetError")
	//r.nodeManager.OnNodeLost(np.ID, np.SID.Type)
	module.NodeManager.OnNodeLost(np.ID, np.SID.Type)
}

func (r *RedisClient) OnServerOk() {

}

func (r *RedisClient) OnNPAdd(np *network.NetPoint) {

}
