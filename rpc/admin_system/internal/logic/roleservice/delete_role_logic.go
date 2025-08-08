package roleservicelogic

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/role"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/user"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteRoleLogic) DeleteRole(in *admin_system.IDsReq) (*admin_system.BaseRes, error) {
	if exist, err := l.svcCtx.DB.User.Query().Where(user.HasRolesWith(role.IDIn(in.Ids...))).Exist(l.ctx); err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	} else if exist {
		l.Logger.Errorw("role used by a user", logx.Field("detail", in))
		return nil, errorc.GRPCInvalidArgumentError(msgc.DELETE_FAILED)
	}

	if err := ent.WithTX(l.ctx, l.svcCtx.DB, func(tx *ent.Tx) error {
		if err := tx.Role.Update().Where(role.IDIn(in.Ids...)).ClearMenus().Exec(l.ctx); err != nil {
			return err
		}
		if _, err := tx.Role.Delete().Where(role.IDIn(in.Ids...)).Exec(l.ctx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	return &admin_system.BaseRes{Msg: msgc.DELETE_SUCCESS}, nil
}
