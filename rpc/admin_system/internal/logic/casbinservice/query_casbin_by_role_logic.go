package casbinservicelogic

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/orm/ent/mixins"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/menu"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/role"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryCasbinByRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryCasbinByRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryCasbinByRoleLogic {
	return &QueryCasbinByRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryCasbinByRoleLogic) QueryCasbinByRole(in *admin_system.CasbinReq) (*admin_system.CasbinRes, error) {
	items, err := l.svcCtx.DB.Role.Query().
		Where(role.CodeIn(in.RoleCode...)).
		WithMenus(func(query *ent.MenuQuery) {
			query = query.Where(menu.StatusEQ(mixins.StatusNormal)).WithResources()
		}).
		All(l.ctx)
	if err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	return &admin_system.CasbinRes{
		List: buildBList(items),
	}, nil
}
