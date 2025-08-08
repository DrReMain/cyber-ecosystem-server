package department

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/helper/common_res"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/svc"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/types"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDepartmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDepartmentLogic {
	return &UpdateDepartmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDepartmentLogic) UpdateDepartment(req *types.DepartmentUpdateReq) (resp *types.DepartmentUpdateRes, err error) {
	data, err := l.svcCtx.RPCAdminSystem.DEPARTMENT.UpdateDepartment(l.ctx, &admin_system.DepartmentBody{
		Id:             req.ID,
		Sort:           req.Sort,
		DepartmentName: req.DepartmentName,
		Remark:         req.Remark,
		ParentId:       req.ParentID,
	})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	return &types.DepartmentUpdateRes{
		CommonRes: common_res.NewYES(data.Msg),
	}, nil
}
