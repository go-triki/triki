/*
Package auth deals with authentication, authorization, logging-in/out and
signing-up.
*/
package auth // import "gopkg.in/triki.v0/internal/auth"

import "time"

var (
	// RequestTimeout is (weak) time limit for HTTP request processing.
	RequestTimeout time.Duration
)
