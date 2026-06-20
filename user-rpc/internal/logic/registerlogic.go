package logic

import (
	"context"

	"github.com/292292yang/go-zero-mall/common/cryptx"
	"github.com/292292yang/go-zero-mall/common/errorx"
	"github.com/292292yang/go-zero-mall/user-rpc/internal/repository"
	"github.com/292292yang/go-zero-mall/user-rpc/internal/svc"
	"github.com/292292yang/go-zero-mall/user-rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	if in.Username == "" || in.Password == "" {
		return nil, errorx.NewCodeError(errorx.InvalidParam, "username or password is empty")
	}
	if in.Nickname == "" {
		in.Nickname = in.Username
	}
	encryptPassword, err := cryptx.EncryptPassword(in.Password)
	if err != nil {
		return nil, errorx.NewCodeError(errorx.ServerError, "系统错误")
	}
	userId, err := l.svcCtx.UserRepository.Create(l.ctx, repository.UserCreate{
		Username: in.Username,
		Password: encryptPassword,
		Nickname: in.Nickname,
		Avatar:   in.Avatar,
		Status:   0,
	})
	if err != nil {
		return nil, err
	}
	return &user.RegisterResp{
		UserId: int64(userId),
	}, nil
}
