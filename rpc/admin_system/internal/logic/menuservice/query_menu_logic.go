package menuservicelogic

import (
	"context"
	"strings"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/pointc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/menu"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryMenuLogic {
	return &QueryMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryMenuLogic) QueryMenu(in *admin_system.MenuListReq) (*admin_system.MenuListRes, error) {
	items, err := l.svcCtx.DB.Menu.Query().
		WithResources().
		Page(l.ctx, nil, nil, 0, 0, func(pager *ent.MenuPager) {
			pager.Order = []menu.OrderOption{
				ent.Desc(menu.FieldSort),
				ent.Asc(menu.FieldCreatedAt),
			}
		})
	if err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	m := make(map[string]*admin_system.MenuBody)

	for _, v := range items.List {
		m[v.ID] = &admin_system.MenuBody{
			Id:         pointc.P(v.ID),
			CreatedAt:  pointc.P(v.CreatedAt.UnixMilli()),
			UpdatedAt:  pointc.P(v.UpdatedAt.UnixMilli()),
			Status:     pointc.PStatus8t32(&v.Status),
			Sort:       pointc.P(v.Sort),
			Title:      pointc.P(v.Title),
			Icon:       pointc.P(v.Icon),
			Code:       pointc.P(v.Code),
			CodePath:   pointc.P(v.CodePath),
			ParentId:   pointc.P(v.ParentID),
			MenuType:   pointc.P(v.MenuType),
			MenuPath:   pointc.P(v.MenuPath),
			Level:      pointc.P(uint32(len(strings.Split(v.CodePath, ".")))),
			Properties: pointc.P(v.Properties),
			Resources:  buildBResources(v.Edges.Resources),
			Children:   []*admin_system.MenuBody{},
		}
	}

	res := &admin_system.MenuListRes{}

	for _, n := range m {
		if p, ok := m[n.GetParentId()]; ok {
			p.Children = append(p.Children, n)
		} else {
			res.List = append(res.List, n)
		}
	}
	sortMenuBySortCreatedAt(res.List)
	return res, nil
}
