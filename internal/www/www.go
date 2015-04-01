/*
Package www supplies http handlers for gotriki www api.
*/
package www // import "gopkg.in/triki.v0/internal/www"

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/triki.v0/internal/log"
)

const (
	contentType         = "Content-Type"
	applicationJSONType = "application/json"
)

// writeError records error in the context (for logging) and sends it to the client
// via http.ResponseWriter.
func writeError(cx context.Context, w http.ResponseWriter, err *log.Error) {
	log.Set(cx, log.FailedReadingRequestErr(err))
	writeJSON(cx, w, Resp{
		Errors: []*log.Error{err},
	})
}

// writeJSON writes JSON representation of v to the http response w.
func writeJSON(cx context.Context, w http.ResponseWriter, v interface{}) *log.Error {
	w.Header().Set(contentType, applicationJSONType)
	enc := json.NewEncoder(w)
	err := enc.Encode(v)
	if err != nil {
		return log.FailedWritingReplyErr(err)
	}
	return nil
}

// readID converts a string representing an ID into bson.ObjectId
func readID(str string) (bson.ObjectId, *log.Error) {
	if !bson.IsObjectIdHex(str) {
		return "", log.InvalidIDErr
	}
	return bson.ObjectIdHex(str), nil
}
