package roleservicelogic

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/pointc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/role"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleLogic {
	return &GetRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRoleLogic) GetRole(in *admin_system.IDReq) (*admin_system.RoleBody, error) {
	item, err := l.svcCtx.DB.Role.Query().
		Where(role.IDEQ(in.Id)).
		WithMenus().
		First(l.ctx)
	if err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	return &admin_system.RoleBody{
		Id:        pointc.P(item.ID),
		CreatedAt: pointc.P(item.CreatedAt.UnixMilli()),
		UpdatedAt: pointc.P(item.UpdatedAt.UnixMilli()),
		Sort:      pointc.P(item.Sort),
		RoleName:  pointc.P(item.RoleName),
		Code:      pointc.P(item.Code),
		Remark:    pointc.P(item.Remark),
		Menus:     buildBMenu(item.Edges.Menus),
	}, nil
}
