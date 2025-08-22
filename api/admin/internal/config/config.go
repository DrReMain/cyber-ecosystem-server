package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/config/captcha"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/config/casbinc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/config/cors"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/config/dbc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/config/redisc"
)

type Config struct {
	rest.RestConf
	Project Project
	CORS    cors.CORSConfig
	Auth    struct {
		AccessSecret  string
		AccessExpire  int64
		RefreshSecret string
		RefreshExpire int64
	}
	RedisC   redisc.RedisConfig
	DBC      dbc.DBConfig
	CasbinC  casbinc.CasbinConfig
	CaptchaC captcha.CaptchaConfig

	RpcAdminSystem zrpc.RpcClientConf
}

type Project struct {
	AppNameHeader string `json:",optional"`
	AppNameValue  string `json:",optional"`
	LangHeader    string `json:",default=X-LANG"`
	DefaultLang   string `json:",env=DEFAULT_LANG"`
}
