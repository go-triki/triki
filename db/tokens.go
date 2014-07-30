package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

const (
	// tokensCName is the name of the tokens collection.
	tokensCName = "tokens"
	//tokensExpireAfter time.Duration = time.Duration(time.Second * 60 * 60 * 24 * 7) // 1 week
	// tokensExpireAfter controls how long authentication tokens are valid.
	tokensExpireAfter time.Duration = time.Duration(time.Second * 60 * 60) // 1h
)

type (
	// Token holds information associated with a given authentication token
	Token struct {
		ID        bson.ObjectId `bson:"_id"`
		CreatedAt time.Time     `bson:"createdAt"`
		UserID    bson.ObjectId `bson:"userID"`
	}
)

// tokensC returns the tokens collection.
func tokensC() *mgo.Collection {
	return adminSession.Copy().DB("").C(tokensCName)
}

func tokensSetup() {
	c := tokensC()
	index := mgo.Index{
		Key:         []string{"createdAt"},
		Unique:      false,
		DropDups:    false,
		Background:  false,
		Sparse:      false,
		ExpireAfter: tokensExpireAfter,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		log.Fatalf("MongoDB ensureIndex createdAt on tokens failed: %s.\n", err)
	}
	usrIndex := mgo.Index{
		Key:        []string{"userID"},
		Unique:     false,
		DropDups:   false,
		Background: false,
		Sparse:     false,
	}
	err = c.EnsureIndex(usrIndex)
	if err != nil {
		log.Fatalf("MongoDB ensureIndex userID on tokens failed: %s.\n", err)
	}
}

// tokenNew creates new token for user usrID.
// Returns token string and error message.
func tokenNew(usrID bson.ObjectId) (string, error) {
	var token Token
	token.ID = bson.NewObjectId()
	token.CreatedAt = time.Now()
	token.UserID = usrID
	err := tokensC().Insert(token)
	return token.ID.Hex(), err
}

// TokenCheck check if the given token is valid.
// Returns authenticated user and error.
func TokenCheck(tknID bson.ObjectId) (bson.ObjectId, error) {
	var token Token
	err := tokensC().Find(bson.M{"_id": tknID}).One(&token)
	return token.UserID, err
}
