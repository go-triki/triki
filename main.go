/*
Gotriki server of the Trikipedia - the truth encyclopedia.
*/
package main // import "gopkg.in/triki.v0"

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"gopkg.in/kornel661/nserv.v0"
	"gopkg.in/triki.v0/internal/db/mongodrv"
	"gopkg.in/triki.v0/internal/log"
)

const (
	apiPrefix = "/api/"
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
		log.StdLog.Infoln("Caught signal. Exiting (waiting for open connections to termiate).")
		server.Stop()
	}()

	// setup database connections, etc.
	mongo.Setup()
	defer mongo.Cleanup()

	r := mux.NewRouter()
	// log.LoggedMux clears contexts for us
	r.KeepContext = true

	// setup API routing
	apiRouter := r.PathPrefix(apiPrefix).Subrouter()
	routeAPI(apiRouter)

	// serve static content from "/"
	r.PathPrefix(staticPrefix).Handler(http.FileServer(http.Dir(optServRoot)))

	// start server
	log.StdLog.Infof("Serving triki via www: http://%s\n", server.Addr)
	server.Handler = &log.LoggedMux{r}
	if err := server.ListenAndServe(); err != nil {
		log.StdLog.Println(err)
		return
	}
	log.StdLog.Infoln("Exiting gracefully, please wait...")
	server.Wait()
}
