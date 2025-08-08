package user

import (
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/types"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/pointc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
)

func buildBPosition(ids []string) (result []*admin_system.PositionBody) {
	result = make([]*admin_system.PositionBody, len(ids))
	for i, id := range ids {
		result[i] = &admin_system.PositionBody{Id: &id}
	}
	return
}

func buildBRole(ids []string) (result []*admin_system.RoleBody) {
	result = make([]*admin_system.RoleBody, len(ids))
	for i, id := range ids {
		result[i] = &admin_system.RoleBody{Id: &id}
	}
	return
}

func buildTDepartment(body *admin_system.DepartmentBody) *types.DepartmentGet {
	if body == nil {
		return nil
	}
	return &types.DepartmentGet{
		ID:             body.Id,
		CreatedAt:      nil,
		UpdatedAt:      nil,
		Sort:           nil,
		DepartmentName: body.DepartmentName,
		Remark:         nil,
		ParentID:       body.ParentId,
		Path:           body.Path,
		Level:          nil,
		Children:       nil,
	}
}

func buildTPositions(b []*admin_system.PositionBody) (result []*types.PositionGet) {
	result = make([]*types.PositionGet, len(b))
	for i, v := range b {
		result[i] = &types.PositionGet{
			ID:           v.Id,
			CreatedAt:    nil,
			UpdatedAt:    nil,
			Sort:         nil,
			PositionName: v.PositionName,
			Code:         v.Code,
			Remark:       nil,
		}
	}
	return
}

func buildTRoles(b []*admin_system.RoleBody) (result []*types.RoleGet) {
	result = make([]*types.RoleGet, len(b))
	for i, v := range b {
		result[i] = &types.RoleGet{
			ID:        v.Id,
			CreatedAt: nil,
			UpdatedAt: nil,
			Sort:      nil,
			RoleName:  v.RoleName,
			Code:      v.Code,
			Remark:    nil,
		}
	}
	return
}

func buildTList(b []*admin_system.UserBody) (result []*types.UserGet) {
	result = make([]*types.UserGet, len(b))
	for i, v := range b {
		result[i] = &types.UserGet{
			ID:         v.Id,
			CreatedAt:  v.CreatedAt,
			UpdatedAt:  v.UpdatedAt,
			Status:     pointc.PStatus32t8(v.Status),
			Email:      v.Email,
			Name:       v.Name,
			NickName:   v.Nickname,
			Phone:      v.Phone,
			Avatar:     v.Avatar,
			Remark:     v.Remark,
			Department: buildTDepartment(v.Department),
			Positions:  buildTPositions(v.Positions),
			Roles:      buildTRoles(v.Roles),
		}
	}
	return
}
