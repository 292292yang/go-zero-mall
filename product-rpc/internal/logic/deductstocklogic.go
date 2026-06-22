package logic

import (
	"context"

	"github.com/292292yang/go-zero-mall/product-rpc/internal/svc"
	"github.com/292292yang/go-zero-mall/product-rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeductStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeductStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeductStockLogic {
	return &DeductStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeductStockLogic) DeductStock(in *product.DeductStockReq) (*product.DeductStockResp, error) {
	// todo: add your logic here and delete this line

	return &product.DeductStockResp{}, nil
}
