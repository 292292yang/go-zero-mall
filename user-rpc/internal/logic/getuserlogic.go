package logic

import (
	"context"

	"github.com/292292yang/go-zero-mall/common/errorx"
	"github.com/292292yang/go-zero-mall/user-rpc/internal/svc"
	"github.com/292292yang/go-zero-mall/user-rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *user.GetUserReq) (*user.GetUserResp, error) {
	if in.UserId <= 0 {
		return nil, errorx.NewCodeError(errorx.InvalidParam, "user id is empty")
	}
	userModel, err := l.svcCtx.UserRepository.FindById(l.ctx, uint64(in.UserId))
	if err != nil {
		return nil, err
	}
	return &user.GetUserResp{
		UserId:   userModel.Id,
		Username: userModel.Username,
		Nickname: userModel.Nickname,
		Avatar:   userModel.Avatar,
		Status:   int64(userModel.Status),
	}, nil
}
