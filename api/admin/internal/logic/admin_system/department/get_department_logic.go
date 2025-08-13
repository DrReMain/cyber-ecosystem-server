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

type GetDepartmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDepartmentLogic {
	return &GetDepartmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDepartmentLogic) GetDepartment(req *types.DepartmentGetReq) (resp *types.DepartmentGetRes, err error) {
	data, err := l.svcCtx.RPCAdminSystem.DEPARTMENT.GetDepartment(l.ctx, &admin_system.IDReq{Id: *req.ID})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	return &types.DepartmentGetRes{
		CommonRes: common_res.NewYES(""),
		Result: &types.DepartmentGet{
			ID:             data.Id,
			CreatedAt:      data.CreatedAt,
			UpdatedAt:      data.UpdatedAt,
			Sort:           data.Sort,
			DepartmentName: data.DepartmentName,
			Remark:         data.Remark,
			ParentID:       data.ParentId,
			Path:           data.Path,
			Level:          data.Level,
			Children:       buildTChildren(data.Children),
		},
	}, nil
}
