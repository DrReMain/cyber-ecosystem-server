package user

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/helper/common_res"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/svc"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/types"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserOneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserOneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserOneLogic {
	return &DeleteUserOneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserOneLogic) DeleteUserOne(req *types.UserDeleteReq) (resp *types.UserDeleteRes, err error) {
	data, err := l.svcCtx.RPCAdminSystem.USER.DeleteUser(l.ctx, &admin_system.IDsReq{Ids: []string{*req.ID}})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	return &types.UserDeleteRes{
		CommonRes: common_res.NewYES(data.Msg),
	}, nil
}
