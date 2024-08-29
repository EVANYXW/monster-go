package redis

type DBRedisParam struct {
	Key  string
	Type string
}

const (
	rpc_NAME_REQUEST = "Request"
	RPC_NAME_RESULT  = "RedisResult"
)

const (
	REDIS_RES_OK    = 0
	REDIS_RES_ERROR = 1000

	REDIS_RES_ARG_ID    = 0
	REDIS_RES_ARG_CODE  = 1
	REDIS_RES_ARG_VALUE = 2
	REDIS_RES_ARG_ERROR = 2
)
