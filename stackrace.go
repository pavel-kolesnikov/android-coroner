package main

import (
	"fmt"
	"log"
	"runtime"
)

type stackTrace struct {
	message  string
	path string
}

func must(e error) {
	if e != nil {
		log.Fatal(newStacktrace(e.Error()))
	}
}
func mustfn(e error, cb func(string)) {
	if e != nil {
		s := fmt.Sprint(newStacktrace(e.Error()))
		cb(s)
		log.Fatal(s)
	}
}

var mustnot = must
var mustnotfn = mustfn

// New function constructs a new `StackTrace` struct by using given panic
// message, absolute path of the caller file and the line number.
func newStacktrace(msg string) *stackTrace {
	_, file, line, _ := runtime.Caller(2)

	return &stackTrace{
		message:  msg,
		path: fmt.Sprintf("%s:%d", file, line),
	}
}

func (s *stackTrace) String() string {
	return fmt.Sprint("[", s.path, "] ", s.message)
}
