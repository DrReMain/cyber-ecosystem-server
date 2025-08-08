package role

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/helper/casbin_rules"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/helper/common_res"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/svc"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/types"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteRoleLogic) DeleteRole(req *types.IDsReq) (resp *types.CommonRes, err error) {
	roles, err := l.svcCtx.RPCAdminSystem.CASBIN.QueryRoleByRole(l.ctx, &admin_system.RoleFromRoleReq{RoleIds: req.IDs})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	data, err := l.svcCtx.RPCAdminSystem.ROLE.DeleteRole(l.ctx, &admin_system.IDsReq{Ids: req.IDs})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	rule, err := l.svcCtx.RPCAdminSystem.CASBIN.QueryCasbinByRole(l.ctx, &admin_system.CasbinReq{RoleCode: roles.RoleCode})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}
	if err := casbin_rules.RefreshCasbinRules(l.svcCtx.Casbin, roles.RoleCode, rule.List); err != nil {
		return nil, errorc.NewUnknownError(err)
	}

	return common_res.NewYES(data.Msg), nil
}
