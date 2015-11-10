package middleware

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/zenazn/goji/web"
)

type CSRF struct {
	ProtectKey string
}

func (m *CSRF) MiddlewareFunc() interface{} {
	key := []byte(m.ProtectKey)
	return func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
		return csrf.Protect(key)(http.HandlerFunc(fn))
	}
}
