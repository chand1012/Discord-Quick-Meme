package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// type for json keys
type Keys struct {
	BotId string `json:"token"`
}

func jsonExtract(filename string) (string, error) {
	jsonfile, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonfile.Close()

	rawjson, _ := ioutil.ReadAll(jsonfile)

	var keys Keys

	json.Unmarshal(rawjson, &keys)

	return keys.BotId, err
}
