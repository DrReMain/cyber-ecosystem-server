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

type UpdateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRoleLogic) UpdateRole(req *types.RoleUpdateReq) (resp *types.RoleUpdateRes, err error) {
	roles, err := l.svcCtx.RPCAdminSystem.CASBIN.QueryRoleByRole(l.ctx, &admin_system.RoleFromRoleReq{RoleIds: []string{*req.ID}})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	data, err := l.svcCtx.RPCAdminSystem.ROLE.UpdateRole(l.ctx, &admin_system.RoleBody{
		Id:       req.ID,
		Sort:     req.Sort,
		RoleName: req.RoleName,
		Code:     req.Code,
		Remark:   req.Remark,
		Menus:    buildBMenus(req.MenuIds),
	})
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

	return &types.RoleUpdateRes{
		CommonRes: common_res.NewYES(data.Msg),
	}, nil
}
