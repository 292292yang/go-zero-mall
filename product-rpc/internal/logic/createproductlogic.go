package logic

import (
	"context"

	"github.com/292292yang/go-zero-mall/product-rpc/internal/repository"
	"github.com/292292yang/go-zero-mall/product-rpc/internal/svc"
	"github.com/292292yang/go-zero-mall/product-rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductLogic {
	return &CreateProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateProductLogic) CreateProduct(in *product.CreateProductReq) (*product.CreateProductResp, error) {
	productId, err := l.svcCtx.ProductRepository.Create(l.ctx, repository.ProductCreate{
		Name:        in.Name,
		Description: in.Description,
		Price:       in.Price,
		Stock:       in.Stock,
		Status:      1,
	})
	if err != nil {
		return nil, err
	}
	return &product.CreateProductResp{
		ProductId: int64(productId),
	}, nil
}
