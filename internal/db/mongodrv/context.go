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

type sessionKey int

const (
	sessKey sessionKey = iota
	adminSessKey
	numOfSessKeys
)

// TODO check
func CloseSessions(req *http.Request) {
	for i := sessionKey(0); i < numOfSessKeys; i++ {
		val, ok := gcontext.GetOk(req, i)
		if ok {
			sess, ok := val.(*mgo.Session)
			if ok {
				sess.Close()
			}
		}
	}
}

func getSession(cx context.Context, typ sessionKey) (*mgo.Session, *log.Error) {
	req, ok := ctx.HTTPRequest(cx)
	if !ok {
		return nil, log.InternalServerErr(errors.New(
			"couldn't find request associated with context"))
	}
	if s, ok := gcontext.GetOk(req, typ); ok {
		return s.(*mgo.Session), nil
	}
	var sess *mgo.Session
	switch typ {
	case sessKey:
		sess = session.Copy()
	case adminSessKey:
		sess = adminSession.Copy()
	default:
		return nil, log.InternalServerErr(fmt.Errorf("unknown session type %v", typ))
	}
	gcontext.Set(req, typ, sess)
	return sess, nil
}
