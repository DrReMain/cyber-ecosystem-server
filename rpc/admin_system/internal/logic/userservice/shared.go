package userservicelogic

import (
	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/pointc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
)

func buildBDepartment(e *ent.Department) *admin_system.DepartmentBody {
	if e == nil {
		return nil
	}
	return &admin_system.DepartmentBody{
		Id:             pointc.P(e.ID),
		CreatedAt:      nil,
		UpdatedAt:      nil,
		Sort:           nil,
		DepartmentName: pointc.P(e.DepartmentName),
		Remark:         nil,
		ParentId:       pointc.P(e.ParentID),
		Path:           pointc.P(e.IDPath),
		Level:          nil,
		Children:       nil,
	}
}

func buildBPositions(e []*ent.Position) (result []*admin_system.PositionBody) {
	result = make([]*admin_system.PositionBody, len(e))
	for i, v := range e {
		result[i] = &admin_system.PositionBody{
			Id:           pointc.P(v.ID),
			CreatedAt:    nil,
			UpdatedAt:    nil,
			Sort:         nil,
			PositionName: pointc.P(v.PositionName),
			Code:         pointc.P(v.Code),
			Remark:       nil,
		}
	}

	return
}

func buildBRoles(e []*ent.Role) (result []*admin_system.RoleBody) {
	result = make([]*admin_system.RoleBody, len(e))
	for i, v := range e {
		result[i] = &admin_system.RoleBody{
			Id:        pointc.P(v.ID),
			CreatedAt: nil,
			UpdatedAt: nil,
			Sort:      nil,
			RoleName:  pointc.P(v.RoleName),
			Code:      pointc.P(v.Code),
			Remark:    pointc.P(v.Remark),
		}
	}
	return
}
