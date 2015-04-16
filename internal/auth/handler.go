package auth

import (
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/triki.v0/internal/ctx"
	"gopkg.in/triki.v0/internal/log"
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
		// close DB sessions (in case the request handler created some)
		defer DBCloseSessions(r)

		// authenticate
		tokn := r.Header.Get("X-AUTHENTICATION-TOKEN")
		// tokn should be of the form "tknID:token"
		if tokn != "" {
			tkns := strings.Split(tokn, ":")
			if len(tkns) != 2 {
				WriteErrorHandler(cx, w, r, log.BadTokenFormatErr)
				return
			}
			if !bson.IsObjectIdHex(tkns[0]) {
				WriteErrorHandler(cx, w, r, log.BadTokenFormatErr)
				return
			}
			tknID := bson.ObjectId(tkns[0])
			tokn := []byte(tkns[1])
			tkn, err := token.Find(cx, tknID)
			if err != nil {
				WriteErrorHandler(cx, w, r, err)
				return
			}
			// check if user-supplied token is correct
			if err := bcrypt.CompareHashAndPassword(tkn.Hash, tokn); err != nil {
				WriteErrorHandler(cx, w, r, log.BadTokenErr)
				return
			} else {
				// TODO ? record successful access (maybe after usr.IsActive)
			}
			usr, err := user.GetByID(cx, tkn.UsrID)
			if err != nil {
				WriteErrorHandler(cx, w, r, err)
				return
			}
			if !usr.IsActive() {
				WriteErrorHandler(cx, w, r, log.UserNotActiveErr)
				return
			}
			ath := T{
				Usr: usr,
				Tkn: tkn,
			}
			Set(cx, &ath)
			// TODO record in log that user logged in
		}
		fun(cx, w, r)
	}
}
