package logic

import (
	"context"

	"github.com/292292yang/go-zero-mall/product-rpc/internal/svc"
	"github.com/292292yang/go-zero-mall/product-rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckStockLogic {
	return &CheckStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckStockLogic) CheckStock(in *product.CheckStockReq) (*product.CheckStockResp, error) {
	// todo: add your logic here and delete this line

	return &product.CheckStockResp{}, nil
}
