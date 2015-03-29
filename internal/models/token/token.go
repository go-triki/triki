/*
Package token provides model for authorization tokens.
*/
package token // import "gopkg.in/triki.v0/internal/models/token"
import (
	"time"

	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/triki.v0/internal/log"
	"gopkg.in/triki.v0/internal/rands"
)

var (
	// DBFind finds given token in the DB.
	DBFind func(cx context.Context, tknID []byte) (*T, *log.Error)
	// DBInsert inserts token tkn into the DB.
	DBInsert func(cx context.Context, tkn *T) *log.Error
)

// T type (token.T) holds information associated with a given authentication token
type T struct {
	Tkn   []byte        `bson:"_id"   json:"tkn"`
	Birth time.Time     `bson:"birth" json:"-"`
	UsrID bson.ObjectId `bson:"usrID" json:"usrID"`
	// TODO add last used array with info
}

// New creates new token for user usrID (and saves it in the DB).
//
// Returns token and error message.
func New(cx context.Context, usrID bson.ObjectId) (tkn *T, err *log.Error) {
	var token T
	token.Tkn = rands.New(30)
	token.Birth = time.Now()
	token.UsrID = usrID
	// TODO add info about request
	err = DBInsert(cx, &token)
	if err != nil {
		return nil, err
	}
	return &token, err
}
