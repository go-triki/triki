package mongo

import (
	"time"

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
func tokensC() *mgo.Collection {
	return adminSession.Copy().DB("").C(tokensCName)
}

// TokenFind finds given token in the DB.
func TokenFind(tknID []byte) (*token.T, *log.Error) {
	var token token.T
	err := tokensC().Find(bson.M{"_id": tknID}).One(&token)
	if err != nil {
		return nil, log.InternalServerErr(err)
	}
	return &token, nil
}

// TokenInsert inserts the token into the DB.
func TokenInsert(tkn *token.T) *log.Error {
	err := tokensC().Insert(tkn)
	if err != nil {
		return log.InternalServerErr(err)
	}
	return nil
}
