package main

import (
	"os"
	"testing"
)

func TestLockFile(t *testing.T) {
	contents, err := lockFileCreate()

	if err != nil {
		t.Errorf("Got an error when creating lock file: %v", err)
	}

	exists := lockFileExists()

	if !exists {
		t.Error("Expecting lock file to exists, it doesn't!")
	}

	same, err := lockFileEqu(contents)

	if err != nil {
		t.Errorf("Got an error while checking lock file: %v", err)
	}

	if !same {
		t.Error("Expected the same bytes, got different values.")
	}

	err = os.Remove(LOCKFILE)

	if err != nil {
		t.Errorf("Got an error while removing the lock file: %v", err)
	}

}
