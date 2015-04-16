package auth

import (
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"gopkg.in/triki.v0/internal/log"
	"gopkg.in/triki.v0/internal/models/token"
	"gopkg.in/triki.v0/internal/models/user"
)

// Authenticate checks if the given credentials can be authenticated.
// Returns (user, token, error).
func Authenticate(cx context.Context, login, pass string) (*user.T, *token.T, *log.Error) {
	usr, err := user.DBFind(cx, login)
	if err != nil {
		return nil, nil, err
	}
	if !usr.IsActive() {
		return nil, nil, log.UserNotActiveErr
	}
	// TODO check recent failed logins, move it to user package?
	if err := bcrypt.CompareHashAndPassword(usr.PassHash, []byte(pass)); err != nil {
		// TODO no more than 3 failed logins in 10 minutes? (use req to get info)
		// add failed login
		return nil, nil, log.IncorrectPassErr
	}
	tkn := &token.T{UsrID: usr.ID}
	err = token.New(cx, tkn)
	if err != nil {
		return nil, nil, log.InternalServerErr(err)
	}
	return usr, tkn, nil
}
