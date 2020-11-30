package main

import (
	"testing"
	"time"
)

func TestMillis(t *testing.T) {
	now := time.Now().UnixNano()
	mills := now / 1000000

	time.Sleep(time.Second)

	test := GetMillis()

	if mills > test {
		t.Errorf("Returned milliseconds less than test time. Expected a time less than %d, got %d.", mills, test)
	}
}
