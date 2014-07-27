package main

import (
	"bitbucket.org/kornel661/triki/gotriki/wwwapi"
	"github.com/gorilla/mux"
)

const (
	authPath = "/auth"
)

func routeAPI(r *mux.Router) {
	r.HandleFunc(authPath, wwwapi.AuthPostHandler).Methods("POST")
	r.HandleFunc(authPath+"/signup", wwwapi.AuthSignupPostHandler).Methods("POST")
}
