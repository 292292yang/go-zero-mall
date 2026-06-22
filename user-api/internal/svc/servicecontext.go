// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"user-api/internal/config"

	"github.com/292292yang/go-zero-mall/order-rpc/orderclient"
	"github.com/292292yang/go-zero-mall/product-rpc/productclient"
	"github.com/292292yang/go-zero-mall/user-rpc/userclient"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	UserRpc    userclient.User
	ProductRpc productclient.Product
	OrderRpc   orderclient.Order
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		UserRpc:    userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		ProductRpc: productclient.NewProduct(zrpc.MustNewClient(c.ProductRpc)),
		OrderRpc:   orderclient.NewOrder(zrpc.MustNewClient(c.OrderRpc)),
	}
}
