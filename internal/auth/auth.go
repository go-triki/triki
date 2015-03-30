/*
Package auth deals with authentication, authorization, logging-in/out and
signing-up.
*/
package auth // import "gopkg.in/triki.v0/internal/auth"

import (
	"net/http"
	"time"
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
