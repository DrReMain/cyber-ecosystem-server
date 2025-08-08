package casbinservicelogic

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/menu"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/role"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryRoleByMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryRoleByMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryRoleByMenuLogic {
	return &QueryRoleByMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryRoleByMenuLogic) QueryRoleByMenu(in *admin_system.RoleFromMenuReq) (*admin_system.RoleFromMenuRes, error) {
	items, err := l.svcCtx.DB.Role.Query().
		Where(role.HasMenusWith(menu.IDIn(in.MenuIds...))).
		Select(role.FieldCode).Strings(l.ctx)
	if err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	return &admin_system.RoleFromMenuRes{
		RoleCode: items,
	}, nil
}
