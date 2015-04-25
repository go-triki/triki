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
	// defence against clickjacking
	w.Header().Set("X-Frame-Options", "DENY")
	defer func() {
		elaps := time.Since(start)
		err := context.Get(r, errKey)
		// log to DB
		dbErr := DBLog(map[string]interface{}{
			"time":   start,
			"elaps":  elaps,
			"method": r.Method,
			"url":    r.URL.Path,
			"error":  err,
		})
		// std log
		str := fmt.Sprintf("%s %s | start: %v elapsed: %v", r.Method, r.URL.Path, start, elaps)
		if dbErr != nil {
			str = fmt.Sprintf("%s | error logging to DB: %v", str, dbErr)
		}
		if err != nil {
			str = fmt.Sprintf("%s | error: %v", str, err)
		}
		StdLog.Println(str)
		context.Clear(r)
	}()
	mux.Router.ServeHTTP(w, r)
}
