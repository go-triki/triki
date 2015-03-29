package auth

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/triki.v0/internal/log"
	"gopkg.in/triki.v0/internal/models/token"
	"gopkg.in/triki.v0/internal/models/user"
)

// Authenticate checks if the given credentials can be authenticated.
// Returns (user, token, error).
func Authenticate(usrID, pass string, req *http.Request) (*user.T, *token.T, *log.Error) {
	usr, err := user.DBFind(usrID)
	if err != nil {
		return nil, nil, err
	}
	// TODO check recent failed logins
	if err := bcrypt.CompareHashAndPassword(usr.PassHash, append([]byte(user.PassSalt+pass), usr.Salt...)); err != nil {
		// TODO no more than 3 failed logins in 10 minutes? (use req to get info)
		return nil, nil, log.IncorrectPassErr
	}
	tkn, err := token.New(usr.ID, req)
	if err != nil {
		return nil, nil, log.InternalServerErr(err)
	}
	return usr, tkn, nil
}
