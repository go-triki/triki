package mongo

import (
	"time"

	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/triki.v0/internal/log"
	"gopkg.in/triki.v0/internal/models/token"
)

const (
	//tokensExpireAfter time.Duration = time.Duration(time.Second * 60 * 60 * 24 * 7) // 1 week
	// tokensExpireAfter controls how long authentication tokens are valid.
	tokensExpireAfter time.Duration = time.Duration(time.Second * 60 * 60) // 1h
)

// tokensC returns the tokens collection.
func tokensC(cx context.Context) (*mgo.Collection, *log.Error) {
	sess, err := getSession(cx, adminSessKey)
	if err != nil {
		return nil, err
	}
	return sess.DB("").C(tokensCName), nil
}

// TokenFind finds given token in the DB.
func TokenFind(cx context.Context, tknID []byte) (*token.T, *log.Error) {
	var tkn token.T
	col, er := tokensC(cx)
	if er != nil {
		return nil, er
	}
	err := col.Find(bson.M{"_id": tknID}).One(&tkn)
	if err != nil {
		return nil, log.InternalServerErr(err)
	}
	return &tkn, nil
}

// TokenExists checks if a token exists in the DB.
func TokenExists(cx context.Context, token []byte) (bool, *log.Error) {
	col, er := tokensC(cx)
	if er != nil {
		return false, er
	}
	n, err := col.FindId(token).Count()
	if err != nil {
		return false, log.InternalServerErr(err)
	}
	return n >= 1, nil
}

// TokenInsert inserts the token into the DB.
func TokenInsert(cx context.Context, tkn *token.T) *log.Error {
	col, er := tokensC(cx)
	if er != nil {
		return er
	}
	err := col.Insert(tkn)
	if err != nil {
		return log.InternalServerErr(err)
	}
	return nil
}
