// random functions that are useful
package main

import (
	"fmt"
	"strings"
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
	var errStr string
	if err != nil {
		errStr = err.Error()
		if strings.Contains(errStr, "redis") {
			if !strings.Contains(errStr, "nil") {
				fmt.Printf("%s: %+v\n", msg, err)
				if shouldPanic {
					panic(err)
				}
			}
		} else {
			fmt.Printf("%s: %+v", msg, err)
			if shouldPanic {
				panic(err)
			}
		}
	}
}
