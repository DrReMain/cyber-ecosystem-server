package menu

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

type CreateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuLogic {
	return &CreateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMenuLogic) CreateMenu(req *types.MenuCreateReq) (resp *types.MenuCreateRes, err error) {

	data, err := l.svcCtx.RPCAdminSystem.MENU.CreateMenu(l.ctx, &admin_system.MenuBody{
		Sort:       req.Sort,
		Title:      req.Title,
		Icon:       req.Icon,
		Code:       req.Code,
		ParentId:   req.ParentID,
		MenuType:   req.MenuType,
		MenuPath:   req.MenuPath,
		Properties: req.Properties,
		Resources:  buildBResources(req.Resources),
	})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	roles, err := l.svcCtx.RPCAdminSystem.CASBIN.QueryRoleByMenu(l.ctx, &admin_system.RoleFromMenuReq{MenuIds: []string{data.Id}})
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

	return &types.MenuCreateRes{
		CommonRes: common_res.NewYES(data.Msg),
		Data:      &data.Id,
	}, nil
}
