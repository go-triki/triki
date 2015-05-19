/*
Package log is a gotriki logging package.
*/
package log // import "gopkg.in/triki.v0/internal/log"

import (
	"bufio"
	"log"
	"os"
)

var (
	// stderr is a buffered std err for log package output.
	stderr = bufio.NewWriter(os.Stderr)
	// StdLog is triki www unbuffered log.Logger.
	//StdLog = log.New(stderr, "triki:", log.LstdFlags) // TODO(km): buffer it?
	StdLog = log.New(os.Stderr, "triki-www:", log.LstdFlags)
	// StdLogMongo is an unbuffered log.Logger used to log mongo-related info.
	StdLogMongo = log.New(os.Stderr, "triki-mongo:", log.LstdFlags)
	// DBLog is a function used by this package to write logs to a database.
	DBLog func(map[string]interface{}) error
)

func init() {
	log.SetOutput(os.Stderr)
	log.SetPrefix("triki:")
	log.SetFlags(log.LstdFlags)
}

// Flush flushes StdLog to ensure all messages went through.
func Flush() {
	stderr.Flush()
}
