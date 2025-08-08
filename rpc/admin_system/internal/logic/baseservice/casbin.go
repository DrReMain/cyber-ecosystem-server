package baseservicelogic

import (
	"errors"

	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/role"
)

func (l *InitDBLogic) initCasbin() error {
	var policies [][]string
	if r, err := l.svcCtx.DB.Role.Query().
		Where(role.CodeIn("SUPER", "ADMIN")).
		WithMenus(func(query *ent.MenuQuery) {
			query = query.WithResources()
		}).
		All(l.ctx); err != nil {
		return err
	} else {
		for _, v1 := range r {
			var resources = make(map[string][]string)
			for _, v2 := range v1.Edges.Menus {
				for _, v3 := range v2.Edges.Resources {
					resources[v1.Code+v3.Method+v3.Path] = []string{v1.Code, v3.Method, v3.Path}
				}
			}
			for _, r := range resources {
				policies = append(policies, []string{r[0], r[1], r[2]})
			}
		}
	}

	csb, err := l.svcCtx.Config.CasbinC.NewCasbin(l.svcCtx.Config.DBC.Type, l.svcCtx.Config.DBC.GetDSN())
	if err != nil {
		return err
	}

	result, err := csb.AddPolicies(policies)
	if err != nil {
		return err
	}
	if !result {
		return errors.New("casbin add policies failed")
	}

	return nil
}
