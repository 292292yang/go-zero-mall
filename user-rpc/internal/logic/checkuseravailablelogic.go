package logic

import (
	"context"

	"github.com/292292yang/go-zero-mall/common/errorx"
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
	if in.UserId <= 0 {
		return nil, errorx.NewCodeError(errorx.InvalidParam, "user id is empty")
	}
	userModel, err := l.svcCtx.UserRepository.FindById(l.ctx, uint64(in.UserId))
	if err != nil {
		return nil, err
	}
	if userModel.Status != 1 {
		return &user.CheckUserAvailableResp{
			Available: false,
			Message:   "用户不可用",
		}, nil
	}
	return &user.CheckUserAvailableResp{
		Available: true,
		Message:   "用户可用",
	}, nil
}
