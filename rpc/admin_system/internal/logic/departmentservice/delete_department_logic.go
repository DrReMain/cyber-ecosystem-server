package departmentservicelogic

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/department"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/user"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDepartmentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDepartmentLogic {
	return &DeleteDepartmentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteDepartmentLogic) DeleteDepartment(in *admin_system.IDsReq) (*admin_system.BaseRes, error) {
	if exist, err := l.svcCtx.DB.Department.Query().Where(department.ParentIDIn(in.Ids...)).Exist(l.ctx); err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	} else if exist {
		l.Logger.Errorw("department has subordinate", logx.Field("detail", in))
		return nil, errorc.GRPCInvalidArgumentError(msgc.DELETE_FAILED)
	}

	if exist, err := l.svcCtx.DB.User.Query().Where(user.DepartmentIDIn(in.Ids...)).Exist(l.ctx); err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	} else if exist {
		l.Logger.Errorw("department used by a user", logx.Field("detail", in))
		return nil, errorc.GRPCInvalidArgumentError(msgc.DELETE_FAILED)
	}

	if _, err := l.svcCtx.DB.Department.Delete().Where(department.IDIn(in.Ids...)).Exec(l.ctx); err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	return &admin_system.BaseRes{Msg: msgc.DELETE_SUCCESS}, nil
}
