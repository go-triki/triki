package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/triki.v0/internal/wwwapi"
)

const (
	authPath = "/auth"
)

// authH wraps handler with wwwapi.AuthenticateHandler.
// Just for convenience (short name).
func authH(handler http.HandlerFunc) http.HandlerFunc {
	return wwwapi.AuthenticateHandler(handler)
}

// routeAPI associates http handlers with corresponding URLs.
func routeAPI(r *mux.Router) {
	r.HandleFunc(authPath, wwwapi.AuthPostHandler).Methods("POST")
	r.HandleFunc(authPath+"/signup", wwwapi.AuthSignupPostHandler).Methods("POST")

	r.HandleFunc("/accounts/{account_id}", authH(wwwapi.AccountsIDGetHandler)).Methods("GET")
}
