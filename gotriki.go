/*
	gotriki server
*/
package main

import (
	"bitbucket.org/kornel661/triki/gotriki/conf"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net"
	"net/http"
)

func homeView(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	headers.Add("Content-Type", "text/html")
	log.Printf("Dealing with a request...\n")
	rest := mux.Vars(r)["rest"]
	io.WriteString(w, "<html><head></head><body><p>It works!<br>rest: "+rest+"</p></body></html>")
}

func main() {
	conf.Setup()

	r := mux.NewRouter()
	r.HandleFunc("/fcgi-test", homeView)
	r.HandleFunc("/fcgi-test/{rest:.*}", homeView)

	var err error

	log.Printf("Serving via www: %s.\n", conf.Server.Addr)
	err = http.ListenAndServe(conf.Server.Addr, r)

	if err != nil {
		log.Fatal(err)
	}
}
