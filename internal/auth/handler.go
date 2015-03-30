package auth

import (
	"net/http"

	"golang.org/x/net/context"
	"gopkg.in/triki.v0/internal/ctx"
)

// Handler wraps handlers with context, session and authentication.
// The wrapped function fun should create mongo session (if requred) and save it
// using ctx.SaveSession(...).
func Handler(fun func(context.Context, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, cancel := context.WithTimeout(context.Background(), RequestTimeout)
		cx := ctx.New(c, r)
		defer cancel()
		defer func() {
			s, _ := ctx.DBSessionFromReq(r)
			if s != nil {
				s.Close()
			}
		}()
		// TODO authenticate...
		fun(cx, w, r)
	}
}
