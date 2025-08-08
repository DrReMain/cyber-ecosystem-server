package userservicelogic

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/pointc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/department"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/position"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/predicate"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/role"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/user"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryUserLogic {
	return &QueryUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryUserLogic) QueryUser(in *admin_system.UserListReq) (*admin_system.UserListRes, error) {
	items, err := l.svcCtx.DB.User.Query().
		Where(
			ent.NewPredicateUser().ApplyCreatedAt(in.CreatedAt).ApplyUpdatedAt(in.UpdatedAt).
				Apply(in.Status != nil, func() predicate.User {
					return user.StatusEQ(*pointc.PStatus32t8(in.Status))
				}).
				Apply(in.Email != nil, func() predicate.User {
					return user.EmailEQ(*in.Email)
				}).
				Apply(in.Name != nil, func() predicate.User {
					return user.NameContains(*in.Name)
				}).
				Apply(in.Nickname != nil, func() predicate.User {
					return user.NicknameContains(*in.Nickname)
				}).
				Apply(in.Phone != nil, func() predicate.User {
					return user.PhoneEQ(*in.Phone)
				}).
				Apply(in.DepartmentIds != nil, func() predicate.User {
					return user.HasDepartmentWith(department.IDIn(in.DepartmentIds...))
				}).
				Apply(in.PositionIds != nil, func() predicate.User {
					return user.HasPositionsWith(position.IDIn(in.PositionIds...))
				}).
				Apply(in.RoleIds != nil, func() predicate.User {
					return user.HasRolesWith(role.IDIn(in.RoleIds...))
				}).
				Submit()...,
		).
		WithDepartment().
		WithPositions().
		WithRoles().
		Page(l.ctx, in.PageNo, in.PageSize, 10, 100, func(pager *ent.UserPager) {
			pager.Order = []user.OrderOption{
				ent.Asc(user.FieldCreatedAt),
			}
		})
	if err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	res := &admin_system.UserListRes{
		PageNo:   items.PageDetail.PageNo,
		PageSize: items.PageDetail.PageSize,
		Total:    items.PageDetail.Total,
		More:     items.PageDetail.More,
	}

	for _, v := range items.List {
		res.List = append(res.List, &admin_system.UserBody{
			Id:         pointc.P(v.ID),
			CreatedAt:  pointc.P(v.CreatedAt.UnixMilli()),
			UpdatedAt:  pointc.P(v.UpdatedAt.UnixMilli()),
			Status:     pointc.PStatus8t32(&v.Status),
			Password:   nil,
			Email:      pointc.P(v.Email),
			Name:       pointc.P(v.Name),
			Nickname:   pointc.P(v.Nickname),
			Phone:      pointc.P(v.Phone),
			Avatar:     pointc.P(v.Avatar),
			Remark:     pointc.P(v.Remark),
			Department: buildBDepartment(v.Edges.Department),
			Positions:  buildBPositions(v.Edges.Positions),
			Roles:      buildBRoles(v.Edges.Roles),
		})
	}

	return res, nil
}
