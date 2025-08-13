package baseservicelogic

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/encrypt"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
)

func (l *InitDBLogic) initBaseData() error {
	if err := ent.WithTX(l.ctx, l.svcCtx.DB, func(tx *ent.Tx) error {
		if count, err := tx.Menu.Query().Count(l.ctx); err != nil {
			return err
		} else if count > 0 {
			return errors.New("database is not empty")
		}

		// 读取json
		menuFile, err := os.ReadFile(filepath.Join(l.svcCtx.Config.Project.Etc, "menu.json"))
		if err != nil {
			return err
		}
		roleFile, err := os.ReadFile(filepath.Join(l.svcCtx.Config.Project.Etc, "role.json"))
		if err != nil {
			return err
		}
		userFile, err := os.ReadFile(filepath.Join(l.svcCtx.Config.Project.Etc, "user.json"))
		if err != nil {
			return err
		}

		// 菜单
		var menuList []map[string]any
		if err := json.Unmarshal(menuFile, &menuList); err != nil {
			return err
		}

		menuCodePathMap := make(map[string]*ent.Menu)
		var createMenus func(menus []map[string]any, parentMenu *ent.Menu) error
		createMenus = func(menus []map[string]any, parentMenu *ent.Menu) error {
			for _, menuData := range menus {
				menuCreate := tx.Menu.Create().
					SetTitle(menuData["title"].(string)).
					SetIcon(menuData["icon"].(string)).
					SetCode(menuData["code"].(string)).
					SetCodePath(menuData["code_path"].(string)).
					SetMenuType(menuData["menu_type"].(string))
				if parentMenu != nil {
					menuCreate.SetParentID(parentMenu.ID)
				}

				menu, err := menuCreate.Save(l.ctx)
				if err != nil {
					return err
				}

				menuCodePathMap[menuData["code_path"].(string)] = menu

				if resources, ok := menuData["resources"].([]interface{}); ok && len(resources) > 0 {
					for _, res := range resources {
						resources := res.(map[string]any)
						_, err := tx.Resource.Create().
							SetMethod(resources["method"].(string)).
							SetPath(resources["path"].(string)).
							SetMenuID(menu.ID).
							Save(l.ctx)
						if err != nil {
							return err
						}
					}
				}

				if children, ok := menuData["children"].([]interface{}); ok && len(children) > 0 {
					childrenList := make([]map[string]any, 0, len(children))
					for _, child := range children {
						childrenList = append(childrenList, child.(map[string]any))
					}
					if err := createMenus(childrenList, menu); err != nil {
						return err
					}
				}
			}
			return nil
		}
		if err := createMenus(menuList, nil); err != nil {
			return err
		}

		// 角色
		var roleList []map[string]any
		if err := json.Unmarshal(roleFile, &roleList); err != nil {
			return err
		}

		roleCodeMap := make(map[string]*ent.Role)
		for _, roleData := range roleList {
			role, err := tx.Role.Create().
				SetRoleName(roleData["role_name"].(string)).
				SetCode(roleData["code"].(string)).
				Save(l.ctx)
			if err != nil {
				return err
			}

			roleCodeMap[roleData["code"].(string)] = role

			if menuCodes, ok := roleData["menu_code_path"].([]interface{}); ok {
				menuIds := make([]string, 0)
				for _, codeI := range menuCodes {
					code := codeI.(string)
					if menu, exists := menuCodePathMap[code]; exists {
						menuIds = append(menuIds, menu.ID)
					}
				}

				if len(menuIds) > 0 {
					_, err := tx.Role.UpdateOneID(role.ID).
						AddMenuIDs(menuIds...).
						Save(l.ctx)
					if err != nil {
						return err
					}
				}
			}
		}

		// 用户
		var userList []map[string]any
		if err := json.Unmarshal(userFile, &userList); err != nil {
			return err
		}

		for _, userData := range userList {
			userCreate := tx.User.Create().
				SetPassword(encrypt.EncryptGenerate(userData["password"].(string))).
				SetEmail(userData["email"].(string)).
				SetName(userData["name"].(string)).
				SetNickname(userData["nickname"].(string)).
				SetPhone(userData["phone"].(string)).
				SetAvatar(userData["avatar"].(string))
			user, err := userCreate.Save(l.ctx)
			if err != nil {
				return err
			}

			if roleCodes, ok := userData["role_code"].([]interface{}); ok {
				roleIDs := []string{}
				for _, codeI := range roleCodes {
					code := codeI.(string)
					if role, exists := roleCodeMap[code]; exists {
						roleIDs = append(roleIDs, role.ID)
					}
				}

				if len(roleIDs) > 0 {
					_, err := tx.User.UpdateOneID(user.ID).
						AddRoleIDs(roleIDs...).
						Save(l.ctx)
					if err != nil {
						return err
					}
				}
			}
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}
