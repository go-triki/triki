// Package user provides structs that model user accounts.
package user // import "gopkg.in/triki.v0/internal/models/user"

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/triki.v0/internal/log"
)

var (
	// PassSalt is used to salt all user passwords.
	PassSalt string
	// MinPassLen is the minimal password length that is accepted.
	MinPassLen int
)

var (
	// Find searches the DB for user with login/email == usr.
	Find func(user string) (*T, *log.Error)
	// Insert inserts user usr into the DB.
	Insert func(usr *T) *log.Error

	// TokenNew returns new authentication token for user with ID == usrId.
	TokenNew func(usrID bson.ObjectId) (tkn string, err *log.Error)
)

// T type (user.T) stores user information (e.g. for authentication), also for MongoDB and JSON.
type T struct {
	ID       bson.ObjectId `json:"id"             bson:"_id"`   // unique ID
	Usr      string        `json:"usr"            bson:"usr"`   // login/email
	Pass     string        `json:"pass,omitempty" bson:"-"`     // password (from www)
	PassHash []byte        `json:"-"              bson:"pass"`  // password hash (from DB)
	Salt     string        `json:"-"              bson:"salt"`  // individual password salt (from DB)
	Nick     string        `json:"nick"           bson:"nick "` // user's nick
	// TODO add list of failed logins
}

// New adds new user to the database. Returns nil on success, error otherwise.
//
// User.Password is hashed into User.PassHash, ID is generated. Some other fields
// are sanitized (usr struct is updated accordingly).
func New(usr *T) *log.Error {
	if usr.Nick == "" {
		usr.Nick = usr.Usr
	}
	err := usr.checkPass()
	if err != nil {
		return err
	}
	err = usr.setPass(usr.Pass)
	if err != nil {
		return err
	}
	usr.ID = bson.NewObjectId()
	return Insert(usr)
}

// Authenticate checks if the given credentials can be authenticated.
// Returns (user, token, error).
func Authenticate(usrID, pass string, req *http.Request) (usr *T, tkn string, err *log.Error) {
	usr, err = Find(usrID)
	if err != nil {
		return nil, "", err
	}
	// TODO check recent failed logins
	if err := bcrypt.CompareHashAndPassword(usr.PassHash, []byte(PassSalt+pass+usr.Salt)); err != nil {
		// TODO no more than 3 failed logins in 10 minutes? (use req to get info)
		return nil, "", log.IncorrectPassErr
	}
	token, err := TokenNew(usr.ID)
	if err != nil {
		return nil, "", log.InternalServerErr(err)
	}
	return usr, token, nil
}
