package user

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/helper/common_res"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/svc"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/types"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/pointc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(req *types.UserUpdateReq) (resp *types.UserUpdateRes, err error) {
	if req.Password != nil && req.Confirm != nil && *req.Password != *req.Confirm {
		return nil, errorc.NewHTTPBadRequest(msgc.CONFIRM_ERROR)
	}

	data, err := l.svcCtx.RPCAdminSystem.USER.UpdateUser(l.ctx, &admin_system.UserBody{
		Id:       req.ID,
		Status:   pointc.PStatus8t32(req.Status),
		Password: req.Password,
		Email:    req.Email,
		Name:     req.Name,
		Nickname: req.NickName,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		Remark:   req.Remark,
		Department: &admin_system.DepartmentBody{
			Id: req.DepartmentID,
		},
		Positions: buildBPosition(req.PositionIDs),
		Roles:     buildBRole(req.RoleIDs),
	})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	return &types.UserUpdateRes{
		CommonRes: common_res.NewYES(data.Msg),
	}, nil
}
