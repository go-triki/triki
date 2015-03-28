/*
Package conf is a gotriki configuration package.

It is resposible for loading configuration file and parsing command line flags.
*/
package main

import (
	"crypto/tls"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
	"time"

	"gopkg.in/go-kornel/go-toml-config.v0"
	"gopkg.in/mgo.v2"
)

//general config
var (
	optShowConf bool
	optConfFile string
)

// server config
var (
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
	config.CommandLine = &config.Set{flag.CommandLine}
	// general config
	config.BoolVar(&optShowConf, "show_config", false, "print currently loaded configuration and exit")
	config.StringVar(&optConfFile, "config", "", "`path` to a TOML configuration file")

	// server config
	config.StringVar(&optServRoot, "server.root", "./www", "directory with static files to serve")
	config.StringVar(&server.Addr, "server.addr", "localhost:8765", "address and port to serve on")

	// mongo config
	config.StringVar(&optMongoAddrs, "mongo.Addrs",
		"localhost:27017", "MongoDB server to connect to, format: host1[:port1][,host2[:port2],...]")
	config.BoolVar(&optMongoSSL, "mongo.SSL",
		true, "use SSL for connections with MongoDB server")
	config.BoolVar(&optMongoSSLInsecure, "mongo.SSLInsecure",
		false, "don't verify MongoDB server's certificates, suspectible to man-in-the-middle attack, insecure!")
	config.BoolVar(&mDialInfo.Direct, "mongo.Direct",
		false, "direct connection with MongoDB?")
	config.DurationVar(&mDialInfo.DialTimeout, "mongo.DialTimeout",
		5*time.Second, "timeout for connecting to MongoDB instance (must be >=0)")
	config.StringVar(&mDialInfo.Database, "mongo.Database",
		"triki", "MongoDB database with triki data")
	config.String(&mDialInfo.Username, "mongo.Usr",
		"triki", "username for authentication to MongoDB")
	config.String(&mDialInfo.Password, "mongo.Pass",
		"triki", "password for authentication to MongoDB")

	// parse flags
	config.ParseArgs()
	// parse config file
	if optConfFile != "" {
		if err := config.Parse(optConfFile); err != nil {
			log.Fatalf("Error reading config file `%s`:\n%v", optConfFile, err)
		}
	}
	// write out option values?
	if optShowConf {
		config.PrintCurrentValues()
		os.Exit(0)
	}

	mDialInfo.Addrs = strings.Split(*optMongoAddrss, ",")

	if *optMongoSSL {
		mDialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), &tls.Config{InsecureSkipVerify: *optMongoSSLInsecure})
			if err != nil {
				log.Printf("MongoDB TLS connection error: %s.\n", err)
			}
			return conn, err
		}
	}

	if mDialInfo.Timeout < 0 {
		log.Fatalln("MongoDB dial timeout must be nonnegative.")
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
}
