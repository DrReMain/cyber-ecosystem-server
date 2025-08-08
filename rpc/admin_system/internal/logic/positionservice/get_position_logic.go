package positionservicelogic

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/pointc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPositionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPositionLogic {
	return &GetPositionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPositionLogic) GetPosition(in *admin_system.IDReq) (*admin_system.PositionBody, error) {
	item, err := l.svcCtx.DB.Position.Get(l.ctx, in.Id)
	if err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	return &admin_system.PositionBody{
		Id:           pointc.P(item.ID),
		CreatedAt:    pointc.P(item.CreatedAt.UnixMilli()),
		UpdatedAt:    pointc.P(item.UpdatedAt.UnixMilli()),
		Sort:         pointc.P(item.Sort),
		PositionName: pointc.P(item.PositionName),
		Code:         pointc.P(item.Code),
		Remark:       pointc.P(item.Remark),
	}, nil
}
