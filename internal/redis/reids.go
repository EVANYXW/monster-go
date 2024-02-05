package redis

import (
	"net"
	"sync"

	"github.com/go-redis/redis"
)

var RedisManagers *RedisManager

type SubFun func(channel, data string)

type RedisConfig struct {
	Addr      []string
	Passwd    string
	PoolSize  int
	IsCluster bool
}

type Redis struct {
	redis.Cmdable
	pubsub   *redis.PubSub
	conf     *RedisConfig
	manager  *RedisManager
	fun      func(channel, data string)
	channels []string
}

type RedisManager struct {
	db   *Redis
	lock sync.RWMutex
}

func (r *RedisManager) AddConf(conf *RedisConfig) {
	r.lock.Lock()
	defer r.lock.Unlock()

	var redisCmdable redis.Cmdable
	if conf.IsCluster {
		redisCmdable = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    conf.Addr,
			Password: conf.Passwd,
			PoolSize: conf.PoolSize,
		})
	} else {
		redisCmdable = redis.NewClient(&redis.Options{
			Addr:     conf.Addr[0],
			Password: conf.Passwd,
			PoolSize: conf.PoolSize,
		})
	}

	re := &Redis{
		Cmdable: redisCmdable,
		conf:    conf,
		manager: r,
	}
	r.db = re
	//r.db
	//LogInfo("connect to redis %v", conf.Addr)
}

func NewRedisManager(conf *RedisConfig) *RedisManager {
	redisManager := &RedisManager{}
	redisManager.AddConf(conf)
	return redisManager
}

func (r *RedisManager) GetDB() *Redis {
	return r.db
}

func (r *RedisManager) Close() {
	r.db.pubsub.Close()
}

// 用redis做pub/sub,要注意缓存限制
// 客户端订阅了channel之后，如果接收消息不及时，可能导致DCS实例消息堆积，当达到消息堆积阈值（默认值为32MB），
// 或者达到某种程度（默认8MB）一段时间（默认为1分钟）后，服务器端会自动断开该客户端连接，避免导致内部内存耗尽.
// https://support.huaweicloud.com/dcs_faq/dcs-faq-0427017.html
func (r *RedisManager) Sub(fun SubFun, channels ...string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.db.channels = channels
	r.db.fun = fun

	redisConf := r.db.conf
	var pubsub *redis.PubSub
	if redisConf.IsCluster {
		db := r.db.Cmdable.(*redis.ClusterClient)
		pubsub = db.Subscribe(channels...)
	} else {
		db := r.db.Cmdable.(*redis.Client)
		pubsub = db.Subscribe(channels...)
	}

	r.db.pubsub = pubsub
	go (func() {
		for {
			msg, err := pubsub.ReceiveMessage()
			if err == nil {
				go (func() { fun(msg.Channel, msg.Payload) })()
			} else if _, ok := err.(net.Error); !ok {
				break
			}
		}
	})()

}
