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

type redisInfo struct {
	Address  string `json:"redis-address"`
	Password string `json:"redis-password"`
	DB       int    `json:"redis-db"`
}

func redisExtract(filename string) (string, string, int, error) {
	jsonfile, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonfile.Close()

	rawjson, _ := ioutil.ReadAll(jsonfile)

	var redisInfo redisInfo

	json.Unmarshal(rawjson, &redisInfo)

	return redisInfo.Address, redisInfo.Password, redisInfo.DB, err
}
