/*
Package mongo wraps MongoDB database for triki.
*/
package mongo // import "gopkg.in/triki.v0/internal/db/mongodrv"

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/triki.v0/internal/ctx"
	tlog "gopkg.in/triki.v0/internal/log"
)

// MongoDB collections' names
const (
	usersCName  = "users"
	tokensCName = "tokens"
	logCName    = "log"
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
	mgo.SetLogger(tlog.StdLog)
	var err error
	session, err = mgo.DialWithInfo(&DialInfo)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %s.\n", err)
	}

	session.SetMode(mgo.Monotonic, true)
	session.SetSafe(&mgo.Safe{})

	logSession = session.Copy()
	logSession.SetMode(mgo.Eventual, true)
	logSession.SetSafe(nil)

	adminSession = session.Copy()
	adminSession.SetMode(mgo.Strong, true)
	adminSession.SetSafe(&mgo.Safe{WMode: "majority", FSync: true})

	// install DB functions
	ctx.DBSaveSession = SaveSession
	ctx.DBSession = Session
	ctx.DBSessionFromReq = SessionFromReq
	// setup collections
	logSetup()
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
