/*
Gotriki server of the Trikipedia - the truth encyclopedia.
*/
package main

import (
	"bitbucket.org/kornel661/triki/gotriki/conf"
	"bitbucket.org/kornel661/triki/gotriki/db"
	"bitbucket.org/kornel661/triki/gotriki/log"
	"github.com/gorilla/mux"
	"github.com/kornel661/manners"
	"net/http"
	"os"
	"os/signal"
)

const (
	staticPrefix = "/static"
	apiPrefix    = "/api/"
)

func main() {
	// panic trap (for fatal errors from log)
	defer func() {
		if r := recover(); r != nil {
			if fatal, ok := r.(log.FatalErrorPanic); ok {
				log.Infof("Panic: %s.\n", fatal)
			} else {
				panic(r)
			}
		}
	}()
	conf.Setup()

	// catch signals & shutdown the server
	server := manners.NewServer()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	// "reroute" signals channel to server.Shutdown
	go func() {
		<-signals
		log.Infoln("Caught signal. Exiting (waiting for open connections).")
		server.Shutdown <- true
	}()

	// setup database connections, etc.
	db.Setup()
	defer db.Cleanup()

	// setup routing
	r := mux.NewRouter()
	apiRouter := r.PathPrefix(apiPrefix).Subrouter()
	routeAPI(apiRouter)

	// serve static content from staticPrefix and "/"
	r.PathPrefix(staticPrefix).Handler(http.FileServer(http.Dir(conf.Server.Root)))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(conf.Server.Root + staticPrefix)))

	// start server
	log.Infof("Serving gotriki via www: http://%s\n", conf.Server.Addr)
	if err := server.ListenAndServe(conf.Server.Addr, r); err != nil {
		log.Fatal(err)
	}
	log.Infoln("Exiting gracefully...")
}
