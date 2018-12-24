package main

import (
	"fmt"
	"runtime"
)

// Assert -- アサート
func Assert(b bool) {
	if !b {
		pc, _, line, _ := runtime.Caller(1)
		name := runtime.FuncForPC(pc).Name()
		panic(fmt.Sprintf("%s():%d: false happened", name, line))
	}
}
