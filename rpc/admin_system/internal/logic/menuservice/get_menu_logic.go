package menuservicelogic

import (
	"context"
	"strings"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/pointc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/menu"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"entgo.io/ent/dialect/sql"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuLogic {
	return &GetMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMenuLogic) GetMenu(in *admin_system.IDReq) (*admin_system.MenuBody, error) {
	items, err := l.svcCtx.DB.Menu.Query().
		Modify(func(s *sql.Selector) {
			t1, t2 := sql.Table(menu.Table), sql.Table(menu.Table)
			with := sql.WithRecursive("tree")
			with.As(
				sql.Select(withAlias(t1.C, columns)...).
					From(t1).
					Where(sql.EQ(t1.C(menu.FieldID), in.Id)).
					UnionAll(
						sql.Select(withAlias(t2.C, columns)...).
							From(t2).
							Join(with).
							On(t2.C(menu.FieldParentID), with.C(menu.FieldID)),
					),
			)
			s.Prefix(with).Select(withAlias(with.C, columns)...).From(with)
		}).
		WithResources().
		All(l.ctx)
	if err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	res := &admin_system.MenuBody{}

	m := make(map[string]*admin_system.MenuBody)

	for _, v := range items {
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

	for _, n := range m {
		if p, ok := m[n.GetParentId()]; ok {
			p.Children = append(p.Children, n)
		} else {
			res = n
		}
	}

	sortMenuBySortCreatedAt(res.Children)
	return res, nil
}
