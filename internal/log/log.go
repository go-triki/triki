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
	stderr = bufio.NewWriter(os.Stderr)
	// StdLog is a log.Logger that is going to be used by this package.
	StdLog = log.New(stderr, "triki", log.LstdFlags)
	// DBLog is a function used by this package to write logs to a database.
	DBLog func(map[string]interface{}) error
)

func init() {
	log.SetOutput(stderr)
}

// Flush flushes StdLog to ensure all messages went through.
func Flush() {
	stderr.Flush()
}
