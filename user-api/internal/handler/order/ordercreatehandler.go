// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package order

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"user-api/internal/logic/order"
	"user-api/internal/svc"
	"user-api/internal/types"
)

func OrderCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrderCreateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := order.NewOrderCreateLogic(r.Context(), svcCtx)
		resp, err := l.OrderCreate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
