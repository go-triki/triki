// Package ctx provides a context.Context implementation whose Value
// method returns the values associated with a specific HTTP request in the
// github.com/gorilla/context package.
//
// Original at https://blog.golang.org/context/gorilla/gorilla.go
package ctx

import (
	"errors"
	"fmt"
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

const (
	reqKey  key = 0
	sessKey key = 1
)

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

// DBSessionFromReq retrieves mongo session associated with the request.
func DBSessionFromReq(r *http.Request) (*mgo.Session, *log.Error) {
	sess, ok := gcontext.GetOk(r, sessKey)
	if !ok {
		return nil, log.InternalServerErr(errors.New(
			fmt.Sprintf("couldn't find session associated with request %v", *r)))
	}
	s, ok := sess.(*mgo.Session)
	if !ok {
		return nil, log.InternalServerErr(errors.New(
			fmt.Sprintf("couldn't find session (type mismatch) associated with request %v", *r)))
	}
	return s, nil
}

// DBSession retrieves mongo session associated with the context.
func DBSession(c context.Context) (*mgo.Session, *log.Error) {
	req, ok := HTTPRequest(c)
	if !ok {
		return nil, log.InternalServerErr(errors.New("couldn't find request associated with context"))
	}
	return DBSessionFromReq(req)
}

// SetDBSession saves mongo session associated with the context.
func SetDBSession(c context.Context, sess *mgo.Session) bool {
	req, ok := HTTPRequest(c)
	if !ok {
		return false
	}
	gcontext.Set(req, sessKey, sess)
	return true
}
