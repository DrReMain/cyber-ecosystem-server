package menuservicelogic

import (
	"context"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/pointc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/resource"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateMenuLogic) UpdateMenu(in *admin_system.MenuBody) (*admin_system.BaseRes, error) {
	if in.Id == nil || *in.Id == "" {
		l.Logger.Errorw("id is empty", logx.Field("detail", in))
		return nil, errorc.GRPCInvalidArgumentError(msgc.UPDATE_FAILED)
	}

	if err := ent.WithTX(l.ctx, l.svcCtx.DB, func(tx *ent.Tx) error {
		prev, err := tx.Menu.Get(l.ctx, *in.Id)
		if err != nil {
			return err
		}
		// 先删除Menu下所有的Resource
		if _, err := tx.Resource.Delete().Where(resource.MenuIDEQ(prev.ID)).Exec(l.ctx); err != nil {
			return err
		}

		var resources = make([]*ent.ResourceCreate, len(in.Resources))
		for i, v := range in.Resources {
			resources[i] = tx.Resource.Create().SetMenuID(prev.ID).SetNotNilMethod(v.Method).SetNotNilPath(v.Path)
		}
		if err := tx.Resource.CreateBulk(resources...).Exec(l.ctx); err != nil {
			return err
		}

		if err := tx.Menu.UpdateOneID(prev.ID).
			SetNotNilSort(in.Sort).
			SetNotNilStatus(pointc.PStatus32t8(in.Status)).
			SetNotNilTitle(in.Title).
			SetNotNilIcon(in.Icon).
			SetNotNilProperties(in.Properties).
			Exec(l.ctx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	return &admin_system.BaseRes{Msg: msgc.UPDATE_SUCCESS}, nil
}
