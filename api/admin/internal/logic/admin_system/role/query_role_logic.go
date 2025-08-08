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

type QueryRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryRoleLogic {
	return &QueryRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryRoleLogic) QueryRole(req *types.RoleQueryReq) (resp *types.RoleQueryRes, err error) {
	data, err := l.svcCtx.RPCAdminSystem.ROLE.QueryRole(l.ctx, &admin_system.RoleListReq{
		PageNo:    req.PageNo,
		PageSize:  req.PageSize,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
		RoleName:  req.RoleName,
		Code:      req.Code,
	})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	resp = &types.RoleQueryRes{
		CommonRes: common_res.NewYES(""),
		Data: &types.RoleQuery{
			CommonPageRes: &types.CommonPageRes{
				PageNo:   data.PageNo,
				PageSize: data.PageSize,
				Total:    data.Total,
				More:     data.More,
			},
			List: make([]*types.RoleGet, len(data.List)),
		},
	}

	for i, v := range data.List {
		resp.Data.List[i] = &types.RoleGet{
			ID:        v.Id,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			Sort:      v.Sort,
			RoleName:  v.RoleName,
			Code:      v.Code,
			Remark:    v.Remark,
			MenuIds:   buildTMenuIDs(v.Menus),
		}
	}
	return
}
