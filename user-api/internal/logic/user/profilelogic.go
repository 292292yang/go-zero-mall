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

type ProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProfileLogic {
	return &ProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProfileLogic) Profile() (resp *types.ProfileResp, err error) {
	userId, err := jwtx.GetUserId(l.ctx)
	if err != nil {
		return nil, errorx.NewCodeError(errorx.Unauthorized, "Unauthorized")
	}
	user, err := l.svcCtx.UserRpc.GetUser(l.ctx, &user.GetUserReq{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}
	return &types.ProfileResp{
		Id:       user.UserId,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}, nil
}
