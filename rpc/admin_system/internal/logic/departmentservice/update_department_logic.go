package departmentservicelogic

import (
	"context"
	"errors"
	"strings"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/pointc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/department"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDepartmentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDepartmentLogic {
	return &UpdateDepartmentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateDepartmentLogic) UpdateDepartment(in *admin_system.DepartmentBody) (*admin_system.BaseRes, error) {
	if in.Id == nil || *in.Id == "" {
		l.Logger.Errorw("id is empty", logx.Field("detail", in))
		return nil, errorc.GRPCInvalidArgumentError(msgc.UPDATE_FAILED)
	}

	var parent *ent.Department
	if in.ParentId != nil && *in.ParentId != "" {
		if item, err := l.svcCtx.DB.Department.Get(l.ctx, *in.ParentId); err != nil {
			return nil, ent.DefaultHandleError(l.Logger, err, in)
		} else {
			parent = item
		}
	}

	id := *in.Id
	var path *string
	// parent_id传了，但是 parent_id == "", 说明当前部门改为1级
	if in.ParentId != nil && *in.ParentId == "" {
		path = pointc.P(id)
	}
	// parent_id传了，并且 parent_id != ""，拼接新的parent下的id_path
	if parent != nil {
		if strings.Contains(parent.IDPath, *in.Id) {
			l.Logger.Errorw("department can't move to subordinate", logx.Field("detail", in))
			return nil, errorc.GRPCInvalidArgumentError(msgc.UPDATE_FAILED)
		}
		path = pointc.P(parent.IDPath + "." + id)
	}

	if err := checkLevel(path); err != nil {
		l.Logger.Errorw("department level should not more then 10", logx.Field("detail", in))
		return nil, errorc.GRPCInvalidArgumentError(msgc.UPDATE_FAILED)
	}

	if err := ent.WithTX(l.ctx, l.svcCtx.DB, func(tx *ent.Tx) error {
		prev, err := tx.Department.Get(l.ctx, id)
		if err != nil {
			return err
		}

		if err := tx.Department.UpdateOneID(*in.Id).
			SetNotNilSort(in.Sort).
			SetNotNilDepartmentName(in.DepartmentName).
			SetNotNilRemark(in.Remark).
			SetNotNilParentID(in.ParentId).
			SetNotNilIDPath(path).
			Exec(l.ctx); err != nil {
			return err
		}

		// 如果有path，就是有parent_id有传入，并且path与原数据不同，需要修改所有子Department
		if path != nil && prev.IDPath != *path {
			children, err := tx.Department.Query().Where(department.IDPathHasPrefix(prev.IDPath + ".")).All(l.ctx)
			if err != nil {
				return err
			}

			for _, child := range children {
				childIdPath := strings.Replace(child.IDPath, prev.IDPath, *path, 1)
				if err := checkLevel(&childIdPath); err != nil {
					return errors.New(msgc.UPDATE_FAILED)
				}
				if err := tx.Department.UpdateOneID(child.ID).SetIDPath(childIdPath).Exec(l.ctx); err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	return &admin_system.BaseRes{Msg: msgc.UPDATE_SUCCESS}, nil
}
