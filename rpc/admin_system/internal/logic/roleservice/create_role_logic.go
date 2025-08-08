package roleservicelogic

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateRoleLogic) CreateRole(in *admin_system.RoleBody) (*admin_system.BaseIDRes, error) {
	create := l.svcCtx.DB.Role.Create().
		SetNotNilSort(in.Sort).
		SetNotNilRoleName(in.RoleName).
		SetNotNilCode(in.Code).
		SetNotNilRemark(in.Remark)

	if in.Menus != nil && len(in.Menus) > 0 {
		create = create.AddMenuIDs(buildEMenuIDs(in.Menus)...)
	}

	item, err := create.Save(l.ctx)
	if err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	return &admin_system.BaseIDRes{Id: item.ID, Msg: msgc.CREATE_SUCCESS}, nil
}
