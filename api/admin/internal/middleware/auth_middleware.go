package middleware

import (
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/helper/usual_err"
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
			httpx.Error(w, err)
		}
		method, path := r.Method, r.URL.Path

		// TODO: jwt blacklist

		if result := check(m.Casbin, roleCode, method, path); result {
			next(w, r)
			return
		}

		httpx.Error(
			w,
			usual_err.HTTPForbidden(fmt.Sprintf("roleCode->%v, method->%v, path->%v", roleCode, method, path)),
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
