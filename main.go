/*
Gotriki server of the Trikipedia - the truth encyclopedia.
*/
package main // import "gopkg.in/triki.v0"

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"gopkg.in/kornel661/nserv.v0"
	"gopkg.in/mgo.v2"
	"gopkg.in/triki.v0/internal/db"
	tlog "gopkg.in/triki.v0/internal/log"
)

const (
	apiPrefix = "/api/"
)

var (
	server    = nserv.Server{}
	mDialInfo = mgo.DialInfo{}
)

func main() {

	// catch signals & shutdown the server
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	// redirect "stop" signals to server.Stop
	go func() {
		<-signals
		log.Infoln("Caught signal. Exiting (waiting for open connections to termiate).")
		server.Stop()
	}()

	// setup database connections, etc.
	db.Setup()
	defer db.Cleanup()

	// setup logger
	// TODO(km): log.DbLog = ...
	defer tlog.Flush()

	r := mux.NewRouter()

	// setup API routing
	apiRouter := r.PathPrefix(apiPrefix).Subrouter()
	routeAPI(apiRouter)

	// serve static content from "/"
	r.PathPrefix(staticPrefix).Handler(http.FileServer(http.Dir(optServRoot)))

	// start server
	log.Infof("Serving triki via www: http://%s\n", *optServerRoot)
	server.Handler = &tlog.LoggedMux{r}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	log.Infoln("Exiting gracefully, please wait...")
	server.Wait()
}
