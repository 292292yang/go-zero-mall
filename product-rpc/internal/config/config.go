package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql      MysqlConf
	CacheRedis redis.RedisConf
	JwtAuth    JwtAuthConf
}

type MysqlConf struct {
	DataSource string
}

type JwtAuthConf struct {
	AccessSecret string
	AccessExpire int64
}
