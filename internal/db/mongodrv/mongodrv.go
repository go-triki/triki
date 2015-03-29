/*
Package mongo wraps MongoDB database for triki.
*/
package mongo // import "gopkg.in/triki.v0/internal/db/mongodrv"

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/triki.v0/internal/models/token"
	"gopkg.in/triki.v0/internal/models/user"
)

// MongoDB collections' names
const (
	usersCName  = "users"
	tokensCName = "tokens"
)

var (
	// DialInfo used to connect MongoBD.
	DialInfo = mgo.DialInfo{}
)

var (
	session, adminSession, logSession *mgo.Session
)

// Setup database connections, etc.
func Setup() {
	var err error
	session, err = mgo.DialWithInfo(&DialInfo)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %s.\n", err)
	}
	session.SetSafe(&mgo.Safe{})
	session.SetMode(mgo.Monotonic, true)

	logSession = session.Copy()
	logSession.SetSafe(nil)
	logSession.SetMode(mgo.Eventual, false)

	adminSession = session.Copy()
	adminSession.SetSafe(&mgo.Safe{WMode: "majority", FSync: true})
	adminSession.SetMode(mgo.Strong, true)

	usersSetup()
	tokensSetup()
}

// Cleanup database connections, etc.
func Cleanup() {
	log.Println("Closing database connections...")
	adminSession.Close()
	session.Close()
	logSession.Close()
	log.Println("Database connections closed.")
}

// usersSetup ensures that the users collection is setup correctly.
func usersSetup() {
	c := usersC()
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
	c := tokensC()
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
