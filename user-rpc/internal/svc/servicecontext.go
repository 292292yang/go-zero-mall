package svc

import (
	"github.com/go-zero-mall/user-rpc/internal/config"
	"github.com/go-zero-mall/user-rpc/internal/repository"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config         config.Config
	Redis          *redis.Redis
	UserRepository repository.UserRepository
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	rds := redis.MustNewRedis(c.CacheRedis)
	return &ServiceContext{
		Config:         c,
		UserRepository: repository.NewUserRepository(conn, rds),
	}
}
