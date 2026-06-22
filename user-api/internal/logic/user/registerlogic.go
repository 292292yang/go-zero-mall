// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"
	"user-api/internal/svc"
	"user-api/internal/types"

	"github.com/292292yang/go-zero-mall/common/errorx"
	"github.com/292292yang/go-zero-mall/common/rpcx"
	"github.com/292292yang/go-zero-mall/user-rpc/user"
	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	registerResp, err := l.svcCtx.UserRpc.Register(l.ctx, &user.RegisterReq{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Avatar:   "",
	})
	if err != nil {
		l.Errorf("call order rpc CreateOrder failed, err=%v", err)
		return nil, rpcx.ConvertRpcError(err, errorx.UserDisabled, "用户注册失败")
	}
	return &types.RegisterResp{
		UserId: registerResp.UserId,
	}, nil
}
