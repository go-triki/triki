package www

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
	"gopkg.in/triki.v0/internal/auth"
	"gopkg.in/triki.v0/internal/log"
	"gopkg.in/triki.v0/internal/models/token"
	"gopkg.in/triki.v0/internal/models/user"
)

type (
	authInJSON struct {
		Login    string `json:"login"`
		Password string `json:"pass"`
	}
	authOutSessionJSON struct {
		AuthToken     string `json:"auth_token"`
		AuthAccountID string `json:"account_id"`
	}
	authOutJSON struct {
		Session authOutSessionJSON `json:"session"`
	}
)

// AuthLoginPost handles user authentication (logging in) in /api/auth/login.
func AuthLoginPost(cx context.Context, w http.ResponseWriter, r *http.Request) {
	// decode request body
	dec := json.NewDecoder(r.Body)
	var authIn authInJSON
	if err := dec.Decode(&authIn); err != nil {
		// cannot read/decode the request
		writeError(cx, w, log.FailedReadingRequestErr(err))
		return
	}

	// authenticate user
	usr, tkn, err := auth.Authenticate(cx, authIn.Login, authIn.Password)
	if err != nil {
		// not authenticated
		writeError(cx, w, err)
		return
	}

	// authentication successful
	err = writeJSON(cx, w, Resp{
		Users:  []*user.T{usr},
		Tokens: []*token.T{tkn},
	})
	if err != nil {
		// error writing response, just log
		log.Set(cx, err)
		return
	}
	log.Set(cx, log.LoginOK(usr.ID, tkn.Tkn))
}

// AuthSignupPost handles user sign-up process in /api/auth/signup.
func AuthSignupPost(cx context.Context, w http.ResponseWriter, r *http.Request) {
	// decode request body
	dec := json.NewDecoder(r.Body)
	var authIn authInJSON
	if err := dec.Decode(&authIn); err != nil {
		// cannot read/decode the request
		writeError(cx, w, log.FailedReadingRequestErr(err))
		return
	}
	// TODO verify email
	// sign-up user
	usr := user.T{
		Usr:  authIn.Login,
		Pass: authIn.Password,
	}
	err := user.New(cx, &usr)
	if err != nil {
		// not signed-up
		writeError(cx, w, err)
		return
	}

	// sign-up successful
	// TODO do we want to put sth in the reply?
	err = writeJSON(cx, w, Resp{})
	if err != nil {
		// error writing response, just log
		log.Set(cx, err)
		return
	}
	log.Set(cx, log.SignupOK(usr.ID))
}
