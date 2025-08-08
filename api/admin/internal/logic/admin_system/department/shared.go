package department

import (
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/types"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
)

func buildTChildren(b []*admin_system.DepartmentBody) (result []*types.DepartmentGet) {
	result = make([]*types.DepartmentGet, len(b))
	for i, v := range b {
		result[i] = &types.DepartmentGet{
			ID:             v.Id,
			CreatedAt:      v.CreatedAt,
			UpdatedAt:      v.UpdatedAt,
			Sort:           v.Sort,
			DepartmentName: v.DepartmentName,
			Remark:         v.Remark,
			ParentID:       v.ParentId,
			Path:           v.Path,
			Level:          v.Level,
			Children:       buildTChildren(v.Children),
		}
	}
	return
}
