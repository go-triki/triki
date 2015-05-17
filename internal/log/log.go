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
	// StdLog is a log.Logger that is going to be used by this package.
	StdLog           = log.New(stderr, "triki:", log.LstdFlags)
	StdLogUnbuffered = log.New(os.Stderr, "triki:", log.LstdFlags)
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
