package departmentservicelogic

import (
	"errors"
	"sort"
	"strings"

	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/department"
)

var columns = []string{
	department.FieldID,
	department.FieldCreatedAt,
	department.FieldUpdatedAt,
	department.FieldSort,
	department.FieldDepartmentName,
	department.FieldRemark,
	department.FieldParentID,
	department.FieldIDPath,
}

func withAlias(C func(string) string, cols []string) (result []string) {
	for _, col := range cols {
		result = append(result, C(col))
	}
	return
}

func checkLevel(path *string) error {
	if path != nil && len(strings.Split(*path, "_")) > 10 {
		return errors.New("level should not more then 10")
	}
	return nil
}

func buildEIDs(e []*ent.Department) (result []any) {
	result = make([]any, len(e))
	for i, v := range e {
		result[i] = v.ID
	}
	return
}

func sortDepartmentBySortCreatedAt(n []*admin_system.DepartmentBody) {
	sort.Slice(n, func(i, j int) bool {
		if n[i].GetSort() != n[j].GetSort() {
			return n[i].GetSort() > n[j].GetSort() // DESC
		}
		return n[i].GetCreatedAt() < n[j].GetCreatedAt() // ASC
	})
	for _, c := range n {
		sortDepartmentBySortCreatedAt(c.Children)
	}
}
