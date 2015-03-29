package mongo

import (
	"log"

	"gopkg.in/mgo.v2"
	tlog "gopkg.in/triki.v0/internal/log"
	"gopkg.in/triki.v0/internal/models/token"
	"gopkg.in/triki.v0/internal/models/user"
)

func logSetup() {
	c := adminSession.DB("").C(logCName)
	c.Create(&mgo.CollectionInfo{
		Capped:   true,
		MaxBytes: 1e6,
	})
	index := mgo.Index{
		Key:        []string{"time"},
		Unique:     false,
		DropDups:   false,
		Background: false,
		Sparse:     false,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		log.Fatalf("MongoDB ensureIndex `time` on log failed: %s.\n", err.Error())
	}
	// install DB functions
	tlog.DBLog = Log
}

// usersSetup ensures that the users collection is setup correctly.
func usersSetup() {
	c := adminSession.DB("").C(usersCName)
	index := mgo.Index{
		Key:        []string{"usr"},
		Unique:     true,
		DropDups:   false,
		Background: false,
		Sparse:     false,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		log.Fatalf("MongoDB ensureIndex `usr` on users failed: %s.\n", err.Error())
	}
	// install DB functions
	user.DBFind = UserFind
	user.DBInsert = UserInsert
	user.DBExists = UserExists
}

func tokensSetup() {
	c := adminSession.DB("").C(tokensCName)
	index := mgo.Index{
		Key:         []string{"birth"},
		Unique:      false,
		DropDups:    false,
		Background:  false,
		Sparse:      false,
		ExpireAfter: tokensExpireAfter,
	}
	if err := c.EnsureIndex(index); err != nil {
		log.Fatalf("MongoDB ensureIndex `birth` on tokens failed: %s.\n", err)
	}
	usrIndex := mgo.Index{
		Key:        []string{"usrID"},
		Unique:     false,
		DropDups:   false,
		Background: false,
		Sparse:     false,
	}
	if err := c.EnsureIndex(usrIndex); err != nil {
		log.Fatalf("MongoDB ensureIndex `usrID` on tokens failed: %s.\n", err)
	}
	// install DB functions
	token.DBFind = TokenFind
	token.DBInsert = TokenInsert
}
