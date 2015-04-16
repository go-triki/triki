/*
Package auth deals with authentication, authorization, logging-in/out and
signing-up.
*/
package auth // import "gopkg.in/triki.v0/internal/auth"

import (
	"net/http"
	"time"

	gctx "github.com/gorilla/context"
	"golang.org/x/net/context"
	"gopkg.in/triki.v0/internal/ctx"
	"gopkg.in/triki.v0/internal/log"
	"gopkg.in/triki.v0/internal/models/token"
	"gopkg.in/triki.v0/internal/models/user"
)

// options set by conf
var (
	// RequestTimeout is (weak) time limit for HTTP request processing.
	RequestTimeout time.Duration
)

// set by DB driver
var (
	// DBCloseSessions closes all sessions associated with the request created
	// by the DB driver.
	DBCloseSessions func(*http.Request)
)

// set by www
var (
	// WriteErrorHandler writes given error to the http respose and saves it in
	// the logs via context.
	WriteErrorHandler func(context.Context, http.ResponseWriter, *http.Request, *log.Error)
)

// T (auth.T) knows what given user is authorized to do.
type T struct {
	Usr *user.T
	Tkn *token.T
}

type authkey int

const authKey authkey = 0

// Set associates given auth info with context (http request).
func Set(cx context.Context, ath *T) {
	req, _ := ctx.HTTPRequest(cx)
	if req == nil {
		return
	}
	gctx.Set(req, authKey, ath)
}

// Get auth info associated with context (http request).
func Get(cx context.Context) *T {
	req, _ := ctx.HTTPRequest(cx)
	if req == nil {
		return &T{}
	}
	ah := gctx.Get(req, authKey)
	ath, _ := ah.(*T)
	if ath == nil {
		return &T{}
	}
	return ath
}
