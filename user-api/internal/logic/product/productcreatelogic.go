// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package product

import (
	"context"

	"user-api/internal/svc"
	"user-api/internal/types"

	"github.com/292292yang/go-zero-mall/common/errorx"
	"github.com/292292yang/go-zero-mall/common/rpcx"
	"github.com/292292yang/go-zero-mall/product-rpc/product"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProductCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProductCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductCreateLogic {
	return &ProductCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProductCreateLogic) ProductCreate(req *types.ProductCreateReq) (resp *types.ProductCreateResp, err error) {
	l.Infof("ProductCreateLogic, product:%+v", req)
	product, err := l.svcCtx.ProductRpc.CreateProduct(l.ctx, &product.CreateProductReq{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	})
	l.Infof("ProductCreateLogic, product:%+v, err:%+v", product, err)
	if err != nil {
		l.Errorf("call product rpc CreateProduct failed, err=%v", err)
		return nil, rpcx.ConvertRpcError(err, errorx.ProductOffShelf, "创建商品失败")
	}
	return &types.ProductCreateResp{
		ProductId: product.ProductId,
	}, nil
}
