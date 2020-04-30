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
	TopGG  string   `json:"topgg"`
}

func loginExtract(filename string) (string, []string, string, error) {
	jsonfile, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonfile.Close()

	rawjson, err := ioutil.ReadAll(jsonfile)

	if err != nil {
		fmt.Println("Error reading login JSON:", err)
		return "", nil, "", err
	}

	var keys Keys

	json.Unmarshal(rawjson, &keys)

	return keys.BotID, keys.Admins, keys.TopGG, err
}

// this is for subreddit extraction. There will be one attribute for "memes" and the like
type subJSON struct {
	Memes      []string `json:"memes"`
	Text       []string `json:"text"`
	Hentai     []string `json:"hentai"`
	News       []string `json:"news"`
	FiftyFifty []string `json:"fiftyfifty"`
}

func subExtract(filename string) map[string][]string {
	jsonfile, err := os.Open(filename)
	if err != nil {
		panic(err) // these are situations that the bot cannot run if it errors out
	}

	defer jsonfile.Close()

	rawjson, err := ioutil.ReadAll(jsonfile)

	if err != nil {
		panic(err)
	}

	var subJSON subJSON

	json.Unmarshal(rawjson, &subJSON)

	subMap := make(map[string][]string)

	subMap["memes"] = subJSON.Memes
	subMap["text"] = subJSON.Text
	subMap["hentai"] = subJSON.Hentai
	subMap["news"] = subJSON.News
	subMap["fiftyfifty"] = subJSON.FiftyFifty

	return subMap
}

func getAllSubs(filename string) []string {
	var subs []string
	subMap := subExtract(filename)
	for key := range subMap {
		subsFromMap := subMap[key]
		for i := 0; i < len(subsFromMap); i++ {
			subs = append(subs, subsFromMap[i])
		}
	}
	return subs
}

type redisInfo struct {
	Address  string `json:"redis-address"`
	Password string `json:"redis-password"`
	DB       int    `json:"redis-db"`
}

func redisExtract(filename string) (string, string, int, error) {
	jsonfile, err := os.Open(filename)

	if err != nil {
		panic(err)
	}

	defer jsonfile.Close()

	rawjson, err := ioutil.ReadAll(jsonfile)

	if err != nil {
		panic(err)
	}

	var redisInfo redisInfo

	json.Unmarshal(rawjson, &redisInfo)

	return redisInfo.Address, redisInfo.Password, redisInfo.DB, err
}

type mrisaInfo struct {
	Address string `json:"mrisa"`
}

func mrisaExtract(filename string) string {
	jsonfile, err := os.Open(filename)

	if err != nil {
		panic(err)
	}

	defer jsonfile.Close()

	rawjson, err := ioutil.ReadAll(jsonfile)

	if err != nil {
		panic(err)
	}

	var mrisainfo mrisaInfo

	json.Unmarshal(rawjson, &mrisainfo)

	return mrisainfo.Address
}

type runMode struct {
	Mode string `json:"mode"`
}

func getMode(filename string) string {
	jsonfile, err := os.Open(filename)

	if err != nil {
		panic(err) // this file NEEDS to be there
	}

	defer jsonfile.Close()

	rawjson, err := ioutil.ReadAll(jsonfile)

	if err != nil {
		return "prod"
	}

	var mode runMode

	json.Unmarshal(rawjson, &mode)

	return mode.Mode
}
