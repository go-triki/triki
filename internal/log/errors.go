package log

import (
	"fmt"
	"net/http"
)

type trikiCode int

// triki codec (TCs)
const (
	incorrectPassTC    trikiCode = 100
	badSignupDetailsTC trikiCode = 200
	badTokenTC         trikiCode = 250
	dbNotFoundTC       trikiCode = 300
)

// Error instances represent error encountered while serving www requests.
type Error struct {
	What       string    // human-readable decription of the error
	TrikiCode  trikiCode // error code passed to triki web interface
	HTTPStatus int       // HTTP status of the reply (0 - don't modify)
}

func (err Error) Error() string {
	return err.What
}

// InternalServerErr returns an error indicating some unexpected condition
// (e.g. DB faliure).
func InternalServerErr(err error) *Error {
	return &Error{
		What:       fmt.Sprintf("internal server error: %s", err.Error()),
		HTTPStatus: http.StatusInternalServerError,
	}
}

// BadSignupDetailsErr returns an error indicating that login/password/nick
// supplied by the user don't conform to triki's standards.
func BadSignupDetailsErr(detail string) *Error {
	return &Error{
		What:      detail,
		TrikiCode: badSignupDetailsTC,
	}
}

// DBNotFoundErr returns an error indicating that either requested item in not
// in the DB or there was a DB error.
func DBNotFoundErr(err error) *Error {
	return &Error{
		What:      err.Error(),
		TrikiCode: dbNotFoundTC,
	}
}

var (
	// IncorrectPassErr is returned when the password supplied by the user
	// doesn't match the one recorded in the DB.
	IncorrectPassErr = &Error{"incorrect password", incorrectPassTC, 0}
	// BadTokenErr is returned when the token supplied by the user
	// is expired/invalid/not in the DB.
	BadTokenErr = &Error{"bad authorization token", badTokenTC, 0}
)
