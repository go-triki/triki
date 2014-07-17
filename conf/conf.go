/*
	Package conf is a gotriki configuration package.

	It is resposible for loading configuration file and parsing command line flags.
*/
package conf

import (
	"flag"
	"runtime"
)

// flag variables
var (
	flagHTTP     = flag.String("addr", "localhost:8765", "address and port for gotriki server")
	flagServRoot = flag.String("root", "./www", "directory with static files to serve")
)

type ServerOpts struct {
	Addr string
	Root string
}

// parsed flags
var (
	// main gotriki server options
	Server ServerOpts
)

//	Setup parses command line flags and config files.
func Setup() {
	flag.Parse()
	Server.Addr = *flagHTTP
	Server.Root = *flagServRoot
	runtime.GOMAXPROCS(runtime.NumCPU())
}
