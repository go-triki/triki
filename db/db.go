/*
Package db wrapps MongoDB database for triki.
*/
package db

import (
	"bitbucket.org/kornel661/triki/gotriki/conf"
	"labix.org/v2/mgo"
	"log"
)

var (
	session, adminSession *mgo.Session
)

// Setup database connections, etc.
func Setup() {
	//session, err := mgo.DialWithTimeout(conf.MongoDB.Addr, conf.MongoDB.DialTimeout)
	session, err := mgo.DialWithInfo(&conf.MDBDialInfo)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %s\n", err)
	}
	adminSession = session.Clone()
	session.SetMode(mgo.Monotonic, true)
	adminSession.SetMode(mgo.Strong, true)
}

// Cleanup database connections, etc.
func Cleanup() {
	log.Println("Closing database connections...")
	adminSession.Close()
	session.Close()
	log.Println("Database connections closed.")
}
