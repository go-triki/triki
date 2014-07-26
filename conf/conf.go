/*
Package conf is a gotriki configuration package.

It is resposible for loading configuration file and parsing command line flags.
*/
package conf

import (
	"crypto/tls"
	"flag"
	"gopkg.in/mgo.v2"
	"log"
	"net"
	"runtime"
	"strings"
	"time"
)

// flag variables
var (
	flagHTTP             = flag.String("addr", "localhost:8765", "address and port for gotriki server")
	flagServRoot         = flag.String("root", "./www", "directory with static files to serve")
	flagMongoAddr        = flag.String("mongoAddr", "localhost:27017", "MongoDB server to connect to, format: host1[:port1][,host2[:port2],...]")
	flagMongoDirect      = flag.Bool("mongoDirect", false, "direct connection with MongoDB?")
	flagMongoDialTimeout = flag.Int("mongoDialTimeout", 5, "timeout for connecting to MongoDB instance (in seconds, >=0)")
	flagMongoDatabase    = flag.String("mongoDatabase", "triki", "MongoDB database with triki data")
	flagMongoUsr         = flag.String("mongoUsr", "triki", "username for authentication to MongoDB")
	flagMongoPass        = flag.String("mongoPass", "triki", "password for authentication to MongoDB")
	flagMongoSSL         = flag.Bool("mongoSSL", true, "use SSL for connections with MongoDB server")
	flagMongoSSLInsecure = flag.Bool("mongoSSLInsecure", false, "don't verify MongoDB server's certificates, suspectible to man-in-the-middle attack, insecure!")
)

// ServerOpts stores main gotriki server options.
type ServerOpts struct {
	Addr string
	Root string
}

// parsed flags
var (
	// Server stores main gotriki server options.
	Server ServerOpts
	// MDBDialInfo holds MongoDB server connection options.
	MDBDialInfo mgo.DialInfo
)

// Setup parses command line flags and config files.
func Setup() {
	flag.Parse()

	MDBDialInfo.Addrs = strings.Split(*flagMongoAddr, ",")
	MDBDialInfo.Direct = *flagMongoDirect
	MDBDialInfo.Timeout = time.Duration(*flagMongoDialTimeout) * time.Second
	MDBDialInfo.Database = *flagMongoDatabase
	MDBDialInfo.Username = *flagMongoUsr
	MDBDialInfo.Password = *flagMongoPass
	if *flagMongoSSL {
		MDBDialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), &tls.Config{InsecureSkipVerify: *flagMongoSSLInsecure})
			if err != nil {
				log.Printf("MongoDB TLS connection error: %s.\n", err)
			}
			return conn, err
		}
	}

	Server.Addr = *flagHTTP
	Server.Root = *flagServRoot

	if MDBDialInfo.Timeout < 0 {
		log.Fatalln("MongoDB dial timeout must be nonnegative.")
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
}
