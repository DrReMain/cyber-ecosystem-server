package menu

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/helper/common_res"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/svc"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/types"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryMenuLogic {
	return &QueryMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryMenuLogic) QueryMenu(req *types.MenuQueryReq) (resp *types.MenuQueryRes, err error) {
	data, err := l.svcCtx.RPCAdminSystem.MENU.QueryMenu(l.ctx, &admin_system.MenuListReq{})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	return &types.MenuQueryRes{
		CommonRes: common_res.NewYES(""),
		Data: &types.MenuQuery{
			List: buildTChildren(data.List),
		},
	}, nil
}
