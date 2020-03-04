// random functions that are useful
package main

import (
	"fmt"
	"time"
)

// GetMillis gets number of milliseconds since epoch as a 64bit integer
func GetMillis() int64 {
	now := time.Now()
	nanos := now.UnixNano()
	return nanos / 1000000
}

// Error checker that I barely used
func errCheck(msg string, err error, shouldPanic bool) {
	if err != nil {
		fmt.Printf("%s: %+v", msg, err)
		if shouldPanic {
			panic(err)
		}
	}
}
