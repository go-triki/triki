package auth

import (
	"net/http"

	"golang.org/x/net/context"
	"gopkg.in/triki.v0/internal/ctx"
)

// Handler wraps handlers with context, session and authentication.
// The wrapped function fun can create extra DB session (if requred) and save it
// using ctx.SaveSession(...).
func Handler(fun func(context.Context, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, cancel := context.WithTimeout(context.Background(), RequestTimeout)
		cx := ctx.New(c, r)
		defer cancel()
		// close DB session (if the request handler created extra one)
		defer func() {
			s, _ := ctx.DBSessionFromReq(r)
			if s != nil {
				s.Close()
			}
		}()
		tkn := r.Header.Get("X-AUTHENTICATION-TOKEN")
		if tkn != "" {
			if !bson.IsObjectIdHex(tkn) {
				Error(w, http.StatusForbidden, statusInvalidToken,
					"Invalid authentication token.")
				return
			}
			tknID := bson.ObjectIdHex(tkn)
			usrID, err := db.TokenCheck(tknID)
			if err != nil {
				Error(w, http.StatusForbidden, statusInvalidToken,
					"Invalid authentication token.")
				return
			}
			context.Set(r, contextKey, usrID)
		}
		// TODO authenticate...
		fun(cx, w, r)
	}
}
