package position

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/helper/common_res"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/svc"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/types"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPositionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPositionLogic {
	return &GetPositionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPositionLogic) GetPosition(req *types.PositionGetReq) (resp *types.PositionGetRes, err error) {
	data, err := l.svcCtx.RPCAdminSystem.POSITION.GetPosition(l.ctx, &admin_system.IDReq{Id: *req.ID})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	return &types.PositionGetRes{
		CommonRes: common_res.NewYES(""),
		Data: &types.PositionGet{
			ID:           data.Id,
			CreatedAt:    data.CreatedAt,
			UpdatedAt:    data.UpdatedAt,
			Sort:         data.Sort,
			PositionName: data.PositionName,
			Code:         data.Code,
			Remark:       data.Remark,
		},
	}, nil
}
