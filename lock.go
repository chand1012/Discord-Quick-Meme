package main

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"os"
)

// LOCKFILE the path to the lockfile
const LOCKFILE string = "./thread.lock"

func lockFileEqu(input []byte) (bool, error) {
	data, err := ioutil.ReadFile(LOCKFILE)
	if err != nil {
		return false, err
	}
	if bytes.Equal(input, data) {
		return true, nil
	}
	return false, nil
}

func lockFileExists() bool {
	info, err := os.Stat(LOCKFILE)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func lockFileCreate() ([]byte, error) {
	fileData := make([]byte, 8)
	_, err := rand.Read(fileData) // skipcq: GSC-G404

	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(LOCKFILE, fileData, 0644)
	return fileData, err
}
