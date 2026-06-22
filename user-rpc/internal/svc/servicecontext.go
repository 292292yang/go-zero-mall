package svc

import (
	"github.com/292292yang/go-zero-mall/user-rpc/internal/config"
	"github.com/292292yang/go-zero-mall/user-rpc/internal/repository"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc"
)

type ServiceContext struct {
	Config         config.Config
	Redis          *redis.Redis
	RpcServerOpts  []grpc.ServerOption
	UserRepository repository.UserRepository
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	rds := redis.MustNewRedis(c.CacheRedis)
	return &ServiceContext{
		Config: c,
		// 注册全局拦截器
		RpcServerOpts: []grpc.ServerOption{
			grpc.UnaryInterceptor(RpcErrorInterceptor),
		},
		UserRepository: repository.NewUserRepository(conn, rds),
	}
}
