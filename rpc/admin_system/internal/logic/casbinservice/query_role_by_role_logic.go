package casbinservicelogic

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/role"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryRoleByRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryRoleByRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryRoleByRoleLogic {
	return &QueryRoleByRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryRoleByRoleLogic) QueryRoleByRole(in *admin_system.RoleFromRoleReq) (*admin_system.RoleFromRoleRes, error) {
	items, err := l.svcCtx.DB.Role.Query().
		Where(role.IDIn(in.RoleIds...)).
		Select(role.FieldCode).Strings(l.ctx)
	if err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	return &admin_system.RoleFromRoleRes{
		RoleCode: items,
	}, nil
}
