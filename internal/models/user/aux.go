package user

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/triki.v0/internal/log"
	"gopkg.in/triki.v0/internal/rands"
)

// checkPass carries out some sanity checks on the given user (pass/nick length, is
// login an email address). Used when creating a new user or changing password
// or email.
func (usr *T) checkPass() *log.Error {
	// password length
	if len(usr.Pass) < MinPassLen {
		return log.BadSignupDetailsErr("password too short")
	}
	// does login look like an email? len('a@b.c') == 5
	if len(usr.Usr) < 5 {
		return log.BadSignupDetailsErr("login must be a valid email address")
	}
	at := strings.Index(usr.Usr, "@")
	if at < 1 || at == len(usr.Usr)-1 {
		return log.BadSignupDetailsErr("login must be a valid email address")
	}
	if len(usr.Nick) < 1 {
		return log.BadSignupDetailsErr("nick cannot be empty")
	}
	// checks passed
	return nil
}

// setPass sets user's password to pass, generates new usr.PassHash and usr.Hash.
// The usr.Pass field is cleared.
func (usr *T) setPass(pass string) *log.Error {
	var err error
	usr.Pass = ""
	salt := rands.New(8)
	usr.Salt = salt
	usr.PassHash, err = bcrypt.GenerateFromPassword(append([]byte(PassSalt+pass), salt...), bcrypt.DefaultCost)
	if err != nil {
		log.InternalServerErr(err)
	}
	return nil
}
