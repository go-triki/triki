/*
Gotriki server of the Trikipedia - the truth encyclopedia.
*/
package main // import "gopkg.in/triki.v0"

import (
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"gopkg.in/kornel661/nserv.v0"
	"gopkg.in/triki.v0/internal/db/mongodrv"
	"gopkg.in/triki.v0/internal/log"
)

var (
	server = nserv.Server{}
)

func main() {
	// flush logger
	defer log.Flush()

	// catch signals & shutdown the server
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	// redirect "stop" signals to server.Stop
	go func() {
		<-signals
		log.StdLog.Println("Caught signal. Exiting (waiting for open connections to termiate).")
		server.Stop()
	}()

	// setup database connections, etc.
	mongo.Setup()
	defer mongo.Cleanup()

	r := mux.NewRouter()
	// log.LoggedMux clears contexts for us
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
	log.StdLog.Printf("Serving triki via www: http://%s\n", server.Addr)
	log.Flush()
	server.Handler = &log.LoggedMux{r}
	server.ErrorLog = log.StdLog
	if err := server.ListenAndServe(); err != nil {
		log.StdLog.Println(err)
		return
	}
	log.StdLog.Println("Exiting gracefully, please wait for active connections...")
	log.Flush()
	server.Wait()
	log.StdLog.Println("All connections terminated.")
	log.Flush()
}
