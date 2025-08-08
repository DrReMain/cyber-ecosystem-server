package config

import (
	"github.com/zeromicro/go-zero/zrpc"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/config/casbinc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/config/dbc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/config/redisc"
)

type Config struct {
	zrpc.RpcServerConf
	RedisC  redisc.RedisConfig
	DBC     dbc.DBConfig
	CasbinC casbinc.CasbinConfig
	Project Project
}

type Project struct {
	Etc string `json:",env=PROJECT_ETC"`
}
