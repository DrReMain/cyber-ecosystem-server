package role

import "github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"

func buildTMenuIDs(b []*admin_system.MenuBody) (result []string) {
	result = make([]string, len(b))
	for i, v := range b {
		result[i] = *v.Id
	}
	return
}

func buildBMenus(ids []string) (result []*admin_system.MenuBody) {
	result = make([]*admin_system.MenuBody, len(ids))
	for i, v := range ids {
		result[i] = &admin_system.MenuBody{Id: &v}
	}
	return
}
