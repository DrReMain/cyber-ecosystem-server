package departmentservicelogic

import (
	"context"
	"strings"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/pointc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/department"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	"entgo.io/ent/dialect/sql"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetDepartmentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDepartmentLogic {
	return &GetDepartmentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDepartmentLogic) GetDepartment(in *admin_system.IDReq) (*admin_system.DepartmentBody, error) {
	// 根据id，递归查询自身及子Department
	items, err := l.svcCtx.DB.Department.Query().
		Modify(func(s *sql.Selector) {
			t1, t2 := sql.Table(department.Table), sql.Table(department.Table)
			with := sql.WithRecursive("tree")
			with.As(
				sql.Select(withAlias(t1.C, columns)...).
					From(t1).
					Where(sql.EQ(t1.C(department.FieldID), in.Id)).
					UnionAll(
						sql.Select(withAlias(t2.C, columns)...).
							From(t2).
							Join(with).
							On(t2.C(department.FieldParentID), with.C(department.FieldID)),
					),
			)
			s.Prefix(with).Select(withAlias(with.C, columns)...).From(with)
		}).All(l.ctx)
	if err != nil {
		return nil, ent.DefaultHandleError(l.Logger, err, in)
	}

	res := &admin_system.DepartmentBody{}

	m := make(map[string]*admin_system.DepartmentBody)

	for _, v := range items {
		m[v.ID] = &admin_system.DepartmentBody{
			Id:             pointc.P(v.ID),
			CreatedAt:      pointc.P(v.CreatedAt.UnixMilli()),
			UpdatedAt:      pointc.P(v.UpdatedAt.UnixMilli()),
			Sort:           pointc.P(v.Sort),
			DepartmentName: pointc.P(v.DepartmentName),
			Remark:         pointc.P(v.Remark),
			ParentId:       v.ParentID,
			Path:           pointc.P(v.IDPath),
			Level:          pointc.P(uint32(len(strings.Split(v.IDPath, "_")))),
			Children:       []*admin_system.DepartmentBody{},
		}
	}

	for _, n := range m {
		if p, ok := m[n.GetParentId()]; ok {
			p.Children = append(p.Children, n)
		} else {
			res = n
		}
	}

	sortDepartmentBySortCreatedAt(res.Children)
	return res, nil
}
