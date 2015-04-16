// Package user provides model for user accounts.
package user // import "gopkg.in/triki.v0/internal/models/user"

import (
	"strings"

	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/triki.v0/internal/log"
)

var (
	// MinPassLen is the minimal password length that is accepted.
	MinPassLen int
)

var (
	// DBFind searches the DB for user with login/email == usr.
	DBFind func(cx context.Context, login string) (usr *T, err *log.Error)
	// DBFindByID finds user with a given _id.
	DBFindByID func(cx context.Context, id bson.ObjectId) (*T, *log.Error)
	// DBInsert inserts user usr into the DB.
	DBInsert func(cx context.Context, usr *T) *log.Error
	// DBExists checks if user with login/email == usr exists in the DB.
	DBExists func(cx context.Context, login string) (exists bool, err *log.Error)
)

// StatusT is user's status
type StatusT int

// Possible user statuses
const (
	SActive StatusT = iota
	SDeleted
	SSuspended
)

// T type (user.T) stores user information (e.g. for authentication), also for MongoDB and JSON.
type T struct {
	ID       bson.ObjectId `json:"id"             bson:"_id"`    // unique ID
	Usr      string        `json:"usr"            bson:"usr"`    // login/email (unique)
	Pass     string        `json:"pass,omitempty" bson:"-"`      // password (from www)
	PassHash []byte        `json:"-"              bson:"pass"`   // password hash (from DB)
	Nick     string        `json:"nick"           bson:"nick "`  // user's nick
	Status   StatusT       `json:"-"              bson:"status"` // user's status (active, deleted, etc.)
	// TODO add list of failed logins
}

// New adds new user to the database. Returns nil on success, error otherwise.
// usr.Usr and usr.Pass should be already initialized.
//
// User.Password is hashed into User.PassHash, ID is generated. Some other fields
// are sanitized (usr struct is updated accordingly).
func New(cx context.Context, usr *T) *log.Error {
	if usr.Nick == "" {
		at := strings.Index(usr.Usr, "@")
		if at == -1 {
			at = len(usr.Usr)
		}
		usr.Nick = usr.Usr[0:at]
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
	return DBInsert(cx, usr)
}

// GetByID retrieves user of a given ID from the DB
func GetByID(cx context.Context, id bson.ObjectId) (*T, *log.Error) {
	return DBFindByID(cx, id)
}

// IsActive returns true if user account is active and the user can log in.
func (usr *T) IsActive() bool {
	return usr.Status == SActive
}
