package casbin_rules

import (
	"errors"

	"github.com/casbin/casbin/v2"

	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
)

func RefreshCasbinRules(csb *casbin.Enforcer, codes []string, b []*admin_system.CasbinBody) error {
	oldPolicies, err := csb.GetFilteredPolicy(0, codes...)
	if err != nil {
		return err
	}

	if len(oldPolicies) != 0 {
		if removeResult, err := csb.RemoveFilteredPolicy(0, codes...); err != nil {
			return err
		} else if !removeResult {
			return errors.New("casbin rules remove failed")
		}
	}

	if b == nil || len(b) == 0 {
		return nil
	}

	var policies = make([][]string, len(b))
	for i, v := range b {
		policies[i] = []string{v.RoleCode, v.Method, v.Path}
	}

	if addResult, err := csb.AddPolicies(policies); err != nil {
		return err
	} else if !addResult {
		return errors.New("casbin rules add failed")
	}
	return nil
}
