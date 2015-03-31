/*
Package conf is a gotriki configuration package.

It is resposible for loading configuration file and parsing command line flags.
*/
package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
	"time"

	"gopkg.in/go-kornel/go-toml-config.v0"
	"gopkg.in/mgo.v2"
	"gopkg.in/triki.v0/internal/auth"
	"gopkg.in/triki.v0/internal/db/mongodrv"
	"gopkg.in/triki.v0/internal/models/token"
	"gopkg.in/triki.v0/internal/models/user"
)

// general config
var (
	optShowConf bool
	optMoreHelp bool
	optConfFile string
	optNumCpus  int
)

// server config
var (
	// directory with static files to serve
	optServRoot string
)

// mongo config
var (
	optMongoAddrs       string
	optMongoSSL         bool
	optMongoSSLInsecure bool
)

// init parses command line flags and config files.
func init() {
	//config.CommandLine = &config.Set{flag.CommandLine}
	// general config
	flag.BoolVar(&optShowConf, "show_config", false, "print currently loaded configuration and exit")
	flag.BoolVar(&optMoreHelp, "more_help", false, "print help for more triki options")
	flag.StringVar(&optConfFile, "config", "", "`path` to a TOML configuration file")

	config.StringVar(&user.PassSalt, "pass_salt",
		"", "used to `salt` passwords in the DB")
	config.IntVar(&optNumCpus, "num_cpus",
		0, "number of CPUs to use, 0 to autodetect")

	// triki config
	config.DurationVar(&token.MaxExpireAfter, "triki.tkn_max_expire_after",
		7*24*time.Hour, "maximum time a user authorization token is valid")

	// server config
	config.StringVar(&optServRoot, "server.root",
		"./www", "directory with static files to serve")
	config.StringVar(&server.Addr, "server.addr",
		"localhost:8765", "address and port to serve on")
	config.DurationVar(&auth.RequestTimeout, "server.request_timeout",
		3*time.Second, "(weak) time limit for HTTP request processing")

	// mongo config
	config.IntVar(&mongo.MaxLogSize, "mongo.max_log_size",
		1e10, "maximum size (in bytes) of the log database")
	config.StringVar(&optMongoAddrs, "mongo.Addrs",
		"localhost:27017", "MongoDB server to connect to, format: host1[:port1][,host2[:port2],...]")
	config.BoolVar(&optMongoSSL, "mongo.SSL",
		true, "use SSL for connections with MongoDB server")
	config.BoolVar(&optMongoSSLInsecure, "mongo.SSLInsecure",
		false, "don't verify MongoDB server's certificates, suspectible to man-in-the-middle attack, insecure!")
	config.BoolVar(&mongo.DialInfo.Direct, "mongo.Direct",
		false, "direct connection with MongoDB (don't connect with the whole cluster)")
	config.DurationVar(&mongo.DialInfo.Timeout, "mongo.DialTimeout",
		5*time.Second, "timeout for connecting to MongoDB instance (must be >=0)")
	config.StringVar(&mongo.DialInfo.Database, "mongo.Database",
		"triki", "MongoDB database with triki data")
	config.StringVar(&mongo.DialInfo.Username, "mongo.Usr",
		"triki", "username for authentication to MongoDB")
	config.StringVar(&mongo.DialInfo.Password, "mongo.Pass",
		"triki", "password for authentication to MongoDB")
	////////////////////////////////////////////////////////////////////////////
	// parse flags
	flag.Parse()
	if optMoreHelp {
		config.CommandLine.PrintDefaults()
		os.Exit(0)
	}
	// parse config file
	if optConfFile != "" {
		if err := config.Parse(optConfFile); err != nil {
			log.Fatalf("Error reading config file `%s`:\n%v", optConfFile, err)
		}
	}
	config.ParseArgs()
	// write out option values?
	if optShowConf {
		config.PrintCurrentValues()
		os.Exit(0)
	}
	////////////////////////////////////////////////////////////////////////////
	// general config
	if user.PassSalt == "" {
		log.Fatalln("Error: `pass_salt` option can't be empty. Best practice is to set it to some random string.")
	}
	if optNumCpus == 0 {
		optNumCpus = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(optNumCpus)
	//server config
	mongo.DialInfo.Addrs = strings.Split(optMongoAddrs, ",")
	// mongo config
	if optMongoSSL {
		mongo.DialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), &tls.Config{InsecureSkipVerify: optMongoSSLInsecure})
			if err != nil {
				log.Printf("MongoDB TLS connection error: %s.\n", err.Error())
			}
			return conn, err
		}
	}

	if mongo.DialInfo.Timeout < 0 {
		log.Fatalln("MongoDB dial timeout `mongo.DialTimeout` must be nonnegative.")
	}
}
