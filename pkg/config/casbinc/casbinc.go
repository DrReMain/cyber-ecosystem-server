package casbinc

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	entadapter "github.com/casbin/ent-adapter"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	zeroredis "github.com/zeromicro/go-zero/core/stores/redis"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/config/redisc"
)

type CasbinConfig struct {
	ModelText string `json:"ModelText,optional,env=CASBIN_MODEL_TEXT"`
}

func (c CasbinConfig) NewCasbin(dbType, dsn string) (*casbin.Enforcer, error) {
	adapter, err := entadapter.NewAdapter(dbType, dsn)
	if err != nil {
		logx.Must(err)
	}

	var text string
	if c.ModelText == "" {
		text = `
		[request_definition]
		r = sub, act, obj
		
		[policy_definition]
		p = sub, act, obj
		
		[role_definition]
		g = _, _
		
		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = r.sub == p.sub && r.act == p.act && keyMatch2(r.obj,p.obj)
		`
	} else {
		text = c.ModelText
	}

	m, err := model.NewModelFromString(text)
	if err != nil {
		logx.Must(err)
	}

	enforce, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		logx.Must(err)
	}

	if err := enforce.LoadPolicy(); err != nil {
		logx.Must(err)
	}

	return enforce, nil
}

func (c CasbinConfig) MustNewCasbin(dbType, dsn string) *casbin.Enforcer {
	if csb, err := c.NewCasbin(dbType, dsn); err != nil {
		logx.Errorw("initialize Casbin failed", logx.Field("detail", err.Error()))
		log.Fatalf("initialize Casbin failed, %s", err.Error())
		return nil
	} else {
		return csb
	}
}

func (c CasbinConfig) MustNewRedisWatcher(rc zeroredis.RedisConf, f func(s string)) persist.Watcher {
	w, err := rediswatcher.NewWatcher(rc.Host, rediswatcher.WatcherOptions{
		Options: redis.Options{
			Network:  "tcp",
			Password: rc.Pass,
		},
		Channel:    redisc.REDIS_CASBIN,
		IgnoreSelf: false,
	})
	if err != nil {
		logx.Must(err)
	}
	if err = w.SetUpdateCallback(f); err != nil {
		logx.Must(err)
	}
	return w
}

func (c CasbinConfig) MustNewOriginalRedisWatcher(rc redisc.RedisConfig, f func(s string)) persist.Watcher {
	w, err := rediswatcher.NewWatcher(rc.Host, rediswatcher.WatcherOptions{
		Options: redis.Options{
			Network:  "tcp",
			Username: rc.Username,
			Password: rc.Password,
		},
		Channel:    fmt.Sprintf("%s-%d", redisc.REDIS_CASBIN, rc.DB),
		IgnoreSelf: false,
	})
	if err != nil {
		logx.Must(err)
	}
	if err = w.SetUpdateCallback(f); err != nil {
		logx.Must(err)
	}
	return w
}

func (c CasbinConfig) MustNewCasbinWithRedisWatcher(dbType, dsn string, rc zeroredis.RedisConf) *casbin.Enforcer {
	csb := c.MustNewCasbin(dbType, dsn)
	w := c.MustNewRedisWatcher(rc, func(s string) {
		rediswatcher.DefaultUpdateCallback(csb)(s)
	})
	if err := csb.SetWatcher(w); err != nil {
		logx.Must(err)
	}
	if err := csb.SavePolicy(); err != nil {
		logx.Must(err)
	}
	return csb
}

func (c CasbinConfig) MustNewCasbinWithOriginalRedisWatcher(dbType, dsn string, rc redisc.RedisConfig) *casbin.Enforcer {
	csb := c.MustNewCasbin(dbType, dsn)
	w := c.MustNewOriginalRedisWatcher(rc, func(s string) {
		rediswatcher.DefaultUpdateCallback(csb)(s)
	})
	if err := csb.SetWatcher(w); err != nil {
		logx.Must(err)
	}
	if err := csb.SavePolicy(); err != nil {
		logx.Must(err)
	}
	return csb
}
