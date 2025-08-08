package casbinservicelogic

import (
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
)

func buildBList(e []*ent.Role) (result []*admin_system.CasbinBody) {
	for _, v1 := range e {
		var resources = make(map[string][]string)
		for _, v2 := range v1.Edges.Menus {
			for _, v3 := range v2.Edges.Resources {
				resources[v1.Code+v3.Method+v3.Path] = []string{v1.Code, v3.Method, v3.Path}
			}
		}
		for _, r := range resources {
			result = append(result, &admin_system.CasbinBody{
				RoleCode: r[0],
				Method:   r[1],
				Path:     r[2],
			})
		}
	}
	return
}
