// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"net/http"

	"user-api/internal/logic/user"
	"user-api/internal/svc"
	"user-api/internal/types"

	"github.com/292292yang/go-zero-mall/common/response"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Fail(w, err)
			return
		}
		l := user.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.Register(&req)
		if err != nil {
			response.Fail(w, err)
		} else {
			response.Ok(w, resp)
		}
	}
}
