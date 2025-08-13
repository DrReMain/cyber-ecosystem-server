package user

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/helper/common_res"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/svc"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/types"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/pointc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryUserLogic {
	return &QueryUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryUserLogic) QueryUser(req *types.UserQueryReq) (resp *types.UserQueryRes, err error) {
	data, err := l.svcCtx.RPCAdminSystem.USER.QueryUser(l.ctx, &admin_system.UserListReq{
		PageNo:        req.PageNo,
		PageSize:      req.PageSize,
		CreatedAt:     req.CreatedAt,
		UpdatedAt:     req.UpdatedAt,
		Status:        pointc.PStatus8t32(req.Status),
		Email:         req.Email,
		Name:          req.Name,
		Nickname:      req.NickName,
		Phone:         req.Phone,
		DepartmentIds: req.DepartmentIDs,
		PositionIds:   req.PositionIDs,
		RoleIds:       req.RoleIDs,
	})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	resp = &types.UserQueryRes{
		CommonRes: common_res.NewYES(""),
		Result: &types.UserQuery{
			CommonPageRes: &types.CommonPageRes{
				PageNo:   data.PageNo,
				PageSize: data.PageSize,
				Total:    data.Total,
				More:     data.More,
			},
			List: buildTList(data.List),
		},
	}

	return
}
