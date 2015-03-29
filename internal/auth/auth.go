/*
Package auth deals with authentication, authorization, logging-in/out and
signing-up.
*/
package auth // import "gopkg.in/triki.v0/internal/auth"

import (
	"time"

	"golang.org/x/net/context"
	"gopkg.in/triki.v0/internal/log"
)

var (
	// RequestTimeout is (weak) time limit for HTTP request processing.
	RequestTimeout time.Duration
)

// DB session management
var (
	DBNewSession     func(context.Context) *log.Error
	DBNewAuthSession func(context.Context) *log.Error
)
