package logic

import (
	"context"

	"github.com/292292yang/go-zero-mall/product-rpc/internal/svc"
	"github.com/292292yang/go-zero-mall/product-rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProductLogic {
	return &ListProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListProductLogic) ListProduct(in *product.ListProductReq) (*product.ListProductResp, error) {
	// todo: add your logic here and delete this line

	return &product.ListProductResp{}, nil
}
