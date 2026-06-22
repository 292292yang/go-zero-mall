package logic

import (
	"context"
	"order-rpc/internal/repository"

	"order-rpc/internal/svc"
	"order-rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateOrderLogic) CreateOrder(in *order.CreateOrderReq) (*order.CreateOrderResp, error) {
	orderId, err := l.svcCtx.OrderRepository.Create(l.ctx, repository.OrderCreate{
		UserId:    in.UserId,
		ProductId: in.ProductId,
		Quantity:  in.Quantity,
		Price:     0,
		Status:    1,
	})
	if err != nil {
		return nil, err
	}
	return &order.CreateOrderResp{
		OrderId: int64(orderId),
	}, nil
}
