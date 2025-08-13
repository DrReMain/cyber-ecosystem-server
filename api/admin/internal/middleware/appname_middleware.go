package middleware

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
)

type AppNameMiddleware struct {
	Header string
	Value  string
}

func NewAppNameMiddleware(header, value string) *AppNameMiddleware {
	return &AppNameMiddleware{
		Header: header,
		Value:  value,
	}
}

func (m *AppNameMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.Header != "" && r.Header.Get(m.Header) != m.Value {
			httpx.Error(
				w,
				errorc.NewHTTPForbidden(msgc.APP_NOTALLOWED, "app name check fail"),
			)
			return
		}
		next.ServeHTTP(w, r)
	})
}
