package menuservicelogic

import (
	"errors"
	"sort"
	"strings"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/pointc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/menu"
)

var columns = []string{
	menu.FieldID,
	menu.FieldCreatedAt,
	menu.FieldUpdatedAt,
	menu.FieldSort,
	menu.FieldStatus,
	menu.FieldTitle,
	menu.FieldIcon,
	menu.FieldCode,
	menu.FieldCodePath,
	menu.FieldParentID,
	menu.FieldMenuType,
	menu.FieldMenuPath,
	menu.FieldProperties,
}

func withAlias(C func(string) string, cols []string) (result []string) {
	for _, col := range cols {
		result = append(result, C(col))
	}
	return
}

func checkLevel(path *string) error {
	if path != nil && len(strings.Split(*path, ".")) > 5 {
		return errors.New("level should not more than 5")
	}
	return nil
}

func buildBResources(e []*ent.Resource) (result []*admin_system.ResourceBody) {
	for _, v := range e {
		if v != nil {
			result = append(result, &admin_system.ResourceBody{
				Method: pointc.P(v.Method),
				Path:   pointc.P(v.Path),
			})
		}
	}
	return
}

func sortMenuBySortCreatedAt(n []*admin_system.MenuBody) {
	sort.Slice(n, func(i, j int) bool {
		if n[i].GetSort() != n[j].GetSort() {
			return n[i].GetSort() > n[j].GetSort()
		}
		return n[i].GetCreatedAt() < n[j].GetCreatedAt()
	})
	for _, c := range n {
		sortMenuBySortCreatedAt(c.Children)
	}
}
