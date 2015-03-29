package auth

import (
	"fmt"
	"net/http"

	"gopkg.in/triki.v0/internal/log"
	"gopkg.in/triki.v0/internal/models/user"
)

// UserSignup signs given user up, pending email verification.
// TODO write email verification
func UserSignup(login, pass string, req *http.Request) *log.Error {
	is, err := user.DBExists(login)
	if err != nil {
		return log.InternalServerErr(err)
	} else if is {
		return log.BadSignupDetailsErr(fmt.Sprintf("User `%s` already exists", login))
	}
	usr := user.T{
		Usr:  login,
		Pass: pass,
	}
	// TODO add logging
	return user.New(&usr)
}
