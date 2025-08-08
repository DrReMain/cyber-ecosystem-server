package positionservicelogic

import (
	"context"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/pointc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/position"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/predicate"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryPositionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryPositionLogic {
	return &QueryPositionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryPositionLogic) QueryPosition(in *admin_system.PositionListReq) (*admin_system.PositionListRes, error) {
	items, err := l.svcCtx.DB.Position.Query().
		Where(
			ent.NewPredicatePosition().ApplyCreatedAt(in.CreatedAt).ApplyUpdatedAt(in.UpdatedAt).
				Apply(in.PositionName != nil, func() predicate.Position {
					return position.PositionNameContains(*in.PositionName)
				}).
				Apply(in.Code != nil, func() predicate.Position {
					return position.CodeEQ(*in.Code)
				}).
				Submit()...,
		).
		Page(l.ctx, in.PageNo, in.PageSize, 0, 0, func(pager *ent.PositionPager) {
			pager.Order = []position.OrderOption{
				ent.Desc(position.FieldSort),
				ent.Asc(position.FieldCreatedAt),
			}
		})
	if err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	res := &admin_system.PositionListRes{
		PageNo:   items.PageDetail.PageNo,
		PageSize: items.PageDetail.PageSize,
		Total:    items.PageDetail.Total,
		More:     items.PageDetail.More,
	}

	for _, v := range items.List {
		res.List = append(res.List, &admin_system.PositionBody{
			Id:           pointc.P(v.ID),
			CreatedAt:    pointc.P(v.CreatedAt.UnixMilli()),
			UpdatedAt:    pointc.P(v.UpdatedAt.UnixMilli()),
			Sort:         pointc.P(v.Sort),
			PositionName: pointc.P(v.PositionName),
			Code:         pointc.P(v.Code),
			Remark:       pointc.P(v.Remark),
		})
	}

	return res, nil
}
