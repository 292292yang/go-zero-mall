package logic

import (
	"context"

	"github.com/292292yang/go-zero-mall/user-rpc/internal/svc"
	"github.com/292292yang/go-zero-mall/user-rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckUserAvailableLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckUserAvailableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckUserAvailableLogic {
	return &CheckUserAvailableLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckUserAvailableLogic) CheckUserAvailable(in *user.CheckUserAvailableReq) (*user.CheckUserAvailableResp, error) {
	// todo: add your logic here and delete this line

	return &user.CheckUserAvailableResp{}, nil
}
