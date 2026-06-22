package svc

import (
	"order-rpc/internal/config"
	"order-rpc/internal/repository"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config          config.Config
	Redis           *redis.Redis
	OrderRepository repository.OrderRepository
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	rds := redis.MustNewRedis(c.CacheRedis)
	return &ServiceContext{
		Config:          c,
		OrderRepository: repository.NewOrderRepository(conn, rds),
	}
}
