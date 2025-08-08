package main

import (
	"flag"
	"fmt"
	casbinserviceServer "github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/server/casbinservice"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/admin_system"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/config"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/svc"

	baseserviceServer "github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/server/baseservice"
	departmentserviceServer "github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/server/departmentservice"
	menuserviceServer "github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/server/menuservice"
	positionserviceServer "github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/server/positionservice"
	roleserviceServer "github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/server/roleservice"
	userserviceServer "github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/server/userservice"
)

var configFile = flag.String("f", "etc/admin_system.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		admin_system.RegisterBaseServiceServer(grpcServer, baseserviceServer.NewBaseServiceServer(ctx))
		admin_system.RegisterCasbinServiceServer(grpcServer, casbinserviceServer.NewCasbinServiceServer(ctx))
		admin_system.RegisterDepartmentServiceServer(grpcServer, departmentserviceServer.NewDepartmentServiceServer(ctx))
		admin_system.RegisterMenuServiceServer(grpcServer, menuserviceServer.NewMenuServiceServer(ctx))
		admin_system.RegisterPositionServiceServer(grpcServer, positionserviceServer.NewPositionServiceServer(ctx))
		admin_system.RegisterRoleServiceServer(grpcServer, roleserviceServer.NewRoleServiceServer(ctx))
		admin_system.RegisterUserServiceServer(grpcServer, userserviceServer.NewUserServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
