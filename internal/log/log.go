/*
Package log is a gotriki logging package.
*/
package log // import "gopkg.in/triki.v0/internal/log"

import (
	"bufio"
	"log"
	"os"
)

// DbLogger is a type of functions that write logs to databases.
type DbLogger func(map[string]interface{}) error

var (
	stderr = bufio.NewWriter(os.Stderr)
	// StdLog is a log.Logger that is going to be used by this package.
	StdLog = log.New(stderr, "triki", LstdFlags)
	// DbLog is a function used by this package to write logs to a database.
	DbLog DbLogger
)

type errkey int

// ErrKey is a key to retreive error list from context.
const errKey errkey = 0

// Flush flushes StdLog to ensure all messages went through.
func Flush() {
	stderr.Flush()
}
