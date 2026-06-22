// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package product

import (
	"net/http"

	"user-api/internal/logic/product"
	"user-api/internal/svc"
	"user-api/internal/types"

	"github.com/292292yang/go-zero-mall/common/errorx"
	"github.com/292292yang/go-zero-mall/common/response"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ProductCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ProductCreateReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Fail(w, errorx.NewCodeError(errorx.InvalidParam, "参数解析出错"))
			return
		}
		l := product.NewProductCreateLogic(r.Context(), svcCtx)
		resp, err := l.ProductCreate(&req)
		if err != nil {
			response.Fail(w, err)
		} else {
			response.Ok(w, resp)
		}
	}
}
