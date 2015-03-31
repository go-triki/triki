package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"gopkg.in/triki.v0/internal/auth"
	"gopkg.in/triki.v0/internal/www"
)

// routeAPI associates http handlers with corresponding URLs.
func routeAPI(r *mux.Router) {
	// TODO r.NotFoundHandler
	add := func(path string, h func(context.Context, http.ResponseWriter, *http.Request)) *mux.Route {
		return r.HandleFunc(path, auth.Handler(h))
	}

	// auth
	add("/auth/login", www.AuthLoginPost).Methods("POST")
	add("/auth/signup", www.AuthSignupPost).Methods("POST")

	// users
	add("/users/{user_id}", www.UsersIDGet).Methods("GET")
}
