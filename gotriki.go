/*
gotriki server of the trikipedia - the truth encyclopedia
*/
package main

import (
	"bitbucket.org/kornel661/triki/gotriki/conf"
	"bitbucket.org/kornel661/triki/gotriki/db"
	"bitbucket.org/kornel661/triki/gotriki/log"
	"github.com/gorilla/mux"
	"github.com/kornel661/manners"
	"io"
	"net/http"
	"os"
	"os/signal"
)

func homeView(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	headers.Add("Content-Type", "text/html")
	log.Infof("Dealing with a request...\n")
	rest := mux.Vars(r)["rest"]
	io.WriteString(w, "<html><head></head><body><p>It works!<br>rest: "+rest+"</p></body></html>")
}

const (
	staticPrefix = "/static"
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
		log.Infoln("Cought signal. Exiting.")
		server.Shutdown <- true
	}()

	// setup database connections, etc.
	db.Setup()
	defer db.Cleanup()

	r := mux.NewRouter()
	r.HandleFunc("/fcgi-test", homeView)
	r.HandleFunc("/fcgi-test/{rest:.*}", homeView)

	// serve static content from staticPrefix and "/"
	r.PathPrefix(staticPrefix).Handler(http.FileServer(http.Dir(conf.Server.Root)))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(conf.Server.Root + staticPrefix)))

	log.Infof("Serving via www: http://%s\n", conf.Server.Addr)
	if err := server.ListenAndServe(conf.Server.Addr, r); err != nil {
		log.Fatal(err)
	}
	log.Infoln("Exiting gracefully...")
}
