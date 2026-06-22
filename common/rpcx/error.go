package rpcx

import (
	"encoding/json"
	"errors"

	"github.com/292292yang/go-zero-mall/common/errorx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type rpcErrorBody struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

func ToRpcError(err error) error {
	if err == nil {
		return nil
	}

	var ce *errorx.CodeError
	if errors.As(err, &ce) {
		bytes, _ := json.Marshal(rpcErrorBody{
			Code:    ce.Code,
			Message: ce.Msg,
		})

		return status.Error(mapCodeToGrpcCode(ce.Code), string(bytes))
	}

	return status.Error(codes.Internal, "internal server error")
}

func FromRpcError(err error) (*errorx.CodeError, bool) {
	if err == nil {
		return nil, false
	}

	st, ok := status.FromError(err)
	if !ok {
		return nil, false
	}

	var body rpcErrorBody
	if json.Unmarshal([]byte(st.Message()), &body) != nil {
		return nil, false
	}

	if body.Code == 0 || body.Message == "" {
		return nil, false
	}

	return &errorx.CodeError{
		Code: body.Code,
		Msg:  body.Message,
	}, true
}

func IsSystemRpcError(err error) bool {
	if err == nil {
		return false
	}

	st, ok := status.FromError(err)
	if !ok {
		return true
	}

	switch st.Code() {
	case codes.Unavailable, codes.DeadlineExceeded, codes.ResourceExhausted:
		return true
	default:
		return false
	}
}

func mapCodeToGrpcCode(code int64) codes.Code {
	switch code {
	case errorx.InvalidParam:
		return codes.InvalidArgument

	case errorx.UserNotFound,
		errorx.ProductNotFound,
		errorx.OrderNotFound,
		errorx.PaymentNotFound:
		return codes.NotFound

	case errorx.UserAlreadyExists,
		errorx.PaymentAlreadyExists:
		return codes.AlreadyExists

	case errorx.UserDisabled,
		errorx.ProductOffShelf,
		errorx.StockNotEnough,
		errorx.OrderStatusInvalid,
		errorx.PaymentStatusInvalid,
		errorx.PaymentAlreadySuccess:
		return codes.FailedPrecondition

	default:
		return codes.Internal
	}
}
