package svc

import (
	"github.com/292292yang/go-zero-mall/product-rpc/internal/config"
	"github.com/292292yang/go-zero-mall/product-rpc/internal/repository"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config            config.Config
	Redis             *redis.Redis
	ProductRepository repository.ProductRepository
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	rds := redis.MustNewRedis(c.CacheRedis)
	return &ServiceContext{
		Config:            c,
		ProductRepository: repository.NewProductRepository(conn, rds),
	}
}
