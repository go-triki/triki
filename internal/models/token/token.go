/*
Package token provides model for authorization tokens.
*/
package token // import "gopkg.in/triki.v0/internal/models/token"
import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/triki.v0/internal/log"
	"gopkg.in/triki.v0/internal/rands"
)

// Options set by conf.
var (
	// MaxExpireAfter controls how long (at most) authentication tokens are valid.
	MaxExpireAfter time.Duration
)

// Set by DB driver.
var (
	// DBFind finds given token in the DB.
	DBFind func(cx context.Context, tknID bson.ObjectId) (*T, *log.Error)
	// DBExists checks if given token is in the DB.
	DBExists func(cx context.Context, tknID bson.ObjectId) (bool, *log.Error)
	// DBInsert inserts token tkn into the DB.
	DBInsert func(cx context.Context, tkn *T) *log.Error
)

// T type (token.T) holds information associated with a given authentication token
type T struct {
	ID          bson.ObjectId `bson:"_id"          json:"id"`
	Tkn         []byte        `bson:"-"            json:"tkn"`
	Hash        []byte        `bson:"hash"         json:"-"`
	Birth       time.Time     `bson:"birth"        json:"-"`
	UsrID       bson.ObjectId `bson:"usrID"        json:"usrID"`
	ExpireAfter time.Duration `bson:"expire_after" json:"-"`
	// TODO add last used array with info
}

// New creates new token for user usrID (and saves it in the DB).
//
// At least tkn.UsrID needs to be set.
func New(cx context.Context, tkn *T) *log.Error {
	tkn.ID = bson.NewObjectId()
	tkn.Tkn = rands.New(20)
	var berr error
	tkn.Hash, berr = bcrypt.GenerateFromPassword([]byte(tkn.Tkn), bcrypt.DefaultCost)
	if berr != nil {
		return log.InternalServerErr(berr)
	}
	tkn.Birth = time.Now()
	if tkn.ExpireAfter == 0 {
		tkn.ExpireAfter = MaxExpireAfter
	}
	if tkn.ExpireAfter > MaxExpireAfter {
		tkn.ExpireAfter = MaxExpireAfter
	}
	// TODO ? add info (from cx) about request
	err := DBInsert(cx, tkn)
	if err != nil {
		return err
	}
	return nil
}

// Find finds given token in the DB.
func Find(cx context.Context, tknID bson.ObjectId) (*T, *log.Error) {
	tkn, err := DBFind(cx, tknID)
	if err != nil {
		is, er := DBExists(cx, tknID)
		if is { // token in the DB but there was an error retrieving it
			return nil, err
		} else if er != nil { // there's no such token in the DB
			return nil, log.BadTokenErr
		} else { // error checking if token exists
			return nil, err // return original error
		}
		return nil, err
	}
	if time.Now().After(tkn.Birth.Add(tkn.ExpireAfter)) {
		// token expired
		return nil, log.BadTokenErr
	}
	return tkn, nil
}
