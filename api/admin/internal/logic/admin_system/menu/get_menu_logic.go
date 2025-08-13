package menu

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

type GetMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuLogic {
	return &GetMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMenuLogic) GetMenu(req *types.MenuGetReq) (resp *types.MenuGetRes, err error) {
	data, err := l.svcCtx.RPCAdminSystem.MENU.GetMenu(l.ctx, &admin_system.IDReq{Id: *req.ID})
	if err != nil {
		return nil, errorc.NewGRPCError(err)
	}

	return &types.MenuGetRes{
		CommonRes: common_res.NewYES(""),
		Result: &types.MenuGet{
			ID:         data.Id,
			CreatedAt:  data.CreatedAt,
			UpdatedAt:  data.UpdatedAt,
			Sort:       data.Sort,
			Status:     pointc.PStatus32t8(data.Status),
			Title:      data.Title,
			Icon:       data.Icon,
			Code:       data.Code,
			CodePath:   data.CodePath,
			ParentID:   data.ParentId,
			MenuType:   data.MenuType,
			Level:      data.Level,
			Properties: data.Properties,
			Resources:  buildTResources(data.Resources),
			Children:   buildTChildren(data.Children),
		},
	}, nil
}
