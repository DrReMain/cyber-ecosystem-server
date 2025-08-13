package account

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/svc"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AccountInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取当前账号信息
func NewAccountInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccountInfoLogic {
	return &AccountInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AccountInfoLogic) AccountInfo(req *types.AccountInfoReq) (resp *types.AccountInfoRes, err error) {
	// todo: add your logic here and delete this line

	return
}
