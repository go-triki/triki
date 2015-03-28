/*
Package wwwapi supplies handlers for gotriki www api.
*/
package wwwapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/triki.v0/internal/db"
)

const (
	contentType                        = "Content-Type"
	applicationJSONType                = "application/json"
	contextKey          contextVarsKey = 47656
)

type (
	// type for context response variables
	contextVarsKey int
	// triki http error codes
	errorCode int
)

// triki http error codes
const (
	statusError           = 0
	statusUnauthorized    = 100
	statusInvalidToken    = 110
	statusResourceInvalid = 200
)

func init() {
	// initialize random number generator
	rand.Seed(time.Now().UTC().UnixNano())
}

// writeJSON writes JSON representation of v to the http response w.
// Errors are written to the log (if not nil).
func writeJSON(w http.ResponseWriter, log *bytes.Buffer, v interface{}) error {
	w.Header().Set(contentType, applicationJSONType)
	enc := json.NewEncoder(w)
	err := enc.Encode(v)
	if err != nil && log != nil {
		fmt.Fprintf(log, " Error writing reply: %s.", err)
	}
	return err
}

// apiAccessLog appends to the buffer API access information.
func apiAccessLog(buf *bytes.Buffer, r *http.Request) {
	fmt.Fprintf(buf, "API access to %s by %s.", r.RequestURI, r.RemoteAddr)
}

// Error writes an error response to w, similarily to http.Error.
// If trikiCode is 0 it behaves as http.Error with string formatting.
// If trikiCode is != 0 an arror code is encoded at the beginning of the response.
func Error(w http.ResponseWriter, code int, trikiCode errorCode, format string, a ...interface{}) {
	errorDesc := fmt.Sprintf(format, a...)
	if trikiCode != 0 {
		errorDesc = fmt.Sprintf("%d.%d: %s", code, trikiCode, errorDesc)
	}
	http.Error(w, errorDesc, code)
}

// AuthenticateHandler wraps regular http handler and handles authentication.
func AuthenticateHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tkn := r.Header.Get("X-AUTHENTICATION-TOKEN")
		if tkn != "" {
			if !bson.IsObjectIdHex(tkn) {
				Error(w, http.StatusForbidden, statusInvalidToken,
					"Invalid authentication token.")
				return
			}
			tknID := bson.ObjectIdHex(tkn)
			usrID, err := db.TokenCheck(tknID)
			if err != nil {
				Error(w, http.StatusForbidden, statusInvalidToken,
					"Invalid authentication token.")
				return
			}
			context.Set(r, contextKey, usrID)
		}
		handler(w, r)
		// context clear not needed when using mux
		//context.Clear(r)
	}
}

// authenticatedUser checks if the request has been authenticated for a user (whose ID is returned in the second value).
func authenticatedUser(r *http.Request) (bool, bson.ObjectId) {
	v := context.Get(r, contextKey)
	if v == nil {
		return false, ""
	}
	usrID, ok := v.(bson.ObjectId)
	return ok, usrID
}

// hasPermission checks if the http request has clerance not below the level of the user with ID usrLevel.
func hasPermission(usrLevel bson.ObjectId, r *http.Request) bool {
	auth, usr := authenticatedUser(r)
	// for now:
	return auth && (usr == usrLevel)
}
