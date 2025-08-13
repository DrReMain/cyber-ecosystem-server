package roleservicelogic

import (
	"strings"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/pointc"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
)

func buildEMenuIDs(b []*admin_system.MenuBody) []string {
	result := make([]string, 0)
	for _, v := range b {
		if v != nil {
			if v.Id != nil && *v.Id != "" {
				result = append(result, *v.Id)
			}
			if v.Children != nil && len(v.Children) > 0 {
				result = append(result, buildEMenuIDs(v.Children)...)
			}
		}
	}
	return result
}

func buildBMenu(e []*ent.Menu) (result []*admin_system.MenuBody) {
	result = make([]*admin_system.MenuBody, len(e))
	for i, v := range e {
		result[i] = &admin_system.MenuBody{
			Id:         pointc.P(v.ID),
			CreatedAt:  nil,
			UpdatedAt:  nil,
			Status:     pointc.PStatus8t32(&v.Status),
			Sort:       nil,
			Title:      pointc.P(v.Title),
			Icon:       nil,
			Code:       pointc.P(v.Code),
			CodePath:   pointc.P(v.CodePath),
			ParentId:   v.ParentID,
			MenuType:   pointc.P(v.MenuType),
			Level:      pointc.P(uint32(len(strings.Split(v.CodePath, "_")))),
			Properties: nil,
			Resources:  nil,
			Children:   nil,
		}
	}
	return
}
