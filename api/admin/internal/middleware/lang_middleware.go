package middleware

import (
	"context"
	"net/http"
)

type LangMiddleware struct {
	Header string
}

func NewLangMiddleware(header string) *LangMiddleware {
	return &LangMiddleware{
		Header: header,
	}
}

func (m *LangMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value := r.Header.Get(m.Header)
		ctx := context.WithValue(r.Context(), m.Header, value)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
