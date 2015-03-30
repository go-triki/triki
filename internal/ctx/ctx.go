// Package ctx provides a context.Context implementation whose Value
// method returns the values associated with a specific HTTP request in the
// github.com/gorilla/context package.
//
// Additionally, access to http.Requests and mgo.Sessions is provided.
//
// Original at https://blog.golang.org/context/gorilla/gorilla.go
package ctx

import (
	"net/http"

	gcontext "github.com/gorilla/context"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/triki.v0/internal/log"
)

// New returns a Context whose Value method returns values associated
// with req using the Gorilla context package:
// http://www.gorillatoolkit.org/pkg/context
func New(parent context.Context, req *http.Request) context.Context {
	return &wrapper{parent, req}
}

type wrapper struct {
	context.Context
	req *http.Request
}

type key int

const reqKey key = 0

// Value returns Gorilla's context package's value for this Context's request
// and key. It delegates to the parent Context if there is no such value.
func (ctx *wrapper) Value(key interface{}) interface{} {
	if key == reqKey {
		return ctx.req
	}
	if val, ok := gcontext.GetOk(ctx.req, key); ok {
		return val
	}
	return ctx.Context.Value(key)
}

// HTTPRequest returns the *http.Request associated with ctx using NewContext,
// if any.
func HTTPRequest(ctx context.Context) (*http.Request, bool) {
	// We cannot use ctx.(*wrapper).req to get the request because ctx may
	// be a Context derived from a *wrapper. Instead, we use Value to
	// access the request if it is anywhere up the Context tree.
	req, ok := ctx.Value(reqKey).(*http.Request)
	return req, ok
}

// SessionType conveys information on a type of session to create.
type SessionType int

// Session Types
const (
	RegularSession SessionType = iota // session with standard consistency requirements
	AdminSession                      // session with higher consistency requirements
)

var (
	// DBSessionFromReq retrieves DB session associated with the request.
	DBSessionFromReq func(r *http.Request) (*mgo.Session, *log.Error)
	// DBSession retrieves DB session associated with the context.
	DBSession func(c context.Context) (*mgo.Session, *log.Error)
	// DBSaveSession creates a new DB session of a given type and associates it
	// with the context. The session is typically closed by auth.Handler.
	DBSaveSession func(context.Context, SessionType) *log.Error
)
