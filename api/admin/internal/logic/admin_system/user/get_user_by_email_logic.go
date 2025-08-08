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

type GetUserByEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserByEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByEmailLogic {
	return &GetUserByEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserByEmailLogic) GetUserByEmail(req *types.UserGetByEmailReq) (resp *types.UserGetRes, err error) {
	data, err := l.svcCtx.RPCAdminSystem.USER.GetUserByEmail(l.ctx, &admin_system.EmailReq{Email: *req.Email})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	return &types.UserGetRes{
		CommonRes: common_res.NewYES(""),
		Data: &types.UserGet{
			ID:         data.Id,
			CreatedAt:  data.CreatedAt,
			UpdatedAt:  data.UpdatedAt,
			Status:     pointc.PStatus32t8(data.Status),
			Email:      data.Email,
			Name:       data.Name,
			NickName:   data.Nickname,
			Phone:      data.Phone,
			Avatar:     data.Avatar,
			Remark:     data.Remark,
			Department: buildTDepartment(data.Department),
			Positions:  buildTPositions(data.Positions),
			Roles:      buildTRoles(data.Roles),
		},
	}, nil
}
