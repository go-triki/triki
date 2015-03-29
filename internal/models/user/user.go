// Package user provides model for user accounts.
package user // import "gopkg.in/triki.v0/internal/models/user"

import (
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
	// DBFind searches the DB for user with login/email == usr.
	DBFind func(login string) (usr *T, err *log.Error)
	// DBInsert inserts user usr into the DB.
	DBInsert func(usr *T) *log.Error
	// DBExists checks if user with login/email == usr exists in the DB.
	DBExists func(login string) (exists bool, err *log.Error)
)

// T type (user.T) stores user information (e.g. for authentication), also for MongoDB and JSON.
type T struct {
	ID       bson.ObjectId `json:"id"             bson:"_id"`   // unique ID
	Usr      string        `json:"usr"            bson:"usr"`   // login/email
	Pass     string        `json:"pass,omitempty" bson:"-"`     // password (from www)
	PassHash []byte        `json:"-"              bson:"pass"`  // password hash (from DB)
	Salt     []byte        `json:"-"              bson:"salt"`  // individual password salt (from DB)
	Nick     string        `json:"nick"           bson:"nick "` // user's nick
	// TODO add list of failed logins
}

// New adds new user to the database. Returns nil on success, error otherwise.
// usr.Usr and usr.Pass should be already initialized.
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
	return DBInsert(usr)
}
