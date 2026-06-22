// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"

	"user-api/internal/svc"
	"user-api/internal/types"

	"github.com/292292yang/go-zero-mall/common/errorx"
	"github.com/292292yang/go-zero-mall/common/jwtx"
	"github.com/292292yang/go-zero-mall/user-rpc/user"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	if req.Username == "" || req.Password == "" {
		return nil, errorx.NewCodeError(errorx.InvalidParam, "username or password is empty")
	}
	loginResp, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginReq{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	token, expireAt, err := jwtx.GenerateToken(
		l.svcCtx.Config.Auth.AccessSecret,
		jwtx.Now(),
		l.svcCtx.Config.Auth.AccessExpire,
		loginResp.UserId,
	)
	if err != nil {
		return nil, err
	}
	return &types.LoginResp{
		Token:    token,
		ExpireAt: expireAt,
	}, nil
}
