/*
Package wwwapi supplies handlers for gotriki www api.
*/
package wwwapi

import (
	"encoding/json"
	"net/http"
)

const (
	contentType         = "Content-Type"
	applicationJSONType = "application/json"
)

// writeJSON writes JSON representation of v to the http response w.
func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set(contentType, applicationJSONType)
	enc := json.NewEncoder(w)
	enc.Encode(v)
}
