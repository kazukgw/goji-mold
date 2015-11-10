package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/zenazn/goji/web"
)

type Session struct {
	Name          string
	Key           string
	EncriptionKey string
}

func (m *Session) MiddlewareFunc() interface{} {
	store := sessions.NewCookieStore(
		[]byte(m.Key),
		[]byte(m.EncriptionKey),
	)
	return func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			sess, err := store.Get(r, m.Name)
			if err != nil {
				panic(err.Error())
			}
			c.Env["session"] = sess
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
