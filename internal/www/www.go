/*
Package www supplies http handlers for gotriki www api.
*/
package www // import "gopkg.in/triki.v0/internal/www"

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/triki.v0/internal/log"
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
func writeJSON(cx context.Context, w http.ResponseWriter, v interface{}) *log.Error {
	w.Header().Set(contentType, applicationJSONType)
	enc := json.NewEncoder(w)
	err := enc.Encode(v)
	if err != nil {
		return log.FailedWritingReplyErr(err)
	}
	return nil
}

func readID(str string) (bson.ObjectId, *log.Error) {
	if !bson.IsObjectIdHex(str) {
		return "", log.InvalidIDErr
	}
	return bson.ObjectIdHex(str), nil
}
