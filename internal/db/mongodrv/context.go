package mongo

import (
	"errors"
	"fmt"
	"net/http"

	gcontext "github.com/gorilla/context"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/triki.v0/internal/ctx"
	"gopkg.in/triki.v0/internal/log"
)

type ctxKey int

const sessKey ctxKey = 1

// SessionFromReq retrieves DB session associated with the request.
func SessionFromReq(r *http.Request) (*mgo.Session, *log.Error) {
	sess, ok := gcontext.GetOk(r, sessKey)
	if !ok {
		return nil, log.InternalServerErr(
			fmt.Errorf("couldn't find session associated with request %v", *r))
	}
	s, ok := sess.(*mgo.Session)
	if !ok {
		return nil, log.InternalServerErr(fmt.Errorf(
			"couldn't find session (type mismatch) associated with request %v", *r))
	}
	return s, nil
}

// Session retrieves DB session associated with the context.
func Session(c context.Context) (*mgo.Session, *log.Error) {
	req, ok := ctx.HTTPRequest(c)
	if !ok {
		return nil, log.InternalServerErr(errors.New(
			"couldn't find request associated with context"))
	}
	return SessionFromReq(req)
}

// SetSession saves DB session associated with the context.
func setSession(c context.Context, sess *mgo.Session) *log.Error {
	req, ok := ctx.HTTPRequest(c)
	if !ok {
		return log.InternalServerErr(errors.New(
			"couldn't associate *mgo.Session with a context.Context"))
	}
	gcontext.Set(req, sessKey, sess)
	return nil
}

// SaveSession creates a new DB session of a given type and associates it
// with the context. The session is typically closed by auth.Handler.
func SaveSession(c context.Context, typ ctx.SessionType) *log.Error {
	var sess *mgo.Session
	switch typ {
	case ctx.RegularSession:
		sess = session.Copy()
	case ctx.AdminSession:
		sess = adminSession.Copy()
	default:
		return log.InternalServerErr(fmt.Errorf("unknown session type %v", typ))
	}
	return setSession(c, sess)
}
