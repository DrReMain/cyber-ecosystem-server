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

type QueryPositionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryPositionLogic {
	return &QueryPositionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryPositionLogic) QueryPosition(req *types.PositionQueryReq) (resp *types.PositionQueryRes, err error) {
	data, err := l.svcCtx.RPCAdminSystem.POSITION.QueryPosition(l.ctx, &admin_system.PositionListReq{
		PageNo:       req.PageNo,
		PageSize:     req.PageSize,
		CreatedAt:    req.CreatedAt,
		UpdatedAt:    req.UpdatedAt,
		PositionName: req.PositionName,
		Code:         req.Code,
	})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	resp = &types.PositionQueryRes{
		CommonRes: common_res.NewYES(""),
		Result: &types.PositionQuery{
			CommonPageRes: &types.CommonPageRes{
				PageNo:   data.PageNo,
				PageSize: data.PageSize,
				Total:    data.Total,
				More:     data.More,
			},
			List: make([]*types.PositionGet, len(data.List)),
		},
	}

	for i, v := range data.List {
		resp.Result.List[i] = &types.PositionGet{
			ID:           v.Id,
			CreatedAt:    v.CreatedAt,
			UpdatedAt:    v.UpdatedAt,
			Sort:         v.Sort,
			PositionName: v.PositionName,
			Code:         v.Code,
			Remark:       v.Remark,
		}
	}

	return
}
