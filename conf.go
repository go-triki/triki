/*
Package conf is a gotriki configuration package.

It is resposible for loading configuration file and parsing command line flags.
*/
package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"gopkg.in/go-kornel/go-toml-config.v0"
	"gopkg.in/mgo.v2"
	"gopkg.in/triki.v0/internal/auth"
	"gopkg.in/triki.v0/internal/db/mongodrv"
	tlog "gopkg.in/triki.v0/internal/log"
	"gopkg.in/triki.v0/internal/models/token"
	"gopkg.in/triki.v0/internal/tls_pass"
)

// general config
var (
	optShowConf bool
	optConfFile string
)

// server config
var (
	// directory with static files to serve
	optServRoot string
	// url to the server with static content (supersedes optServRoot)
	optServStaticServ string
	// url to the server with static resources (parsed optServStaticServ)
	staticServerURL *url.URL
)

// mongo config
var (
	optNumCpus          int
	optMongoAddrs       string
	optMongoSSL         bool
	optMongoSSLInsecure bool
	optMongoCAFile      string
	optMongoCertFile    string
	optMongoKeyFile     string
	optMongoSSLPass     string
)

// init parses command line flags and config files.
func init() {
	//config.CommandLine = &config.Set{flag.CommandLine}
	// general config
	config.BoolVar(&optShowConf, "show_config", false, "print currently loaded configuration and exit")
	config.StringVar(&optConfFile, "config", "", "`path` to a TOML configuration file")

	// triki config
	config.IntVar(&optNumCpus, "triki.num_cpus",
		0, "number of CPUs to use, 0 to autodetect")
	config.DurationVar(&token.MaxExpireAfter, "triki.tkn_max_expire_after",
		7*24*time.Hour, "maximum time a user authorization token is valid, change requres to rebuild indexes in the MongoDB")

	// server config
	config.StringVar(&optServRoot, "server.root",
		"./www", "directory with static files to serve")
	config.StringVar(&optServStaticServ, "server.static_server",
		"", "redirect requests for static contenct to this `server address` (e.g., http://localhost:8080/), supersedes server.root")
	config.StringVar(&server.Addr, "server.addr",
		"localhost:8765", "address and port to serve on")
	config.DurationVar(&auth.RequestTimeout, "server.request_timeout",
		3*time.Second, "(weak) time limit for HTTP request processing")

	// mongo config
	config.IntVar(&mongo.MaxLogSize, "mongo.max_log_size",
		1e10, "maximum size (in bytes) of the log database, change requres to rebuild indexes in the MongoDB")
	config.StringVar(&optMongoAddrs, "mongo.address",
		"localhost:27017", "MongoDB server to connect to, format: host1[:port1][,host2[:port2],...]")
	config.BoolVar(&optMongoSSL, "mongo.SSL",
		true, "use SSL for connections with MongoDB server")
	config.BoolVar(&optMongoSSLInsecure, "mongo.SSL_insecure",
		false, "don't verify MongoDB server's certificates, suspectible to man-in-the-middle attack, insecure!")
	config.StringVar(&optMongoCAFile, "mongo.SSL_CA_file",
		"", "`path` to file with CA SSL certificates, e.g., './CA.pem'")
	config.StringVar(&optMongoCertFile, "mongo.SSL_cert_file",
		"", "`path` to file with SSL certificate")
	config.StringVar(&optMongoKeyFile, "mongo.SSL_key_file",
		"", "`path` to file with SSL key")
	config.StringVar(&optMongoSSLPass, "mongo.SSL_password",
		"", "`password` for the encrypted SSL key")
	config.BoolVar(&mongo.DialInfo.Direct, "mongo.direct",
		false, "direct connection with MongoDB (don't connect with the whole cluster)")
	config.DurationVar(&mongo.DialInfo.Timeout, "mongo.dial_timeout",
		5*time.Second, "timeout for connecting to MongoDB instance (must be >=0)")
	config.StringVar(&mongo.DialInfo.Database, "mongo.database",
		"triki", "MongoDB database with triki data")
	config.StringVar(&mongo.DialInfo.Username, "mongo.user",
		"triki", "username for authentication to MongoDB")
	config.StringVar(&mongo.DialInfo.Password, "mongo.password",
		"triki", "password for authentication to MongoDB")
	////////////////////////////////////////////////////////////////////////////
	// parse flags
	config.ParseArgs()
	// parse config file
	if optConfFile != "" {
		if err := config.Parse(optConfFile); err != nil {
			log.Fatalf("Error reading config file `%s`:\n%v\n", optConfFile, err)
		}
	}
	////////////////////////////////////////////////////////////////////////////
	// general config
	if optNumCpus == 0 {
		optNumCpus = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(optNumCpus)

	//server config
	if url, err := url.Parse(optServStaticServ); err == nil {
		staticServerURL = url
	} else {
		log.Fatalf("Error parsing url `%s`:\n%v\n", optServStaticServ, err)
	}

	// mongo config
	mongo.DialInfo.Addrs = strings.Split(optMongoAddrs, ",")
	if optMongoSSL {
		// CA pool
		// TODO
		var CAPool *x509.CertPool
		if optMongoCAFile != "" {
			CAPool = x509.NewCertPool()
			cert, err := ioutil.ReadFile(optMongoCAFile)
			if err != nil {
				log.Fatalf("Error reading Mongo CA certificate `%s`: %v\n", optMongoCAFile, err)
			}
			if !CAPool.AppendCertsFromPEM(cert) {
				log.Fatalf("Error adding Mongo CA certificate `%s`.\n", optMongoCAFile)
			}
		}
		// certificate
		// TODO
		var cert tls.Certificate
		loadCert := optMongoCertFile != "" || optMongoKeyFile != ""
		if loadCert {
			var err error
			cert, err = tls_pass.LoadX509KeyPair(optMongoCertFile, optMongoKeyFile, []byte(optMongoSSLPass))
			if err != nil {
				log.Fatalf("Error loading certificate pair `%s`, `%s`: %v\n", optMongoCertFile, optMongoKeyFile, err)
			}
		}
		// TLS dial
		mongo.DialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			config := tls.Config{
				InsecureSkipVerify: optMongoSSLInsecure,
				RootCAs:            CAPool,
			}
			if loadCert {
				config.Certificates = []tls.Certificate{cert}
			}
			conn, err := tls.Dial("tcp", addr.String(), &config)
			if err != nil {
				tlog.StdLog.Printf("MongoDB TLS connection error: %s.\n", err.Error())
				tlog.Flush()
			}
			return conn, err
		}
	}

	if mongo.DialInfo.Timeout < 0 {
		log.Fatalln("MongoDB dial timeout `mongo.DialTimeout` must be nonnegative.")
	}

	////////////////////////////////////////////////////////////////////////////
	// write out option values?
	if optShowConf {
		config.PrintCurrentValues()
		tlog.Flush()
		os.Exit(0)
	}
}
