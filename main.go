/*
Gotriki server of the Trikipedia - the truth encyclopedia.
*/
package main // import "gopkg.in/triki.v0"

import (
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"gopkg.in/kornel661/nserv.v0"
	"gopkg.in/triki.v0/internal/db/mongodrv"
	tlog "gopkg.in/triki.v0/internal/log"
)

var (
	server = nserv.Server{}
)

func main() {
	// flush logger
	defer tlog.Flush()

	// catch signals & shutdown the server
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	// redirect "stop" signals to server.Stop
	go func() {
		<-signals
		log.Println("Caught signal. Exiting (waiting for open connections to termiate).")
		server.Stop()
	}()

	// setup database connections, etc.
	mongo.Setup()
	defer mongo.Cleanup()

	r := mux.NewRouter()
	// tlog.LoggedMux clears contexts for us
	r.KeepContext = true

	// setup API routing
	apiRouter := r.PathPrefix("/api").Subrouter()
	routeAPI(apiRouter)

	if staticServerURL == nil {
		// serve static content from "/"
		r.PathPrefix("/").Handler(http.FileServer(http.Dir(optServRoot)))
	} else {
		// serve static content from staticServerURL
		r.PathPrefix("/").Handler(httputil.NewSingleHostReverseProxy(staticServerURL))
	}
	// start server
	log.Printf("Serving triki via www: http://%s\n", server.Addr)
	server.Handler = &tlog.LoggedMux{Router: r}
	server.ErrorLog = tlog.StdLog
	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
		return
	}
	log.Println("Exiting gracefully, please wait for active connections...")
	server.Wait()
	log.Println("All connections terminated.")
}
