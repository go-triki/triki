/*
Package db wraps MongoDB database for triki.
*/
package mongo // import "gopkg.in/triki.v0/internal/db/mongodrv"

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/triki.v0/internal/conf"
)

// MongoDB collections' names
const (
	usersCName  = "users"
	tokensCName = "tokens"
)

var (
	session, adminSession, logSession *mgo.Session
)

// Setup database connections, etc.
func Setup() {
	var err error
	session, err = mgo.DialWithInfo(&conf.MDBDialInfo)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %s.\n", err)
	}
	session.SetSafe(&mgo.Safe{})
	session.SetMode(mgo.Monotonic, true)

	adminSession = session.Copy()
	adminSession.SetSafe(&mgo.Safe{WMode: "majority", FSync: true})
	adminSession.SetMode(mgo.Strong, true)

	// TODO
	logSession = session.Copy()
	//logSession.SetMode(mgo.)

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
