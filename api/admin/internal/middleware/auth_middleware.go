package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/config/redisc"
	"github.com/casbin/casbin/v2"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/utils/ctx/role_ctx"
)

type AuthMiddleware struct {
	Casbin *casbin.Enforcer
	Redis  redis.UniversalClient
}

func NewAuthMiddleware(csb *casbin.Enforcer, r redis.UniversalClient) *AuthMiddleware {
	return &AuthMiddleware{
		Casbin: csb,
		Redis:  r,
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roleCode, err := role_ctx.ValueFromCtx(r.Context())
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		}

		jwt := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		v, err := m.Redis.Get(context.Background(), redisc.REDIS_TOKEN_PREFIX+jwt).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			httpx.ErrorCtx(r.Context(), w, errorc.NewHTTPInternal(msgc.SYSTEM_ERROR, err.Error()))
			return

		}
		if v == redisc.REDIS_TOKEN_BANNED {
			httpx.ErrorCtx(
				r.Context(),
				w,
				errorc.NewHTTPUnauthorized(
					msgc.TOKEN_INVALID,
					fmt.Sprintf("token has banned: %s", jwt),
				),
			)
			return
		}

		method, path := r.Method, r.URL.Path
		if result := check(m.Casbin, roleCode, method, path); result {
			next(w, r)
			return
		}

		httpx.ErrorCtx(
			r.Context(),
			w,
			errorc.NewHTTPForbidden(
				msgc.AUTH_NOTALLOWED,
				fmt.Sprintf("roleCode->%v, method->%v, path->%v", roleCode, method, path),
			),
		)
		return
	}
}

func check(csb *casbin.Enforcer, roleCode []string, act, obj string) bool {
	var list = make([][]any, len(roleCode))
	for i, v := range roleCode {
		list[i] = []any{v, act, obj}
	}

	result, err := csb.BatchEnforce(list)
	if err != nil {
		logx.Errorw("[CASBIN]", logx.Field("detail", err.Error()))
		return false
	}

	for _, v := range result {
		if v {
			return true
		}
	}
	return false
}
