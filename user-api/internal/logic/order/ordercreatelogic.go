// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package order

import (
	"context"

	"user-api/internal/svc"
	"user-api/internal/types"

	"github.com/292292yang/go-zero-mall/order-rpc/orderclient"
	"github.com/zeromicro/go-zero/core/logx"
)

type OrderCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderCreateLogic {
	return &OrderCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderCreateLogic) OrderCreate(req *types.OrderCreateReq) (resp *types.OrderCreateResp, err error) {
	l.svcCtx.OrderRpc.CreateOrder(l.ctx, &orderclient.CreateOrderReq{})
	return
}
