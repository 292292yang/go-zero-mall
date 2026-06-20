package response

import (
	"errors"
	"net/http"

	"github.com/go-zero-mall/common/errorx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Ok(w http.ResponseWriter, data interface{}) {
	httpx.OkJson(w, Body{
		Code:    errorx.Success,
		Message: "success",
		Data:    data,
	})
}

func Fail(w http.ResponseWriter, err error) {
	var e *errorx.CodeError
	if errors.As(err, &e) {
		httpx.OkJson(w, Body{
			Code:    e.Code,
			Message: e.Msg,
		})
		return
	}
	httpx.OkJson(w, Body{
		Code:    errorx.ServerError,
		Message: "系统繁忙，请稍后重试",
	})
}
