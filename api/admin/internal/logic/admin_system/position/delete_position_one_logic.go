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

type DeletePositionOneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletePositionOneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePositionOneLogic {
	return &DeletePositionOneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletePositionOneLogic) DeletePositionOne(req *types.PositionDeleteReq) (resp *types.PositionDeleteRes, err error) {
	data, err := l.svcCtx.RPCAdminSystem.POSITION.DeletePosition(l.ctx, &admin_system.IDsReq{Ids: []string{*req.ID}})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	return &types.PositionDeleteRes{
		CommonRes: common_res.NewYES(data.Msg),
	}, nil
}
