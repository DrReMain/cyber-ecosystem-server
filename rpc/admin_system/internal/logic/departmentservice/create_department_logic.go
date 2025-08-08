package departmentservicelogic

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/pointc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"github.com/rs/xid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDepartmentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDepartmentLogic {
	return &CreateDepartmentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateDepartmentLogic) CreateDepartment(in *admin_system.DepartmentBody) (*admin_system.BaseIDRes, error) {
	// 获取父级Department，没有则为nil
	var parent *ent.Department
	if in.ParentId != nil && *in.ParentId != "" {
		if item, err := l.svcCtx.DB.Department.Get(l.ctx, *in.ParentId); err != nil {
			return nil, ent.DefaultHandleError(l.Logger, err, in)
		} else {
			parent = item
		}
	}

	// 没有父级，path为自身id
	id := xid.New().String()
	var path *string
	if parent != nil {
		path = pointc.P(parent.IDPath + "." + id)
	} else {
		path = pointc.P(id)
	}

	if err := checkLevel(path); err != nil {
		l.Logger.Errorw("department level should not more then 10", logx.Field("detail", in))
		return nil, errorc.GRPCInvalidArgumentError(msgc.CREATE_FAILED)
	}

	item, err := l.svcCtx.DB.Department.Create().
		SetID(id).
		SetNotNilSort(in.Sort).
		SetNotNilDepartmentName(in.DepartmentName).
		SetNotNilRemark(in.Remark).
		SetNotNilParentID(in.ParentId).
		SetNotNilIDPath(path).
		Save(l.ctx)
	if err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	return &admin_system.BaseIDRes{Id: item.ID, Msg: msgc.CREATE_SUCCESS}, nil
}
