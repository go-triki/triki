/*
Package log is a gotriki logging package.
*/
package log

import (
	//"github.com/golang/glog"
	"log"
)

// Type that is "paniced" in fatal log functions
type FatalErrorPanic string

func (fatal FatalErrorPanic) Error() string {
	return string(fatal)
}

func (fatal FatalErrorPanic) FatalErrorPanicType() {
	return
}

const (
	fatalErr FatalErrorPanic = "Fatal server error"
)

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
