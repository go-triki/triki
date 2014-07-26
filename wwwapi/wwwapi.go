/*
Package wwwapi supplies handlers for gotriki www api.
*/
package wwwapi

import (
	"encoding/json"
	"net/http"
)

// writeJSON writes JSON representation of v to the http response w.
func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(v)
}
