package log

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// LoggedMux will log requests in StdLog and in a DB via DbLog. Uses contexts to
// retreive errors.
type LoggedMux struct {
	*mux.Router
}

func (mux *LoggedMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	mux.Router.ServeHTTP(w, r)
	elaps := time.Since(start)
	errors := context.Get(r, errKey)
	// log to DB
	err := DBLog(map[string]interface{}{
		"time":   start,
		"elaps":  elaps,
		"method": r.Method,
		"url":    r.URL.Path,
		"errors": errors,
	})
	// std log
	str := fmt.Sprintf("%s %s | start: %v elapsed: %v")
	if err != nil {
		str = fmt.Sprintf("%s | error logging to DB: %v", str, err)
	}
	if errors != nil {
		str = fmt.Sprintf("%s | errors: %v", str, errors)
	}
	StdLog.Println(str)
	context.Clear(r)
}
