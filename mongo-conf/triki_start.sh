#!/bin/sh

# GOPATH needs to be set and "$GOPATH/bin" be in PATH.

go install "gopkg.in/triki.v0"

triki.v0 -config="./triki.conf"
