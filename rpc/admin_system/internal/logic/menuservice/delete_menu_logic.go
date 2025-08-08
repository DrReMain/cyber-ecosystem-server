package menuservicelogic

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/menu"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/resource"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/role"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMenuLogic) DeleteMenu(in *admin_system.IDsReq) (*admin_system.BaseRes, error) {
	if exist, err := l.svcCtx.DB.Menu.Query().Where(menu.ParentIDIn(in.Ids...)).Exist(l.ctx); err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	} else if exist {
		l.Logger.Errorw("menu has subordinate", logx.Field("detail", in))
		return nil, errorc.GRPCInvalidArgumentError(msgc.DELETE_FAILED)
	}

	if exist, err := l.svcCtx.DB.Role.Query().Where(role.HasMenusWith(menu.IDIn(in.Ids...))).Exist(l.ctx); err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	} else if exist {
		l.Logger.Errorw("menu used by a role", logx.Field("detail", in))
		return nil, errorc.GRPCInvalidArgumentError(msgc.DELETE_FAILED)
	}

	if err := ent.WithTX(l.ctx, l.svcCtx.DB, func(tx *ent.Tx) error {
		// 先删除Menu下所有的Resource
		if _, err := tx.Resource.Delete().Where(resource.MenuIDIn(in.Ids...)).Exec(l.ctx); err != nil {
			return err
		}
		if _, err := tx.Menu.Delete().Where(menu.IDIn(in.Ids...)).Exec(l.ctx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	return &admin_system.BaseRes{Msg: msgc.DELETE_SUCCESS}, nil
}
