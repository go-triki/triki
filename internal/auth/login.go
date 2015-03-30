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
		// TODO return an error
	}
	// TODO check recent failed logins
	if err := bcrypt.CompareHashAndPassword(usr.PassHash, append([]byte(user.PassSalt+pass), usr.Salt...)); err != nil {
		// TODO no more than 3 failed logins in 10 minutes? (use req to get info)
		return nil, nil, log.IncorrectPassErr
	}
	tkn := &token.T{UsrID: usr.ID}
	err = token.New(cx, tkn)
	if err != nil {
		return nil, nil, log.InternalServerErr(err)
	}
	return usr, tkn, nil
}
