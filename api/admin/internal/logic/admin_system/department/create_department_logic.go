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

type CreateDepartmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDepartmentLogic {
	return &CreateDepartmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDepartmentLogic) CreateDepartment(req *types.DepartmentCreateReq) (resp *types.DepartmentCreateRes, err error) {
	data, err := l.svcCtx.RPCAdminSystem.DEPARTMENT.CreateDepartment(l.ctx, &admin_system.DepartmentBody{
		Sort:           req.Sort,
		DepartmentName: req.DepartmentName,
		Remark:         req.Remark,
		ParentId:       req.ParentID,
	})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	return &types.DepartmentCreateRes{
		CommonRes: common_res.NewYES(data.Msg),
		Result:    &data.Id,
	}, nil
}
