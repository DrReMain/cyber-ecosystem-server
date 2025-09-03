package middleware

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
)

type AppNameMiddleware struct {
	Header string
	Value  []string
}

func NewAppNameMiddleware(header string, value []string) *AppNameMiddleware {
	return &AppNameMiddleware{
		Header: header,
		Value:  value,
	}
}

func (m *AppNameMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get(m.Header)
		if len(m.Value) > 0 && !includes(m.Value, h) {
			httpx.ErrorCtx(
				r.Context(),
				w,
				errorc.NewHTTPForbidden(msgc.APP_NOTALLOWED, "app name check fail"),
			)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func includes(slice []string, target string) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}
