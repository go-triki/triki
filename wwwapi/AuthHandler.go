package wwwapi

import (
	"bitbucket.org/kornel661/triki/gotriki/db"
	"bitbucket.org/kornel661/triki/gotriki/log"
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
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

const (
	authDelayMs = 150
)

func authDelay() {
	time.Sleep(time.Duration(authDelayMs+rand.Intn(authDelayMs)) * time.Millisecond)
}

// AuthPostHandler handles user authentication (logging in) in /api/auth.
func AuthPostHandler(w http.ResponseWriter, r *http.Request) {
	// simple guard, FIXME
	authDelay()

	// log info
	var info bytes.Buffer
	defer func() { log.Infoln(info.String()) }()
	apiAccessLog(&info, r)

	// decode request body
	dec := json.NewDecoder(r.Body)
	var authIn authInJSON
	err := dec.Decode(&authIn)
	if err != nil {
		// server doesn't understand the request
		fmt.Fprintf(&info, " Bad request syntax: %s.", err)
		Error(w, http.StatusBadRequest, statusError,
			"Bad request syntax: %s.", err)
		return
	}

	// authenticate user
	usr, token, err := db.UserAuthenticate(authIn.Session.Login, authIn.Session.Password)
	if err != nil {
		// not authenticated
		Error(w, http.StatusForbidden, statusError,
			"Authentication failed for user `%s`.", authIn.Session.Login)
		fmt.Fprintf(&info, " Authentication failed for `%s`: %s.", authIn.Session.Login, err)
		return
	}

	// authentication successful
	writeJSON(w, &info, authOutJSON{authOutSessionJSON{token, usr.ID.Hex()}})
	fmt.Fprintf(&info, " User `%s` authenticated.", usr.Login)
}

// AuthSignupPostHandler handles user sign-up process in /api/auth/signup.
func AuthSignupPostHandler(w http.ResponseWriter, r *http.Request) {
	// simple guard, FIXME
	authDelay()

	// log info
	var info bytes.Buffer
	defer func() { log.Infoln(info.String()) }()
	apiAccessLog(&info, r)

	// decode request body
	dec := json.NewDecoder(r.Body)
	var authIn authInJSON
	err := dec.Decode(&authIn)
	if err != nil {
		// server doesn't understand the request
		fmt.Fprintf(&info, " Bad request syntax: %s.", err)
		Error(w, http.StatusBadRequest, statusError,
			"Bad request syntax: %s.", err)
		return
	}

	// sign-up user
	err = db.UserSignup(authIn.Session.Login, authIn.Session.Password)
	if err != nil {
		// not authenticated
		Error(w, http.StatusForbidden, statusError,
			"Sign-up failed for user `%s`. Reason: %s.", authIn.Session.Login, err)
		fmt.Fprintf(&info, " Sign-up failed for `%s`: %s.", authIn.Session.Login, err)
		return
	}

	// sign-up successful
	Error(w, http.StatusAccepted, statusError,
		"Sign-up successful for user `%s`. Email verification required.", authIn.Session.Login)
	fmt.Fprintf(&info, " User `%s` signed-up.", authIn.Session.Login)
}
