package wwwapi

import (
	"bitbucket.org/kornel661/triki/gotriki/db"
	"bitbucket.org/kornel661/triki/gotriki/log"
	"bytes"
	"encoding/json"
	"fmt"
	//"github.com/gorilla/mux"
	"net/http"
)

type (
	authInJSON struct {
		Session struct {
			Login    string `json:"login"`
			Password string `json:"password"`
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

// AuthPostHandler handles user authentication in /api/auth
func AuthPostHandler(w http.ResponseWriter, r *http.Request) {
	// log info
	var info bytes.Buffer
	defer func() { log.Infoln(info.String()) }()
	fmt.Fprintf(&info, "API access to %s by %s.", r.RequestURI, r.RemoteAddr)

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
	if usr, err := db.UserAuthenticate(authIn.Session.Login, authIn.Session.Password); err != nil {
		// not authenticated
		http.Error(w, "User authentication failed: "+err.Error(), http.StatusForbidden)
		fmt.Fprintf(&info, " Authentication failed for `%s`: %s.", usr.Login, err)
		return
	} else {
		// authentication successful
		writeJSON(w, authOutJSON{authOutSessionJSON{"token", usr.Login}})
		fmt.Fprintf(&info, " User `%s` authenticated.", usr.Login)
		return
	}
}
