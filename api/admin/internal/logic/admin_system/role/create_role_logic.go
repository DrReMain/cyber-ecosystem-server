package role

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

type CreateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRoleLogic) CreateRole(req *types.RoleCreateReq) (resp *types.RoleCreateRes, err error) {
	data, err := l.svcCtx.RPCAdminSystem.ROLE.CreateRole(l.ctx, &admin_system.RoleBody{
		Sort:     req.Sort,
		RoleName: req.RoleName,
		Code:     req.Code,
		Remark:   req.Remark,
		Menus:    buildBMenus(req.MenuIds),
	})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	roles, err := l.svcCtx.RPCAdminSystem.CASBIN.QueryRoleByRole(l.ctx, &admin_system.RoleFromRoleReq{RoleIds: []string{data.Id}})
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

	return &types.RoleCreateRes{
		CommonRes: common_res.NewYES(data.Msg),
		Result:    &data.Id,
	}, nil
}
