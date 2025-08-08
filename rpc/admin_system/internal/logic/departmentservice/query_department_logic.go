package departmentservicelogic

import (
	"context"
	"strings"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/pointc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/department"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/predicate"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"entgo.io/ent/dialect/sql"
	"github.com/zeromicro/go-zero/core/logx"
)

type QueryDepartmentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryDepartmentLogic {
	return &QueryDepartmentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryDepartmentLogic) QueryDepartment(in *admin_system.DepartmentListReq) (*admin_system.DepartmentListRes, error) {
	// 根据page_no, page_size，查询所有根Department
	items, err := l.svcCtx.DB.Department.Query().
		Where(
			ent.NewPredicateDepartment().
				Apply(true, func() predicate.Department {
					return department.ParentIDEQ("") // 查询一级部门
				}).
				Submit()...,
		).
		Page(l.ctx, in.PageNo, in.PageSize, 0, 0, func(pager *ent.DepartmentPager) {
			pager.Order = []department.OrderOption{
				ent.Desc(department.FieldSort),
				ent.Asc(department.FieldCreatedAt),
			}
		})
	if err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	// 根据根Department的ids，查询所有子Department
	children, err := l.svcCtx.DB.Department.Query().
		Modify(func(s *sql.Selector) {
			t1, t2 := sql.Table(department.Table), sql.Table(department.Table)
			with := sql.WithRecursive("tree")
			with.As(
				sql.Select(withAlias(t1.C, columns)...).
					From(t1).
					Where(sql.In(t1.C(department.FieldParentID), buildEIDs(items.List)...)).
					UnionAll(
						sql.Select(withAlias(t2.C, columns)...).
							From(t2).
							Join(with).
							On(t2.C(department.FieldParentID), with.C(department.FieldID)),
					),
			)
			s.Prefix(with).Select(withAlias(with.C, columns)...).From(with)
		}).
		All(l.ctx)
	if err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	m := make(map[string]*admin_system.DepartmentBody)

	for _, v := range append(items.List, children...) {
		m[v.ID] = &admin_system.DepartmentBody{
			Id:             pointc.P(v.ID),
			CreatedAt:      pointc.P(v.CreatedAt.UnixMilli()),
			UpdatedAt:      pointc.P(v.UpdatedAt.UnixMilli()),
			Sort:           pointc.P(v.Sort),
			DepartmentName: pointc.P(v.DepartmentName),
			Remark:         pointc.P(v.Remark),
			ParentId:       pointc.P(v.ParentID),
			Path:           pointc.P(v.IDPath),
			Level:          pointc.P(uint32(len(strings.Split(v.IDPath, ".")))),
			Children:       []*admin_system.DepartmentBody{},
		}
	}

	res := &admin_system.DepartmentListRes{
		PageNo:   items.PageDetail.PageNo,
		PageSize: items.PageDetail.PageSize,
		Total:    items.PageDetail.Total,
		More:     items.PageDetail.More,
	}

	for _, n := range m {
		if p, ok := m[n.GetParentId()]; ok {
			p.Children = append(p.Children, n)
		} else {
			res.List = append(res.List, n)
		}
	}

	sortDepartmentBySortCreatedAt(res.List)
	return res, nil
}
