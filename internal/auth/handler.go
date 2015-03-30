package auth

import (
	"net/http"

	"golang.org/x/net/context"
	"gopkg.in/triki.v0/internal/ctx"
	"gopkg.in/triki.v0/internal/models/token"
	"gopkg.in/triki.v0/internal/models/user"
)

// Handler wraps handlers with context, session and authentication.
// The wrapped function fun can create extra DB session (if requred) and save it
// using ctx.SaveSession(...).
func Handler(fun func(context.Context, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, cancel := context.WithTimeout(context.Background(), RequestTimeout)
		defer cancel()
		cx := ctx.New(c, r)
		// close DB session (if the request handler created extra one)
		defer DBCloseSessions(r)

		// authenticate
		tokn := r.Header.Get("X-AUTHENTICATION-TOKEN")
		if tokn != "" {
			tkn, err := token.Find(cx, []byte(tokn))
			if err != nil {
				// TODO
				return
			}
			usr, err := user.GetByID(cx, tkn.UsrID)
			if err != nil {
				// TODO
				return
			}
			if !usr.IsActive() {
				// TODO
				return
			}
			//
			//context.Set(r, contextKey, usrID)
		}
		// TODO authenticate...
		fun(cx, w, r)
	}
}
