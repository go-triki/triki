/*
Package log is a gotriki logging package.
*/
package log

import (
	"log"
)

// Type that is "paniced" in fatal log functions
type FatalErrorPanic string

// Error returns error description.
func (fatal FatalErrorPanic) Error() string {
	return string(fatal)
}

// FatalErrorPanicType method makes this type unique.
func (fatal FatalErrorPanic) FatalErrorPanicType() {
	return
}

const (
	fatalErr FatalErrorPanic = "Fatal server error"
)

func Warningf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func Warningln(msg string) {
	Warningf("%s\n", msg)
}

func Infof(format string, args ...interface{}) {
	//glog.Infof(format, args...)
	log.Printf(format, args...)
}

func Infoln(msg string) {
	Infof("%s\n", msg)
}

func Fatalf(format string, args ...interface{}) {
	//glog.Fatalf(format, args...)
	log.Printf(format, args...)
	panic(fatalErr)
}

func Fatal(v interface{}) {
	//glog.Fatalf(format, args...)
	log.Print(v)
	panic(fatalErr)
}
