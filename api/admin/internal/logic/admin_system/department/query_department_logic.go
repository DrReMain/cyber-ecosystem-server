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

type QueryDepartmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryDepartmentLogic {
	return &QueryDepartmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryDepartmentLogic) QueryDepartment(req *types.DepartmentQueryReq) (resp *types.DepartmentQueryRes, err error) {
	data, err := l.svcCtx.RPCAdminSystem.DEPARTMENT.QueryDepartment(l.ctx, &admin_system.DepartmentListReq{
		PageNo:   req.PageNo,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	return &types.DepartmentQueryRes{
		CommonRes: common_res.NewYES(""),
		Data: &types.DepartmentQuery{
			CommonPageRes: &types.CommonPageRes{
				PageNo:   data.PageNo,
				PageSize: data.PageSize,
				Total:    data.Total,
				More:     data.More,
			},
			List: buildTChildren(data.List),
		},
	}, nil
}
