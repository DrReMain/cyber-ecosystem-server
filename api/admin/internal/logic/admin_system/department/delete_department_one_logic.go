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

type DeleteDepartmentOneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteDepartmentOneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDepartmentOneLogic {
	return &DeleteDepartmentOneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteDepartmentOneLogic) DeleteDepartmentOne(req *types.DepartmentDeleteReq) (resp *types.DepartmentDeleteRes, err error) {
	data, err := l.svcCtx.RPCAdminSystem.DEPARTMENT.DeleteDepartment(l.ctx, &admin_system.IDsReq{Ids: []string{*req.ID}})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	return &types.DepartmentDeleteRes{
		CommonRes: common_res.NewYES(data.Msg),
	}, nil
}
