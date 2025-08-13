package menu

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/helper/casbin_rules"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/helper/common_res"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/svc"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/types"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMenuOneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMenuOneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuOneLogic {
	return &DeleteMenuOneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMenuOneLogic) DeleteMenuOne(req *types.MenuDeleteReq) (resp *types.MenuDeleteRes, err error) {
	roles, err := l.svcCtx.RPCAdminSystem.CASBIN.QueryRoleByMenu(l.ctx, &admin_system.RoleFromMenuReq{MenuIds: []string{*req.ID}})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	data, err := l.svcCtx.RPCAdminSystem.MENU.DeleteMenu(l.ctx, &admin_system.IDsReq{Ids: []string{*req.ID}})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	rule, err := l.svcCtx.RPCAdminSystem.CASBIN.QueryCasbinByRole(l.ctx, &admin_system.CasbinReq{RoleCode: roles.RoleCode})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}
	if err := casbin_rules.RefreshCasbinRules(l.svcCtx.Casbin, roles.RoleCode, rule.List); err != nil {
		return nil, errorc.NewHTTPInternal(msgc.SYSTEM_ERROR, err.Error())
	}

	return &types.MenuDeleteRes{
		CommonRes: common_res.NewYES(data.Msg),
	}, nil
}
