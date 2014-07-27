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

// AuthPostHandler handles user authentication in /api/auth.
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
		http.Error(w, "Bad request syntax: "+err.Error(), http.StatusBadRequest)
		return
	}

	// authenticate user
	usr, err := db.UserAuthenticate(authIn.Session.Login, authIn.Session.Password)
	if err != nil {
		// not authenticated
		http.Error(w, "Authentication failed for user `"+authIn.Session.Login+"`.", http.StatusForbidden)
		fmt.Fprintf(&info, " Authentication failed for `%s`: %s.", authIn.Session.Login, err)
		return
	}

	// authentication successful
	err = writeJSON(w, authOutJSON{authOutSessionJSON{"token", usr.Login}})
	if err != nil {
		fmt.Fprintf(&info, "Error writing reply: %s.", err)
		log.Warningln(info.String())
	}
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
		http.Error(w, "Bad request syntax: "+err.Error(), http.StatusBadRequest)
		return
	}

	// sign-up user
	err = db.UserSignup(authIn.Session.Login, authIn.Session.Password)
	if err != nil {
		// not authenticated
		http.Error(w, "Sign-up failed for user `"+authIn.Session.Login+
			"`. Reason: "+err.Error()+".",
			http.StatusForbidden)
		fmt.Fprintf(&info, " Sign-up failed for `%s`: %s.", authIn.Session.Login, err)
		return
	}

	// sign-up successful
	http.Error(w, "Sign-up successful for user `"+authIn.Session.Login+"`. Email verification required.", http.StatusAccepted)
	fmt.Fprintf(&info, " User `%s` signed-up.", authIn.Session.Login)
}
