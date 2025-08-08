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

type UpdatePositionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePositionLogic {
	return &UpdatePositionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePositionLogic) UpdatePosition(req *types.PositionUpdateReq) (resp *types.PositionUpdateRes, err error) {
	data, err := l.svcCtx.RPCAdminSystem.POSITION.UpdatePosition(l.ctx, &admin_system.PositionBody{
		Id:           req.ID,
		Sort:         req.Sort,
		PositionName: req.PositionName,
		Code:         req.Code,
		Remark:       req.Remark,
	})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	return &types.PositionUpdateRes{
		CommonRes: common_res.NewYES(data.Msg),
	}, nil
}
