package menu

import (
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/types"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/pointc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
)

func buildBResources(t []*types.Resource) (result []*admin_system.ResourceBody) {
	result = make([]*admin_system.ResourceBody, len(t))
	for i, v := range t {
		result[i] = &admin_system.ResourceBody{
			Method: v.Method,
			Path:   v.Path,
		}
	}
	return
}

func buildTResources(b []*admin_system.ResourceBody) (result []*types.Resource) {
	result = make([]*types.Resource, len(b))
	for i, v := range b {
		result[i] = &types.Resource{
			Method: v.Method,
			Path:   v.Path,
		}
	}
	return
}

func buildTChildren(b []*admin_system.MenuBody) (result []*types.MenuGet) {
	result = make([]*types.MenuGet, len(b))
	for i, v := range b {
		result[i] = &types.MenuGet{
			ID:         v.Id,
			CreatedAt:  v.CreatedAt,
			UpdatedAt:  v.UpdatedAt,
			Sort:       v.Sort,
			Status:     pointc.PStatus32t8(v.Status),
			Title:      v.Title,
			Icon:       v.Icon,
			Code:       v.Code,
			CodePath:   v.CodePath,
			ParentID:   v.ParentId,
			MenuType:   v.MenuType,
			Level:      v.Level,
			Properties: v.Properties,
			Resources:  buildTResources(v.Resources),
			Children:   buildTChildren(v.Children),
		}
	}
	return
}
