package roleservicelogic

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/pointc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/predicate"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/role"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryRoleLogic {
	return &QueryRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryRoleLogic) QueryRole(in *admin_system.RoleListReq) (*admin_system.RoleListRes, error) {
	items, err := l.svcCtx.DB.Role.Query().
		Where(
			ent.NewPredicateRole().ApplyCreatedAt(in.CreatedAt).ApplyUpdatedAt(in.UpdatedAt).
				Apply(in.RoleName != nil, func() predicate.Role {
					return role.RoleNameContains(*in.RoleName)
				}).
				Apply(in.Code != nil, func() predicate.Role {
					return role.CodeEQ(*in.Code)
				}).
				Submit()...,
		).
		WithMenus().
		Page(l.ctx, in.PageNo, in.PageSize, 0, 0, func(pager *ent.RolePager) {
			pager.Order = []role.OrderOption{
				ent.Desc(role.FieldSort),
				ent.Asc(role.FieldCreatedAt),
			}
		})
	if err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	res := &admin_system.RoleListRes{
		PageNo:   items.PageDetail.PageNo,
		PageSize: items.PageDetail.PageSize,
		Total:    items.PageDetail.Total,
		More:     items.PageDetail.More,
	}

	for _, v := range items.List {
		res.List = append(res.List, &admin_system.RoleBody{
			Id:        pointc.P(v.ID),
			CreatedAt: pointc.P(v.CreatedAt.UnixMilli()),
			UpdatedAt: pointc.P(v.UpdatedAt.UnixMilli()),
			Sort:      pointc.P(v.Sort),
			RoleName:  pointc.P(v.RoleName),
			Code:      pointc.P(v.Code),
			Remark:    pointc.P(v.Remark),
			Menus:     buildBMenu(v.Edges.Menus),
		})
	}

	return res, nil
}
