package userservicelogic

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/encrypt"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/pointc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserLogic) UpdateUser(in *admin_system.UserBody) (*admin_system.BaseRes, error) {
	if in.Id == nil || *in.Id == "" {
		l.Logger.Errorw("id is empty", logx.Field("detail", in))
		return nil, errorc.GRPCInvalidArgumentError(msgc.UPDATE_FAILED)
	}

	if err := ent.WithTX(l.ctx, l.svcCtx.DB, func(tx *ent.Tx) error {
		update := tx.User.UpdateOneID(*in.Id).
			SetNotNilStatus(pointc.PStatus32t8(in.Status)).
			SetNotNilEmail(in.Email).
			SetNotNilName(in.Name).
			SetNotNilNickname(in.Nickname).
			SetNotNilPhone(in.Phone).
			SetNotNilAvatar(in.Avatar).
			SetNotNilRemark(in.Remark)

		if in.Password != nil {
			update = update.SetNotNilPassword(pointc.P(encrypt.EncryptGenerate(*in.Password)))
		}

		if in.Department != nil {
			update = update.SetNotNilDepartmentID(in.Department.Id)
		}

		if in.Positions != nil && len(in.Positions) > 0 {
			if err := tx.User.UpdateOneID(*in.Id).ClearPositions().Exec(l.ctx); err != nil {
				return err
			}

			ids := make([]string, len(in.Positions))
			for i, p := range in.Positions {
				ids[i] = *p.Id
			}
			update = update.AddPositionIDs(ids...)
		}
		if in.Roles != nil && len(in.Roles) > 0 {
			if err := tx.User.UpdateOneID(*in.Id).ClearRoles().Exec(l.ctx); err != nil {
				return err
			}

			ids := make([]string, len(in.Roles))
			for i, p := range in.Roles {
				ids[i] = *p.Id
			}
			update = update.AddRoleIDs(ids...)
		}

		return update.Exec(l.ctx)
	}); err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	return &admin_system.BaseRes{Msg: msgc.UPDATE_SUCCESS}, nil
}
