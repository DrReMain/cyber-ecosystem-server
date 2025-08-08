package roleservicelogic

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateRoleLogic) UpdateRole(in *admin_system.RoleBody) (*admin_system.BaseRes, error) {
	if in.Id == nil || *in.Id == "" {
		l.Logger.Errorw("id is empty", logx.Field("detail", in))
		return nil, errorc.GRPCInvalidArgumentError(msgc.UPDATE_FAILED)
	}

	if err := ent.WithTX(l.ctx, l.svcCtx.DB, func(tx *ent.Tx) error {
		update := tx.Role.UpdateOneID(*in.Id).
			SetNotNilSort(in.Sort).
			SetNotNilRoleName(in.RoleName).
			SetNotNilCode(in.Code).
			SetNotNilRemark(in.Remark)

		if in.Menus != nil && len(in.Menus) > 0 {
			if err := tx.Role.UpdateOneID(*in.Id).ClearMenus().Exec(l.ctx); err != nil {
				return err
			}
			update = update.AddMenuIDs(buildEMenuIDs(in.Menus)...)
		}

		return update.Exec(l.ctx)
	}); err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	return &admin_system.BaseRes{Msg: msgc.UPDATE_SUCCESS}, nil
}
