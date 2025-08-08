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

type CreatePositionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePositionLogic {
	return &CreatePositionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePositionLogic) CreatePosition(req *types.PositionCreateReq) (resp *types.PositionCreateRes, err error) {
	data, err := l.svcCtx.RPCAdminSystem.POSITION.CreatePosition(l.ctx, &admin_system.PositionBody{
		Sort:         req.Sort,
		PositionName: req.PositionName,
		Code:         req.Code,
		Remark:       req.Remark,
	})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	return &types.PositionCreateRes{
		CommonRes: common_res.NewYES(data.Msg),
		Data:      &data.Id,
	}, nil
}
