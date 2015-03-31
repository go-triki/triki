package www

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
	"gopkg.in/triki.v0/internal/auth"
	"gopkg.in/triki.v0/internal/log"
	"gopkg.in/triki.v0/internal/models/user"
)

type (
	authInJSON struct {
		Session struct {
			Login    string `json:"login"`
			Password string `json:"pass"`
		} `json:"session"`
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
		log.Set(cx, log.FailedReadingRequestErr(err))
		// TODO
		return
	}

	// authenticate user
	usr, tkn, err := auth.Authenticate(cx, authIn.Session.Login, authIn.Session.Password)
	if err != nil {
		// not authenticated
		log.Set(cx, err)
		// TODO
		return
	}

	// authentication successful
	err = writeJSON(cx, w, authOutJSON{authOutSessionJSON{string(tkn.Tkn), usr.ID.Hex()}})
	if err != nil {
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
		log.Set(cx, log.FailedReadingRequestErr(err))
		// TODO
		return
	}
	// TODO verify email
	// sign-up user
	usr := user.T{
		Usr:  authIn.Session.Login,
		Pass: authIn.Session.Password,
	}
	err := user.New(cx, &usr)
	if err != nil {
		// not signed-up
		log.Set(cx, err)
		// TODO
		return
	}

	// sign-up successful
	log.Set(cx, log.SignupOK(usr.ID))
	// TODO do we want to put sth in the reply?
}
