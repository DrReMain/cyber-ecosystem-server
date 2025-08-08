package svc

import (
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent"
	_ "github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/ent/runtime"
	"github.com/DrReMain/cyber-ecosystem-server/rpc/admin_system/internal/config"
)

type ServiceContext struct {
	Config config.Config
	DB     *ent.Client
	Redis  redis.UniversalClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := ent.NewClient(
		ent.Log(logx.Info),
		ent.Driver(c.DBC.NewNoCacheDriver()),
	)

	if c.DBC.Debug {
		db = db.Debug()
	}
	return &ServiceContext{
		Config: c,
		DB:     db,
		Redis:  c.RedisC.MustNewUniversalClient(),
	}
}
