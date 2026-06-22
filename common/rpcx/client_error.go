package rpcx

import (
	"github.com/292292yang/go-zero-mall/common/errorx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ConvertRpcError(err error, fallbackCode int64, fallbackMsg string) error {
	if err == nil {
		return nil
	}

	if ce, ok := FromRpcError(err); ok {
		return ce
	}

	st, ok := status.FromError(err)
	if !ok {
		return errorx.NewCodeError(errorx.ServerError, "系统繁忙，请稍后重试")
	}

	switch st.Code() {
	case codes.Unavailable, codes.DeadlineExceeded:
		return errorx.NewCodeError(errorx.DownstreamFailure, fallbackMsg)
	default:
		return errorx.NewCodeError(fallbackCode, fallbackMsg)
	}
}
