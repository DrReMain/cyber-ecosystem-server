package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
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
				errors.New(fmt.Sprintf("'%s' in header is invalid", m.Header)),
			)
			return
		}
		next.ServeHTTP(w, r)
	})
}
