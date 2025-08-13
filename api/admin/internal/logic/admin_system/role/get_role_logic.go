package role

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/helper/common_res"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/svc"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/types"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleLogic {
	return &GetRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRoleLogic) GetRole(req *types.RoleGetReq) (resp *types.RoleGetRes, err error) {
	data, err := l.svcCtx.RPCAdminSystem.ROLE.GetRole(l.ctx, &admin_system.IDReq{Id: *req.ID})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	return &types.RoleGetRes{
		CommonRes: common_res.NewYES(""),
		Result: &types.RoleGet{
			ID:        data.Id,
			CreatedAt: data.CreatedAt,
			UpdatedAt: data.UpdatedAt,
			Sort:      data.Sort,
			RoleName:  data.RoleName,
			Code:      data.Code,
			Remark:    data.Remark,
			MenuIds:   buildTMenuIDs(data.Menus),
		},
	}, nil
}
