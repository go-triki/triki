/*
Package db wraps MongoDB database for triki.
*/
package db

import (
	"gopkg.in/triki.v0/conf"
	"gopkg.in/mgo.v2"
	"log"
)

// MongoDB collections' names
const (
	usersCName  = "users"
	tokensCName = "tokens"
)

var (
	session, adminSession *mgo.Session
)

// options:
var (
	// prefix passwords before hashing
	usersPassPrefix string = "" // TODO: set this up from conf
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

	usersSetup()
	tokensSetup()
}

// Cleanup database connections, etc.
func Cleanup() {
	log.Println("Closing database connections...")
	adminSession.Close()
	session.Close()
	log.Println("Database connections closed.")
}
