package svc

import (
	"github.com/casbin/casbin/v2"
	"github.com/redis/go-redis/v9"
	"github.com/wenlng/go-captcha/v2/rotate"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/config"
	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/middleware"
	admin_system_baseservice "github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/client/baseservice"
	admin_system_casbinservice "github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/client/casbinservice"
	admin_system_departmentservice "github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/client/departmentservice"
	admin_system_menuservice "github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/client/menuservice"
	admin_system_positionservice "github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/client/positionservice"
	admin_system_roleservice "github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/client/roleservice"
	admin_system_userservice "github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/client/userservice"
)

type ServiceContext struct {
	Config  config.Config
	Auth    rest.Middleware
	Redis   redis.UniversalClient
	Casbin  *casbin.Enforcer
	Captcha *rotate.Captcha

	RPCAdminSystem RPCAdminSystem
}

type RPCAdminSystem struct {
	BASE       admin_system_baseservice.BaseService
	CASBIN     admin_system_casbinservice.CasbinService
	DEPARTMENT admin_system_departmentservice.DepartmentService
	MENU       admin_system_menuservice.MenuService
	POSITION   admin_system_positionservice.PositionService
	ROLE       admin_system_roleservice.RoleService
	USER       admin_system_userservice.UserService
}

func NewServiceContext(c config.Config) *ServiceContext {
	r := c.RedisC.MustNewUniversalClient()

	csb := c.CasbinC.MustNewCasbinWithOriginalRedisWatcher(c.DBC.Type, c.DBC.GetDSN(), c.RedisC)
	csb.EnableLog(c.Mode == service.DevMode)

	clientAdminSystem := zrpc.MustNewClient(c.RpcAdminSystem)
	return &ServiceContext{
		Config:  c,
		Auth:    middleware.NewAuthMiddleware(csb, r).Handle,
		Redis:   r,
		Casbin:  csb,
		Captcha: c.CaptchaC.MustNew(),
		RPCAdminSystem: RPCAdminSystem{
			BASE:       admin_system_baseservice.NewBaseService(clientAdminSystem),
			CASBIN:     admin_system_casbinservice.NewCasbinService(clientAdminSystem),
			DEPARTMENT: admin_system_departmentservice.NewDepartmentService(clientAdminSystem),
			MENU:       admin_system_menuservice.NewMenuService(clientAdminSystem),
			POSITION:   admin_system_positionservice.NewPositionService(clientAdminSystem),
			ROLE:       admin_system_roleservice.NewRoleService(clientAdminSystem),
			USER:       admin_system_userservice.NewUserService(clientAdminSystem),
		},
	}
}
