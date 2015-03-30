package auth

import (
	"fmt"

	"golang.org/x/net/context"
	"gopkg.in/triki.v0/internal/log"
	"gopkg.in/triki.v0/internal/models/user"
)

// UserSignup signs given user up, pending email verification.
// TODO write email verification
func UserSignup(cx context.Context, login, pass string) *log.Error {
	is, err := user.DBExists(cx, login)
	if err != nil {
		return log.InternalServerErr(err)
	} else if is {
		return log.BadSignupDetailsErr(fmt.Sprintf("user `%s` already exists", login))
	}
	usr := user.T{
		Usr:  login,
		Pass: pass,
	}
	// TODO add logging
	return user.New(cx, &usr)
}
