/*
gotriki server
*/
package main

import (
	"bitbucket.org/kornel661/triki/gotriki/conf"
	"bitbucket.org/kornel661/triki/gotriki/db"
	"bitbucket.org/kornel661/triki/gotriki/log"
	"github.com/gorilla/mux"
	//"github.com/kornel661/manners"
	"io"
	"net/http"
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
	conf.Setup()
	db.Setup()

	r := mux.NewRouter()
	r.HandleFunc("/fcgi-test", homeView)
	r.HandleFunc("/fcgi-test/{rest:.*}", homeView)

	// serve static content
	r.PathPrefix(staticPrefix).Handler(http.FileServer(http.Dir(conf.Server.Root)))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(conf.Server.Root + staticPrefix)))

	log.Infof("Serving via www: http://%s\n", conf.Server.Addr)
	if err := http.ListenAndServe(conf.Server.Addr, r); err != nil {
		log.Fatal(err)
	}
}
