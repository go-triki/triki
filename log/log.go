/*
Package log is a gotriki logging package.
*/
package log

import (
	//"github.com/golang/glog"
	"log"
)

func Infof(format string, args ...interface{}) {
	//glog.Infof(format, args...)
	log.Printf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	//glog.Fatalf(format, args...)
	log.Fatalf(format, args...)
}

func Fatal(v interface{}) {
	//glog.Fatalf(format, args...)
	log.Fatal(v)
}
