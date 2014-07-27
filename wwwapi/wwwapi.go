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
)

const (
	contentType         = "Content-Type"
	applicationJSONType = "application/json"
)

func init() {
	// initialize random number generator
	rand.Seed(time.Now().UTC().UnixNano())
}

// writeJSON writes JSON representation of v to the http response w.
func writeJSON(w http.ResponseWriter, v interface{}) error {
	w.Header().Set(contentType, applicationJSONType)
	enc := json.NewEncoder(w)
	return enc.Encode(v)
}

// apiAccessLog appends to the buffer API access information.
func apiAccessLog(buf *bytes.Buffer, r *http.Request) {
	fmt.Fprintf(buf, "API access to %s by %s.", r.RequestURI, r.RemoteAddr)
}
