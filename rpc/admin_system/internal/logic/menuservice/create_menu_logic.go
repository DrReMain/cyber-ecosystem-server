package menuservicelogic

import (
	"context"
	"strings"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/pointc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuLogic {
	return &CreateMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateMenuLogic) CreateMenu(in *admin_system.MenuBody) (*admin_system.BaseIDRes, error) {
	if in.Code == nil || *in.Code == "" {
		l.Logger.Errorw("code is empty", logx.Field("detail", in))
		return nil, errorc.GRPCInvalidArgumentError(msgc.CREATE_FAILED)
	} else {
		c := strings.ReplaceAll(*in.Code, ".", "")
		in.Code = &c
	}

	var parent *ent.Menu
	if in.ParentId != nil && *in.ParentId != "" {
		if item, err := l.svcCtx.DB.Menu.Get(l.ctx, *in.ParentId); err != nil {
			return nil, ent.DefaultHandleError(l.Logger, err, in)
		} else {
			parent = item
		}
	}

	// 如果没传父级id，path为自身Code
	var path *string
	if parent != nil {
		path = pointc.P(parent.CodePath + "." + *in.Code)
	} else {
		path = in.Code
	}

	if err := checkLevel(path); err != nil {
		l.Logger.Errorw("menu level should not more than 5", logx.Field("detail", in))
		return nil, errorc.GRPCInvalidArgumentError(msgc.CREATE_FAILED)
	}

	var item *ent.Menu
	if err := ent.WithTX(l.ctx, l.svcCtx.DB, func(tx *ent.Tx) error {
		menu, err := tx.Menu.Create().
			SetNotNilSort(in.Sort).
			SetNotNilStatus(pointc.PStatus32t8(in.Status)).
			SetNotNilTitle(in.Title).
			SetNotNilIcon(in.Icon).
			SetNotNilCode(in.Code).
			SetNotNilCodePath(path).
			SetNotNilParentID(in.ParentId).
			SetNotNilMenuType(in.MenuType).
			SetNotNilMenuPath(in.MenuPath).
			SetNotNilProperties(in.Properties).
			Save(l.ctx)
		if err != nil {
			return err
		}

		var resources = make([]*ent.ResourceCreate, len(in.Resources))
		for i, v := range in.Resources {
			resources[i] = tx.Resource.Create().SetMenuID(menu.ID).SetNotNilMethod(v.Method).SetNotNilPath(v.Path)
		}
		if err := tx.Resource.CreateBulk(resources...).Exec(l.ctx); err != nil {
			return err
		}

		item = menu
		return nil
	}); err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	return &admin_system.BaseIDRes{Id: item.ID, Msg: msgc.CREATE_SUCCESS}, nil
}
