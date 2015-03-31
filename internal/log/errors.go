package log

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	golangctx "golang.org/x/net/context"
	"gopkg.in/triki.v0/internal/ctx"
)

type errkey int

// errKey - a key to retreive errors. from context.
const errKey errkey = 0

type trikiCode int

// triki codes (TCs), for errors >0
const (
	incorrectPassTC        trikiCode = 100
	userNotActiveTC        trikiCode = 150
	badSignupDetailsTC     trikiCode = 200
	badTokenTC             trikiCode = 250
	dbNotFoundTC           trikiCode = 300
	internalServerErrTC    trikiCode = 400
	failedWritingReplyTC   trikiCode = 420
	failedReadingRequestTC trikiCode = 460
	invalidIDTC            trikiCode = 500
)

// Error instances represent error encountered while serving www requests.
type Error struct {
	What       string        // human-readable decription of the error
	TrikiCode  trikiCode     // error code passed to triki web interface
	HTTPStatus int           // HTTP status of the reply (0 - don't modify)
	Info       []interface{} // additional info
}

func (err Error) Error() string {
	return err.What
}

// Set associates err with the context cx so that it can be later retreived
// for logging purposes.
func Set(cx golangctx.Context, err *Error) {
	req, _ := ctx.HTTPRequest(cx)
	if req == nil {
		return
	}
	context.Set(req, errKey, err)
}

// InternalServerErr returns an error indicating some unexpected condition
// (e.g. DB faliure).
func InternalServerErr(err error) *Error {
	return &Error{
		What:       fmt.Sprintf("internal server error: %s", err.Error()),
		TrikiCode:  internalServerErrTC,
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

// FailedWritingReplyErr error is returned when writing to http.ResponseWriter failed.
func FailedWritingReplyErr(err error) *Error {
	return &Error{
		What:      err.Error(),
		TrikiCode: failedWritingReplyTC,
	}
}

// FailedReadingRequestErr error is returned when reading from http.Request failed.
func FailedReadingRequestErr(err error) *Error {
	return &Error{
		What:      err.Error(),
		TrikiCode: failedReadingRequestTC,
	}
}

var (
	// IncorrectPassErr is returned when the password supplied by the user
	// doesn't match the one recorded in the DB.
	IncorrectPassErr = &Error{What: "incorrect password", TrikiCode: incorrectPassTC}
	// BadTokenErr is returned when the token supplied by the user
	// is expired/invalid/not in the DB.
	BadTokenErr = &Error{What: "bad authorization token", TrikiCode: badTokenTC}
	// UserNotActiveErr indicates that user account is not in the "active" state
	// and the user cannot log in.
	UserNotActiveErr = &Error{What: "this user account is not active", TrikiCode: userNotActiveTC}
	// InvalidIDErr is returned when ID (of a user, article, etc.) supplied in
	// http.Request is not a valid ID.
	InvalidIDErr = &Error{What: "requested ID is invalid", TrikiCode: invalidIDTC}
)
