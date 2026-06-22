package rpcinterceptor

import (
	"context"

	"github.com/292292yang/go-zero-mall/common/rpcx"
	"google.golang.org/grpc"
)

// RpcErrorInterceptor 全局统一错误转换拦截器
func RpcErrorInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		// 调用公共包统一转换业务error为grpc status error
		return nil, rpcx.ToRpcError(err)
	}
	return resp, nil
}
