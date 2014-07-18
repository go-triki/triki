/*
Package conf is a gotriki configuration package.

It is resposible for loading configuration file and parsing command line flags.
*/
package conf

import (
	"flag"
	"log"
	"runtime"
	"time"
)

// flag variables
var (
	flagHTTP             = flag.String("addr", "localhost:8765", "address and port for gotriki server")
	flagServRoot         = flag.String("root", "./www", "directory with static files to serve")
	flagMongoAddr        = flag.String("mongo", "triki:triki@localhost:27017/triki", "MongoDB server to connect to, format: [user:pass@]host1[:port1][,host2[:port2],...][/database][?options]")
	flagMongoDialTimeout = flag.Int("mongoDialTimeout", 10, "timeout for connecting to MongoDB instance (in seconds, >=0)")
)

// ServerOpts stores main gotriki server options
type ServerOpts struct {
	Addr string
	Root string
}

// MongoOpts stores options of connection to the MongoDB
type MongoOpts struct {
	Addr        string
	DialTimeout time.Duration
}

// parsed flags
var (
	// Server stores main gotriki server options
	Server ServerOpts
	// MongoDB server options
	MongoDB MongoOpts
)

// Setup parses command line flags and config files.
func Setup() {
	flag.Parse()

	Server.Addr = *flagHTTP
	Server.Root = *flagServRoot

	MongoDB.Addr = *flagMongoAddr
	MongoDB.DialTimeout = time.Duration(*flagMongoDialTimeout) * time.Second
	if MongoDB.DialTimeout < 0 {
		log.Fatalln("MongoDB dial timeout must be nonnegative.")
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
}
