package logic

import (
	"context"

	"github.com/292292yang/go-zero-mall/common/cryptx"
	"github.com/292292yang/go-zero-mall/common/errorx"
	"github.com/292292yang/go-zero-mall/user-rpc/internal/svc"
	"github.com/292292yang/go-zero-mall/user-rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	userModel, err := l.svcCtx.UserRepository.FindByUsername(l.ctx, in.Username)
	if err != nil {
		return nil, err
	}
	if cryptx.CheckPassword(in.Password, userModel.Password) {
		return nil, errorx.NewCodeError(errorx.UserInvalidPassword, "username or password error")
	}
	return &user.LoginResp{
		UserId: userModel.Id,
	}, nil
}
