/*
Package mongo wraps MongoDB database for triki.
*/
package mongo // import "gopkg.in/triki.v0/internal/db/mongodrv"

import (
	"log"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/triki.v0/internal/auth"
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
	// MaxLogSize is maximum size (in bytes) of the log DB.
	MaxLogSize int
)

var (
	session, adminSession, logSession *mgo.Session
)

// Setup database connections, etc. Run by main program.
func Setup() {
	mgo.SetLogger(tlog.StdLogMongo)
	mgo.SetDebug(false)
	mgo.SetStats(false)
	var err error
	session, err = mgo.DialWithInfo(&DialInfo)
	if err != nil {
		log.Printf("Error connecting to MongoDB: %s.\n", err)
		tlog.Flush()
		os.Exit(1)
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
	auth.DBCloseSessions = CloseSessions
	// setup collections
	logSetup()
	usersSetup()
	tokensSetup()
}

// Cleanup database connections, etc. Run by main program.
func Cleanup() {
	log.Println("Closing database connections...")
	adminSession.Close()
	session.Close()
	logSession.Close()
	log.Println("Database connections closed.")
}
