package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Keys type for json keys
type Keys struct {
	BotID  string   `json:"token"`
	Admins []string `json:"admin"`
}

func jsonExtract(filename string) (string, []string, error) {
	jsonfile, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonfile.Close()

	rawjson, _ := ioutil.ReadAll(jsonfile)

	var keys Keys

	json.Unmarshal(rawjson, &keys)

	return keys.BotID, keys.Admins, err
}
